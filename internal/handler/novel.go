package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"ai_novel/internal/pkg/sse"
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
}

func NewNovelHandler(db *gorm.DB, director *agent.DirectorAgent, outliner *agent.OutlinerAgent, writer *agent.WriterAgent, state *agent.StateAgent, foresight *agent.ForesightAgent, consistency *agent.ConsistencyAgent, summarizer *agent.SummarizerAgent, ctxMgr *svccontext.ContextManager, ragService *rag.MemoryRecallService) *NovelHandler {
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

	// 2. 调用 Outliner Agent
	outline, err := h.outliner.GenerateOutline(
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate outline: " + err.Error()})
		return
	}

	// 3. 更新章节大纲
	outlineJSON, _ := json.Marshal(outline)
	updates := map[string]interface{}{
		"outline":              string(outlineJSON),
		"is_outline_confirmed": false, // 重新生成大纲时，重置确认状态
	}
	if err := h.db.Model(&chapter).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save outline: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, outline)
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
	for textChunk := range streamChan {
		sse.SendText(c.Writer, textChunk)
		fullContent += textChunk
		c.Writer.Flush()
	}

	// 4. 更新章节内容并创建版本
	chapter.Content = fullContent
	wordCount := len([]rune(fullContent))

	err = h.db.Transaction(func(tx *gorm.DB) error {
		// 获取最新版本号
		var lastVersion models.ChapterVersion
		nextVersionNum := 1
		if err := tx.Where("chapter_id = ?", chapter.ID).Order("version desc").First(&lastVersion).Error; err == nil {
			nextVersionNum = lastVersion.Version + 1
		}

		// 创建版本记录
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

	if err != nil {
		fmt.Printf("Failed to update chapter content and version: %v\n", err)
	}

	// 5. 更新角色动态状态 (章节后处理)
	var book models.Book
	if err := h.db.Preload("Characters").First(&book, chapter.BookID).Error; err == nil {
		// 准备基础资料
		var baseProfiles []string
		var previousStates []string
		for _, char := range book.Characters {
			baseProfiles = append(baseProfiles, fmt.Sprintf("%s: %s (%s)", char.Name, char.Role, char.Description))
			stateJSON, _ := json.Marshal(char.DynamicState)
			previousStates = append(previousStates, fmt.Sprintf("%s: %s", char.Name, string(stateJSON)))
		}

		var events []models.StoryEvent
		updates, err := h.state.ExtractDynamicStateChanges(
			ctx,
			book.WorldSetting.Summary,
			strings.Join(baseProfiles, "\n"),
			strings.Join(previousStates, "\n"),
			fullContent,
		)

		if err == nil {
			for _, char := range book.Characters {
				if update, ok := updates[char.Name]; ok {
					char.DynamicState = update
					// 1. 更新当前状态
					h.db.Model(&char).Update("dynamic_state", char.DynamicState)

					// 2. 记录历史轨迹 (Task 18)
					stateRecord := models.CharacterStateRecord{
						CharacterID: char.ID,
						ChapterID:   uint(chapterID),
						State:       update,
					}
					h.db.Create(&stateRecord)
				}
			}

			// 6. 关键事件抽取与伏笔追踪 (Task 19)
			events, err = h.foresight.ExtractEvents(ctx, fullContent)
			if err == nil {
				// 1. 持久化 StoryEvent 记录 (Task 19)
				for i := range events {
					events[i].BookID = chapter.BookID
					events[i].ChapterID = uint(chapterID)
					events[i].ChapterIndex = chapter.Order
					h.db.Create(&events[i])
				}

				// 2. 准备状态变化的 JSON 用于伏笔回收判断
				updatesJSON, _ := json.Marshal(updates)
				h.foresight.UpdateForeshadowing(ctx, chapter.BookID, uint(chapterID), chapter.Order, fullContent, events, string(updatesJSON))
			}

			// 7. OOC 评分 (章节后处理)
			for _, char := range book.Characters {
				if update, ok := updates[char.Name]; ok {
					// 获取性格锚点
					var anchor models.CharacterAnchor
					if err := h.db.Where("character_id = ?", char.ID).First(&anchor).Error; err != nil {
						// 如果没有锚点，则尝试提取初始锚点
						newAnchor, err := h.consistency.ExtractCharacterAnchor(ctx, &char, "")
						if err == nil {
							anchor = *newAnchor
							h.db.Create(&anchor)
						}
					}

					if anchor.ID != 0 {
						// 准备行为描述
						behavior := fmt.Sprintf("目标: %s\n行为: %s\n情绪: %s", update.Goal, update.KeyActions, update.EmotionalState)
						// 评估 OOC
						score, err := h.consistency.EvaluateOOC(ctx, &anchor, "", behavior)
						if err == nil {
							score.ChapterID = uint(chapterID)
							h.db.Create(score)
						}
					}
				}
			}
		}

		// 8. 剧情矛盾检测 (章节后处理)
		// 使用 RAG 召回历史事实 (Task 23)
		recallQuery := fmt.Sprintf("分析当前章节内容是否存在与历史事实冲突: %s", chapter.Title)
		historyMemory, err := h.rag.Recall(c.Request.Context(), book.ID, recallQuery, 10, "")
		if err != nil {
			fmt.Printf("Warning: failed to recall history memory: %v\n", err)
			historyMemory = "无法召回历史记忆"
		}

		contradictions, err := h.consistency.DetectContradictions(
			c.Request.Context(),
			book.WorldSetting.Rules,
			historyMemory,
			strings.Join(previousStates, "\n"),
			fullContent,
		)
		if err == nil {
			for _, con := range contradictions {
				con.BookID = book.ID
				con.ChapterID = uint(chapterID)
				h.db.Create(&con)
			}
		}

		// 9. 综合健康度评分 (Task 25)
		// 获取 OOC 评分
		var oocScores []models.OOCScore
		h.db.Where("chapter_id = ?", chapterID).Find(&oocScores)

		// 获取伏笔状态 (open 和本章回收的)
		var openForeshadows []models.Foreshadowing
		h.db.Where("book_id = ? AND status = ?", book.ID, "open").Find(&openForeshadows)

		var resolvedForeshadows []models.Foreshadowing
		h.db.Where("book_id = ? AND status = ? AND resolved_chapter_index = ?", book.ID, "resolved", chapter.Order).Find(&resolvedForeshadows)

		healthScore := h.consistency.EvaluateChapterHealth(c.Request.Context(), oocScores, contradictions, openForeshadows, resolvedForeshadows)
		healthScore.BookID = book.ID
		healthScore.ChapterID = uint(chapterID)
		h.db.Create(healthScore)

		// 10. 异步向量化分类入库 (Task 21 + 状态写回)
		go func(bID, cID uint, title, content string, evs []models.StoryEvent, states map[string]models.CharacterDynamicState) {
			ctx := context.Background()
			
			// A. 生成本章摘要 (状态写回)
			summary, err := h.summarizer.SummarizeChapter(ctx, title, content)
			if err == nil {
				// 将摘要存入历史库，作为更精简的记忆
				h.rag.IndexChapter(ctx, bID, cID, title, "[章节摘要] "+summary)
			}

			// B. 章节正文分段向量化
			if err := h.rag.IndexChapter(ctx, bID, cID, title, content); err != nil {
				fmt.Printf("Warning: failed to index chapter: %v\n", err)
			}
			
			// C. 剧情事件向量化
			for _, ev := range evs {
				if err := h.rag.IndexEvent(ctx, bID, cID, ev); err != nil {
					fmt.Printf("Warning: failed to index event: %v\n", err)
				}
			}
			
			// D. 角色状态向量化 (人设库更新)
			for name, state := range states {
				// 构造详细的角色档案片段
				charContent := fmt.Sprintf("角色: %s\n最新状态: %s\n当前目标: %s\n性格/情绪: %s\n能力变化: %s", 
					name, state.IdentityLocation, state.Goal, state.EmotionalState, state.AbilityResourceChanges)
				
				if err := h.rag.IndexCharacter(ctx, bID, name, charContent); err != nil {
					fmt.Printf("Warning: failed to index character state: %v\n", err)
				}
				
				// 保留原有的 SQLite 记录
				h.rag.IndexCharacterState(ctx, bID, cID, name, state)
			}
		}(chapter.BookID, uint(chapterID), chapter.Title, fullContent, events, updates)
	}

	sse.Send(c.Writer, sse.Message{Event: sse.EventEnd, Data: "Chapter generation completed"})
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
