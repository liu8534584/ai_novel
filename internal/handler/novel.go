package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ai_novel/internal/model"
	"ai_novel/internal/pkg/sse"
	"ai_novel/internal/service"
	"ai_novel/internal/service/agent"
	svccontext "ai_novel/internal/service/context"
	"ai_novel/internal/service/llm/core"
	"ai_novel/internal/service/rag"
	"ai_novel/models"
	"ai_novel/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NovelHandler struct {
	db             *gorm.DB
	director       *agent.DirectorAgent
	outliner       *agent.OutlinerAgent
	writer         *agent.WriterAgent
	state          *agent.StateAgent
	foresight      *agent.ForesightAgent
	consistency    *agent.ConsistencyAgent
	summarizer     *agent.SummarizerAgent
	contextManager *svccontext.ContextManager
	rag            *rag.MemoryRecallService
	processor      *service.PostWriteProcessor
}

func NewNovelHandler(db *gorm.DB, director *agent.DirectorAgent, outliner *agent.OutlinerAgent, writer *agent.WriterAgent, state *agent.StateAgent, foresight *agent.ForesightAgent, consistency *agent.ConsistencyAgent, summarizer *agent.SummarizerAgent, ctxMgr *svccontext.ContextManager, ragService *rag.MemoryRecallService, processor *service.PostWriteProcessor) *NovelHandler {
	return &NovelHandler{
		db:             db,
		director:       director,
		outliner:       outliner,
		writer:         writer,
		state:          state,
		foresight:      foresight,
		consistency:    consistency,
		summarizer:     summarizer,
		contextManager: ctxMgr,
		rag:            ragService,
		processor:      processor,
	}
}

// CreateBook handles the creation of a new book (World Initialization).
func (h *NovelHandler) CreateBook(c *gin.Context) {
	var req struct {
		Idea     string `json:"idea" binding:"required"`
		Genre    string `json:"genre"`
		Chapters int    `json:"chapters"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Chapters <= 0 {
		req.Chapters = 100 // Default
	}

	// 1. Call Director Agent to initialize world
	worldConfig, err := h.director.InitWorld(c.Request.Context(), req.Idea, req.Genre, req.Chapters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize world: " + err.Error()})
		return
	}

	// 2. Create Book record
	book := models.Book{
		Title:         "New Story", // Default title, will be updated by plan
		Genre:         req.Genre,
		Description:   req.Idea,
		TotalChapters: req.Chapters,
		Status:        "planning",
		WorldSetting: models.WorldSetting{
			Genre:       req.Genre,
			Description: req.Idea,
			Summary:     worldConfig.Content, // 存储生成的世界观文档
		},
		CurrentState: models.CurrentState{
			ChapterIndex: 0,
			Summary:      "Story Start",
		},
	}

	if err := h.db.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save book: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, book)
}

// GetState returns the current state of the book.
func (h *NovelHandler) GetState(c *gin.Context) {
	bookID := c.Param("id")
	var book models.Book
	if err := h.db.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, book.CurrentState)
}

// GenerateOutline generates the outline for a chapter.
func (h *NovelHandler) GenerateOutline(c *gin.Context) {
	chapterID := c.Param("id")
	var chapter models.Chapter
	if err := h.db.First(&chapter, chapterID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
		return
	}

	var book models.Book
	if err := h.db.First(&book, chapter.BookID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Book not found"})
		return
	}

	// 1. 获取上下文
	var prevChapter models.Chapter
	prevSummary := "无"
	if chapter.Order > 1 {
		if err := h.db.Where("book_id = ? AND order = ?", chapter.BookID, chapter.Order-1).First(&prevChapter).Error; err == nil {
			prevSummary = prevChapter.Summary
		}
	}

	// 获取选中的计划以获取主线大纲和角色设定
	var plan models.OutlineVersion
	if err := h.db.Where("book_id = ? AND is_selected = ?", book.ID, true).First(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Selected plan not found, please select a plan first"})
		return
	}

	currentStateJSON, _ := json.Marshal(book.CurrentState)
	ctx := core.WithBookID(c.Request.Context(), book.ID)

	// 2. 设置 SSE
	sse.SetHeaders(c.Writer)

	// 3. 调用 Outliner Agent（流式）
	streamChan, err := h.outliner.GenerateOutlineStream(
		ctx,
		book.WorldSetting.Summary,
		plan.Characters,
		plan.Outline,
		chapter.Title,
		plan.Titles,
		string(currentStateJSON),
		prevSummary,
		chapter.Order,
		chapter.UserIntent,
	)
	if err != nil {
		sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Failed to generate outline: " + err.Error()})
		return
	}

	// 4. 流式输出并收集完整文本
	fullOutlineText := ""
	streamErr := ""
	for chunk := range streamChan {
		if chunk.Error != "" {
			streamErr = chunk.Error
			break
		}
		if chunk.Content == "" {
			continue
		}
		sse.SendText(c.Writer, chunk.Content)
		fullOutlineText += chunk.Content
		c.Writer.Flush()
	}
	if streamErr != "" {
		sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Outline generation failed: " + streamErr})
		return
	}
	if strings.TrimSpace(fullOutlineText) == "" {
		sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Outline generation returned empty content"})
		return
	}

	// 5. 解析结果并更新章节大纲
	var outline model.ChapterOutline
	cleanJSON := core.ParseJSON(fullOutlineText)
	if err := json.Unmarshal([]byte(cleanJSON), &outline); err != nil {
		sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Failed to parse outline JSON: " + err.Error()})
		return
	}

	outlineJSON, _ := json.Marshal(outline)
	updates := map[string]interface{}{
		"outline":              string(outlineJSON),
		"is_outline_confirmed": false, // 重新生成大纲时，重置确认状态
	}
	if err := h.db.Model(&chapter).Updates(updates).Error; err != nil {
		sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Failed to save outline: " + err.Error()})
		return
	}

	sse.Send(c.Writer, sse.Message{Event: sse.EventEnd, Data: "Outline generation completed"})
}

// ConfirmOutline 确认章节大纲
func (h *NovelHandler) ConfirmOutline(c *gin.Context) {
	chapterID := c.Param("id")
	var chapter models.Chapter
	if err := h.db.First(&chapter, chapterID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
		return
	}

	if chapter.Outline == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "必须先生成大纲才能确认"})
		return
	}

	if err := h.db.Model(&chapter).Update("is_outline_confirmed", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm outline: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Outline confirmed", "is_outline_confirmed": true})
}

// WriteChapter handles the streaming generation of a chapter using layered context.
func (h *NovelHandler) WriteChapter(c *gin.Context) {
	chapterIDStr := c.Param("id")
	chapterID, _ := strconv.Atoi(chapterIDStr)

	var chapter models.Chapter
	if err := h.db.First(&chapter, chapterID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
		return
	}

	// 校验大纲是否已确认
	if !chapter.IsOutlineConfirmed {
		c.JSON(http.StatusForbidden, gin.H{"error": "必须先确认章节大纲才能续写"})
		return
	}

	// 1. 构建分层上下文
	ctx := core.WithBookID(c.Request.Context(), chapter.BookID)
	wCtx, err := h.contextManager.BuildChapterContext(ctx, chapter.BookID, uint(chapterID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to build context: " + err.Error()})
		return
	}

	// 2. 设置 SSE
	sse.SetHeaders(c.Writer)

	// 3. 流式生成
	streamChan, err := h.writer.WriteChapterStream(ctx, wCtx)
	if err != nil {
		sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: err.Error()})
		return
	}

	fullContent := ""
	streamErr := ""
	for chunk := range streamChan {
		if chunk.Error != "" {
			streamErr = chunk.Error
			break
		}
		if chunk.Content == "" {
			continue
		}
		sse.SendText(c.Writer, chunk.Content)
		fullContent += chunk.Content
		c.Writer.Flush()
	}
	if streamErr != "" {
		sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Chapter generation failed: " + streamErr})
		return
	}
	if fullContent == "" {
		sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Chapter generation returned empty content"})
		return
	}

	// 4. 更新章节内容并创建版本
	chapter.Content = fullContent
	wordCount := len([]rune(fullContent))

	const maxVersionRetries = 3
	for attempt := 0; attempt < maxVersionRetries; attempt++ {
		err = h.db.Transaction(func(tx *gorm.DB) error {
			// 获取最新版本号
			var lastVersion models.ChapterVersion
			nextVersionNum := 1
			if err := tx.Where("chapter_id = ?", chapter.ID).Order("version desc").First(&lastVersion).Error; err == nil {
				nextVersionNum = lastVersion.Version + 1
			}

			// 创建版本记录（并发下依赖唯一索引兜底）
			version := models.ChapterVersion{
				ChapterID: chapter.ID,
				Version:   nextVersionNum,
				Title:     chapter.Title,
				Content:   fullContent,
				WordCount: wordCount,
			}
			if err := tx.Create(&version).Error; err != nil {
				return err
			}

			// 更新章节主表
			if err := tx.Model(&chapter).Updates(map[string]interface{}{
				"content":         fullContent,
				"current_version": nextVersionNum,
			}).Error; err != nil {
				return err
			}

			// 更新书籍的当前状态
			if err := tx.Model(&models.Book{}).Where("id = ?", chapter.BookID).Update("current_state", models.CurrentState{
				ChapterIndex: chapter.Order,
				Summary:      fmt.Sprintf("Completed Chapter %d: %s (Version: %d, Words: %d)", chapter.Order, chapter.Title, nextVersionNum, wordCount),
			}).Error; err != nil {
				return err
			}

			return nil
		})
		if err == nil {
			break
		}
		if !isDuplicateVersionErr(err) {
			break
		}
	}

	if err != nil {
		sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Failed to persist generated chapter: " + err.Error()})
		return
	}

	// 5. 异步后处理（统一后处理器，包含状态提取/事件/伏笔/OOC/矛盾/健康度/摘要/RAG）
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		h.processor.Process(bgCtx, chapter.BookID, uint(chapterID), chapter.Order, fullContent)
	}()

	sse.Send(c.Writer, sse.Message{Event: sse.EventEnd, Data: "Chapter generation completed"})
}

func isDuplicateVersionErr(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique constraint") || strings.Contains(msg, "duplicate key")
}

// GetOOCScores 获取章节 OOC 评分
func (h *NovelHandler) GetOOCScores(c *gin.Context) {
	chapterID := c.Param("id")
	var scores []models.OOCScore
	if err := h.db.Where("chapter_id = ?", chapterID).Find(&scores).Error; err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, scores)
}

// GetContradictions 获取章节剧情矛盾
func (h *NovelHandler) GetContradictions(c *gin.Context) {
	chapterID := c.Param("id")
	var contradictions []models.StoryContradiction
	if err := h.db.Where("chapter_id = ?", chapterID).Find(&contradictions).Error; err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, contradictions)
}

// GetChapterHealth 获取章节健康度评分
func (h *NovelHandler) GetChapterHealth(c *gin.Context) {
	chapterID := c.Param("id")
	var health models.ChapterHealthScore
	if err := h.db.Where("chapter_id = ?", chapterID).First(&health).Error; err != nil {
		response.Error(c, 404, "Health score not found")
		return
	}
	response.Success(c, health)
}

// ListChapterVersions 获取章节的所有版本
func (h *NovelHandler) ListChapterVersions(c *gin.Context) {
	chapterID := c.Param("id")
	var versions []models.ChapterVersion
	if err := h.db.Where("chapter_id = ?", chapterID).Order("version desc").Find(&versions).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch chapter versions")
		return
	}
	response.Success(c, versions)
}

// RollbackChapter 回滚章节到指定版本
func (h *NovelHandler) RollbackChapter(c *gin.Context) {
	chapterID := c.Param("id")
	versionIDStr := c.Param("versionId")
	versionID, _ := strconv.Atoi(versionIDStr)

	var version models.ChapterVersion
	if err := h.db.Where("chapter_id = ? AND version = ?", chapterID, versionID).First(&version).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Version not found")
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		// 更新章节主表内容
		if err := tx.Model(&models.Chapter{}).Where("id = ?", chapterID).Updates(map[string]interface{}{
			"content":         version.Content,
			"current_version": version.Version,
		}).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to rollback chapter")
		return
	}

	response.Success(c, gin.H{"chapter_id": chapterID, "version": versionID})
}

// GetForeshadowingAlerts 获取书籍的伏笔预警列表
func (h *NovelHandler) GetForeshadowingAlerts(c *gin.Context) {
	bookIDStr := c.Param("id")
	bookID, _ := strconv.Atoi(bookIDStr)

	// 获取书籍当前进度
	var book models.Book
	if err := h.db.First(&book, bookID).Error; err != nil {
		response.Error(c, 404, "Book not found")
		return
	}

	alerts, err := h.foresight.GetOpenForeshadowingWithAlerts(uint(bookID), book.CurrentState.ChapterIndex)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, alerts)
}
