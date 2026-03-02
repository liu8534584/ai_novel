package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"ai_novel/internal/service/agent"
	"ai_novel/internal/service/llm/core"
	"ai_novel/internal/service/rag"
	"ai_novel/models"
	"ai_novel/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"fmt"
)

type PlanHandler struct {
	db           *gorm.DB
	agent        *agent.PlanAgent
	director     *agent.DirectorAgent
	character    *agent.CharacterAgent
	chapterTitle *agent.ChapterTitleAgent
	rag          *rag.MemoryRecallService
}

func NewPlanHandler(db *gorm.DB, agent *agent.PlanAgent, director *agent.DirectorAgent, character *agent.CharacterAgent, chapterTitle *agent.ChapterTitleAgent, ragService *rag.MemoryRecallService) *PlanHandler {
	return &PlanHandler{
		db:           db,
		agent:        agent,
		director:     director,
		character:    character,
		chapterTitle: chapterTitle,
		rag:          ragService,
	}
}

func (h *PlanHandler) ListPlans(c *gin.Context) {
	bookID := c.Param("id")
	var plans []models.OutlineVersion
	if err := h.db.Where("book_id = ?", bookID).Order("version desc, created_at desc").Find(&plans).Error; err != nil {
		response.ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, "Failed to fetch plans")
		return
	}
	response.Success(c, plans)
}

func (h *PlanHandler) GeneratePlans(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid book id")
		return
	}

	var book models.Book
	if err := h.db.First(&book, bookID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.ErrorWithStatus(c, http.StatusNotFound, http.StatusNotFound, "Book not found")
			return
		}
		response.ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, "Failed to fetch book")
		return
	}

	var req struct {
		Description string `json:"description"`
		Genre       string `json:"genre"`
		Chapters    int    `json:"chapters"`
		Count       int    `json:"count"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	description := req.Description
	if description == "" {
		description = book.Description
	}
	genre := req.Genre
	if genre == "" {
		genre = book.Genre
	}
	chapters := req.Chapters
	if chapters <= 0 {
		chapters = book.TotalChapters
	}

	// 0. Build context with book ID
	ctx := core.WithBookID(c.Request.Context(), uint(bookID))

	// 1. Generate World Setting if not already generated or if requested
	worldConfig, err := h.director.InitWorld(ctx, description, genre, chapters)
	if err != nil {
		response.ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, "Failed to generate world setting: "+err.Error())
		return
	}

	// 2. Generate Plan Versions based on World Setting
	versions, err := h.agent.GeneratePlanVersions(ctx, description, genre, worldConfig.Content, chapters, req.Count)
	if err != nil {
		response.ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		// Update book info if changed
		if req.Genre != "" || req.Description != "" || req.Chapters > 0 {
			updates := make(map[string]interface{})
			if req.Genre != "" {
				updates["genre"] = req.Genre
			}
			if req.Description != "" {
				updates["description"] = req.Description
			}
			if req.Chapters > 0 {
				updates["total_chapters"] = req.Chapters
			}
			if len(updates) > 0 {
				if err := tx.Model(&book).Updates(updates).Error; err != nil {
					return err
				}
			}
		}

		// Get current max version
		var lastVersion models.OutlineVersion
		var nextVersionNum int = 1
		if err := tx.Where("book_id = ?", book.ID).Order("version desc").First(&lastVersion).Error; err == nil {
			nextVersionNum = lastVersion.Version + 1
		}

		// Delete old plans that are NOT locked and NOT selected
		if err := tx.Where("book_id = ? AND is_locked = ? AND is_selected = ?", book.ID, false, false).Delete(&models.OutlineVersion{}).Error; err != nil {
			return err
		}

		for _, version := range versions {
			record := models.OutlineVersion{
				BookID:     book.ID,
				Version:    nextVersionNum,
				WorldView:  version.WorldView,
				Outline:    version.Outline,
				Characters: "", // Characters will be generated after selection
				IsSelected: false,
				IsLocked:   false,
			}
			if err := tx.Create(&record).Error; err != nil {
				return err
			}
			nextVersionNum++
		}
		return nil
	})
	if err != nil {
		response.ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, "Failed to save plans")
		return
	}

	var plans []models.OutlineVersion
	if err := h.db.Where("book_id = ?", book.ID).Order("version desc").Find(&plans).Error; err != nil {
		response.ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, "Failed to load plans")
		return
	}
	response.Success(c, plans)
}

func (h *PlanHandler) SelectPlan(c *gin.Context) {
	bookID := c.Param("id")
	planID := c.Param("planId")

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.OutlineVersion{}).Where("book_id = ?", bookID).Update("is_selected", false).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.OutlineVersion{}).Where("id = ? AND book_id = ?", planID, bookID).Update("is_selected", true).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		response.ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, "Failed to select plan")
		return
	}

	// 异步索引计划内容到向量库
	go func() {
		var plan models.OutlineVersion
		if err := h.db.First(&plan, planID).Error; err == nil {
			h.rag.IndexFullPlan(context.Background(), plan.BookID, plan.WorldView, plan.Characters, plan.Outline)
		}
	}()

	response.Success(c, gin.H{"id": planID})
}

func (h *PlanHandler) LockPlan(c *gin.Context) {
	bookID := c.Param("id")
	planID := c.Param("planId")

	var req struct {
		Locked bool `json:"locked"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.db.Model(&models.OutlineVersion{}).Where("id = ? AND book_id = ?", planID, bookID).Update("is_locked", req.Locked).Error; err != nil {
		response.ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, "Failed to update lock status")
		return
	}
	response.Success(c, gin.H{"id": planID, "is_locked": req.Locked})
}

func (h *PlanHandler) UpdatePlan(c *gin.Context) {
	bookID := c.Param("id")
	planID := c.Param("planId")

	var req struct {
		WorldView  string `json:"world_view"`
		Outline    string `json:"outline"`
		Characters string `json:"characters"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.WorldView != "" {
		updates["world_view"] = req.WorldView
	}
	if req.Outline != "" {
		updates["outline"] = req.Outline
	}
	if req.Characters != "" {
		updates["characters"] = req.Characters
	}

	if err := h.db.Model(&models.OutlineVersion{}).Where("id = ? AND book_id = ?", planID, bookID).Updates(updates).Error; err != nil {
		response.ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, "Failed to update plan")
		return
	}
	response.Success(c, gin.H{"id": planID})
}

// GenerateCharacters 为选中的计划生成角色
func (h *PlanHandler) GenerateCharacters(c *gin.Context) {
	bookID := c.Param("id")
	var plan models.OutlineVersion
	if err := h.db.Where("book_id = ? AND is_selected = ?", bookID, true).First(&plan).Error; err != nil {
		response.Error(c, http.StatusNotFound, "No selected plan found")
		return
	}

	if plan.IsLocked {
		response.Error(c, http.StatusForbidden, "Plan is locked and cannot be regenerated")
		return
	}

	bookIDUint, _ := strconv.ParseUint(bookID, 10, 64)
	ctx := core.WithBookID(c.Request.Context(), uint(bookIDUint))
	characters, err := h.character.GenerateCharacters(ctx, plan.WorldView, plan.Outline)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to generate characters: "+err.Error())
		return
	}

	// 序列化为 JSON 存储在 OutlineVersion 表中作为备份
	charJSON, _ := json.Marshal(characters)
	if err := h.db.Model(&plan).Update("characters", string(charJSON)).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to save characters to plan")
		return
	}

	// 持久化到 Character 表和 CharacterAnchor 表
	err = h.db.Transaction(func(tx *gorm.DB) error {
		// 先清理该书旧的角色（如果需要覆盖生成）
		if err := tx.Where("book_id = ?", bookID).Delete(&models.Character{}).Error; err != nil {
			return err
		}

		for _, char := range characters {
			// 1. 创建角色基础信息
			dbChar := models.Character{
				BookID:      plan.BookID,
				Name:        char.Name,
				Role:        char.Role,
				Description: char.Description,
			}
			if err := tx.Create(&dbChar).Error; err != nil {
				return err
			}

			// 2. 创建角色性格锚点
			dbAnchor := models.CharacterAnchor{
				CharacterID:        dbChar.ID,
				PersonalityLabels:  char.Anchor.PersonalityLabels,
				CoreMotivation:     char.Anchor.CoreMotivation,
				BehaviorBottomLine: char.Anchor.BehaviorBottomLine,
				DecisionTendency:   char.Anchor.DecisionTendency,
				EmotionalTriggers:  char.Anchor.EmotionalTriggers,
			}
			if err := tx.Create(&dbAnchor).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to persist characters: "+err.Error())
		return
	}

	response.Success(c, gin.H{"characters": characters})
}

// GenerateChapterTitles 为选中的计划生成章节标题
func (h *PlanHandler) GenerateChapterTitles(c *gin.Context) {
	bookID := c.Param("id")
	var book models.Book
	if err := h.db.First(&book, bookID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Book not found")
		return
	}

	var plan models.OutlineVersion
	if err := h.db.Where("book_id = ? AND is_selected = ?", bookID, true).First(&plan).Error; err != nil {
		response.Error(c, http.StatusNotFound, "No selected plan found")
		return
	}

	if plan.IsLocked {
		response.Error(c, http.StatusForbidden, "Plan is locked and cannot be regenerated")
		return
	}

	ctx := core.WithBookID(c.Request.Context(), book.ID)
	titlePlan, err := h.chapterTitle.GenerateChapterTitlePlan(ctx, plan.Outline, book.TotalChapters)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to generate chapter title plan: "+err.Error())
		return
	}

	batchSize := 10
	totalChapters := book.TotalChapters
	currentCount := 0
	var allTitles strings.Builder
	var previousTitles strings.Builder

	for currentCount < totalChapters {
		remaining := totalChapters - currentCount
		currentBatch := batchSize
		if remaining < batchSize {
			currentBatch = remaining
		}
		startChapter := currentCount + 1

		batchTitles, err := h.chapterTitle.GenerateChapterTitlesBatch(
			ctx,
			plan.WorldView,
			titlePlan,
			plan.Characters,
			startChapter,
			currentBatch,
			currentCount,
			previousTitles.String(),
		)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to generate chapter titles batch: "+err.Error())
			return
		}

		batchTitles = strings.TrimSpace(batchTitles)
		if batchTitles != "" {
			if allTitles.Len() > 0 {
				allTitles.WriteString("\n")
			}
			allTitles.WriteString(batchTitles)
			if previousTitles.Len() > 0 {
				previousTitles.WriteString("\n")
			}
			previousTitles.WriteString(batchTitles)
		}

		currentCount += currentBatch
	}

	titles := allTitles.String()
	if err := h.db.Model(&plan).Update("titles", titles).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to save chapter titles")
		return
	}

	// 尝试解析章节标题并同步到 chapters 表
	go h.syncChapters(book.ID, titles)

	response.Success(c, gin.H{"titles": titles})
}

// GenerateChapterTitlesBatch 分批生成章节标题
func (h *PlanHandler) GenerateChapterTitlesBatch(c *gin.Context) {
	bookID := c.Param("id")
	var req struct {
		BatchSize int `json:"batchSize"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.BatchSize <= 0 {
		req.BatchSize = 10 // 默认生成10章
	}

	var book models.Book
	if err := h.db.First(&book, bookID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Book not found")
		return
	}

	var plan models.OutlineVersion
	if err := h.db.Where("book_id = ? AND is_selected = ?", bookID, true).First(&plan).Error; err != nil {
		response.Error(c, http.StatusNotFound, "No selected plan found")
		return
	}

	// 获取当前已有的章节
	var currentChapters []models.Chapter
	if err := h.db.Where("book_id = ?", bookID).Order("`order` ASC").Find(&currentChapters).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch current chapters")
		return
	}

	currentCount := len(currentChapters)
	startChapter := currentCount + 1

	// 拼接前文标题作为参考
	var previousTitles strings.Builder
	for _, ch := range currentChapters {
		previousTitles.WriteString(fmt.Sprintf("第%d章：%s\n", ch.Order, ch.Title))
	}

	ctx := core.WithBookID(c.Request.Context(), book.ID)
	titlePlan, err := h.chapterTitle.GenerateChapterTitlePlan(ctx, plan.Outline, book.TotalChapters)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to generate chapter title plan: "+err.Error())
		return
	}
	titles, err := h.chapterTitle.GenerateChapterTitlesBatch(
		ctx,
		plan.WorldView,
		titlePlan,
		plan.Characters,
		startChapter,
		req.BatchSize,
		currentCount,
		previousTitles.String(),
	)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to generate chapter titles batch: "+err.Error())
		return
	}

	// 解析并保存新生成的章节
	newChapters := h.parseBatchTitles(book.ID, titles, startChapter)
	if len(newChapters) > 0 {
		if err := h.db.Create(&newChapters).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to save new chapters: "+err.Error())
			return
		}
	}

	// 更新 plan 中的 titles（追加或更新）
	updatedTitles := plan.Titles
	if updatedTitles != "" {
		updatedTitles += "\n" + titles
	} else {
		updatedTitles = titles
	}
	h.db.Model(&plan).Update("titles", updatedTitles)

	response.Success(c, gin.H{
		"newTitles":   titles,
		"newChapters": newChapters,
	})
}

func (h *PlanHandler) parseBatchTitles(bookID uint, titlesText string, startOrder int) []models.Chapter {
	lines := strings.Split(titlesText, "\n")
	var chapters []models.Chapter
	order := startOrder
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// 简单的清理：如果包含 "第x章："，尝试提取后面的内容
		title := line
		if idx := strings.Index(line, "："); idx != -1 {
			title = strings.TrimSpace(line[idx+len("："):])
		} else if idx := strings.Index(line, ":"); idx != -1 {
			title = strings.TrimSpace(line[idx+1:])
		} else if idx := strings.Index(line, " "); idx != -1 {
			// 如果有空格，尝试看前面是不是第x章
			firstPart := line[:idx]
			if strings.Contains(firstPart, "第") && strings.Contains(firstPart, "章") {
				title = strings.TrimSpace(line[idx+1:])
			}
		}

		chapters = append(chapters, models.Chapter{
			BookID: bookID,
			Title:  title,
			Order:  order,
		})
		order++
	}
	return chapters
}

func (h *PlanHandler) syncChapters(bookID uint, titlesText string) {

	// 简单的正则表达式解析 "第x章：标题" 或 "第x章 标题"
	// 这里简化处理，按行分割
	lines := strings.Split(titlesText, "\n")
	var chapters []models.Chapter
	order := 1
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// 简单的清理，去掉前面的序号
		title := line
		chapters = append(chapters, models.Chapter{
			BookID: bookID,
			Title:  title,
			Order:  order,
		})
		order++
	}

	if len(chapters) > 0 {
		h.db.Transaction(func(tx *gorm.DB) error {
			// 如果已有章节且有内容，则不覆盖，这里简单处理：只在没有章节时创建
			var count int64
			tx.Model(&models.Chapter{}).Where("book_id = ?", bookID).Count(&count)
			if count == 0 {
				return tx.Create(&chapters).Error
			}
			return nil
		})
	}
}
