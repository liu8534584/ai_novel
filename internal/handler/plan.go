package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	internalmodel "ai_novel/internal/model"
	"ai_novel/internal/pkg/sse"
	"ai_novel/internal/service/agent"
	"ai_novel/internal/service/llm/core"
	"ai_novel/internal/service/rag"
	"ai_novel/models"
	"ai_novel/pkg/logger"
	"ai_novel/pkg/response"

	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	count := req.Count
	if count <= 0 {
		count = 3
	}

	// 0. Build context with book ID
	ctx := core.WithBookID(c.Request.Context(), uint(bookID))
	sse.SetHeaders(c.Writer)
	_ = sse.Send(c.Writer, sse.Message{Event: sse.EventStateUpdate, Data: gin.H{"phase": "world_start"}})

	// 1. Generate World Setting (Streaming)
	worldStream, err := h.director.InitWorldStream(ctx, description, genre, chapters)
	if err != nil {
		_ = sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Failed to start world generation: " + err.Error()})
		return
	}

	worldConfigContent, err := streamTextToSSE(c, worldStream)
	if err != nil {
		_ = sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Failed to generate world setting: " + err.Error()})
		return
	}
	worldConfig := &internalmodel.WorldConfig{Content: worldConfigContent}

	_ = sse.Send(c.Writer, sse.Message{Event: sse.EventStateUpdate, Data: gin.H{"phase": "world_done"}})

	// 2. Stream Generate Plan Versions based on World Setting
	versions := make([]internalmodel.OutlineVersion, 0, count)
	for i := 0; i < count; i++ {
		_ = sse.Send(c.Writer, sse.Message{Event: sse.EventStateUpdate, Data: gin.H{"phase": "plan_version_start", "index": i + 1}})
		streamChan, streamErr := h.agent.GeneratePlanVersionStream(ctx, description, genre, worldConfig.Content, chapters, i)
		if streamErr != nil {
			_ = sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Failed to generate plan version: " + streamErr.Error()})
			return
		}
		text, err := streamTextToSSE(c, streamChan)
		if err != nil {
			_ = sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Failed to generate plan version: " + err.Error()})
			return
		}
		versions = append(versions, internalmodel.OutlineVersion{
			WorldView: worldConfig.Content,
			Outline:   strings.TrimSpace(text),
		})
		_ = sse.Send(c.Writer, sse.Message{Event: sse.EventStateUpdate, Data: gin.H{"phase": "plan_version_end", "index": i + 1}})
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
		_ = sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Failed to save plans"})
		return
	}

	_ = sse.Send(c.Writer, sse.Message{Event: sse.EventEnd, Data: "Plans generation completed"})
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

	var plan models.OutlineVersion
	if err := h.db.First(&plan, planID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Plan not found")
		return
	}

	if plan.IsLocked {
		response.Error(c, http.StatusForbidden, "Plan is locked and cannot be edited manually. Please unlock first.")
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

	// Note: We allow generating/regenerating characters even if the plan is locked.
	// The lock is primarily to prevent accidental manual changes to the outline/worldview.

	bookIDUint, _ := strconv.ParseUint(bookID, 10, 64)
	ctx := core.WithBookID(c.Request.Context(), uint(bookIDUint))
	sse.SetHeaders(c.Writer)

	streamChan, err := h.character.GenerateCharactersStream(ctx, plan.WorldView, plan.Outline)
	if err != nil {
		_ = sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Failed to generate characters: " + err.Error()})
		return
	}
	rawText, err := streamTextToSSE(c, streamChan)
	if err != nil {
		_ = sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Failed to generate characters: " + err.Error()})
		return
	}

	var characters []internalmodel.Character
	if strings.TrimSpace(rawText) == "" {
		generated, genErr := h.character.GenerateCharacters(ctx, plan.WorldView, plan.Outline)
		if genErr != nil {
			_ = sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Failed to generate characters: " + genErr.Error()})
			return
		}
		for _, c := range generated {
			characters = append(characters, internalmodel.Character{
				Name:        c.Name,
				Role:        c.Role,
				Description: c.Description,
				Anchor:      c.Anchor,
			})
		}
		charJSONBytes, _ := json.Marshal(characters)
		_ = sse.SendText(c.Writer, string(charJSONBytes))
	} else {
		parsed, parseErr := parseCharactersFromText(rawText)
		if parseErr != nil {
			_ = sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Failed to parse characters: " + parseErr.Error()})
			return
		}
		characters = parsed
	}

	// 序列化为 JSON 存储在 OutlineVersion 表中作为备份
	charJSON, _ := json.Marshal(characters)
	if err := h.db.Model(&plan).Update("characters", string(charJSON)).Error; err != nil {
		_ = sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Failed to save characters to plan"})
		return
	}

	// 持久化到 Character 表和 CharacterAnchor 表
	err = h.db.Transaction(func(tx *gorm.DB) error {
		// 级联清理：先删除角色关联数据
		var charIDs []uint
		tx.Model(&models.Character{}).Where("book_id = ?", bookID).Pluck("id", &charIDs)
		if len(charIDs) > 0 {
			tx.Where("character_id IN ?", charIDs).Delete(&models.CharacterAnchor{})
			tx.Where("character_id IN ?", charIDs).Delete(&models.CharacterStateRecord{})
			tx.Where("character_id IN ?", charIDs).Delete(&models.OOCScore{})
		}
		// 再删除角色本身
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
		_ = sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Failed to persist characters: " + err.Error()})
		return
	}

	_ = sse.Send(c.Writer, sse.Message{Event: sse.EventEnd, Data: "Characters generation completed"})
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

	// Note: We allow generating/regenerating titles even if the plan is locked.
	// The lock is primarily for the outline/worldview content itself.

	ctx := core.WithBookID(c.Request.Context(), book.ID)
	sse.SetHeaders(c.Writer)

	_ = sse.Send(c.Writer, sse.Message{Event: sse.EventStateUpdate, Data: gin.H{"phase": "title_plan_start"}})
	titlePlanStream, err := h.chapterTitle.GenerateChapterTitlePlanStream(ctx, plan.Outline, book.TotalChapters)
	if err != nil {
		logger.Error("Failed to start chapter title plan stream: %v", err)
		_ = sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Failed to generate chapter title plan: " + err.Error()})
		return
	}
	titlePlan, err := streamTextToSSE(c, titlePlanStream)
	if err != nil {
		logger.Error("Failed to generate chapter title plan from stream: %v", err)
		_ = sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Failed to generate chapter title plan: " + err.Error()})
		return
	}
	logger.Info("Chapter title plan generated successfully (len: %d)", len(titlePlan))
	_ = sse.Send(c.Writer, sse.Message{Event: sse.EventStateUpdate, Data: gin.H{"phase": "title_plan_done"}})

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

		_ = sse.Send(c.Writer, sse.Message{Event: sse.EventStateUpdate, Data: gin.H{"phase": "titles_batch_start", "start": startChapter, "size": currentBatch}})
		streamChan, streamErr := h.chapterTitle.GenerateChapterTitlesBatchStream(
			ctx,
			plan.WorldView,
			titlePlan,
			plan.Characters,
			startChapter,
			currentBatch,
			currentCount,
			previousTitles.String(),
		)
		if streamErr != nil {
			_ = sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Failed to generate chapter titles batch: " + streamErr.Error()})
			return
		}
		batchTitles, err := streamTextToSSE(c, streamChan)
		if err != nil {
			_ = sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Failed to generate chapter titles batch: " + err.Error()})
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
		_ = sse.Send(c.Writer, sse.Message{Event: sse.EventStateUpdate, Data: gin.H{"phase": "titles_batch_end", "start": startChapter}})
	}

	titles := allTitles.String()
	if err := h.db.Model(&plan).Update("titles", titles).Error; err != nil {
		_ = sse.Send(c.Writer, sse.Message{Event: sse.EventError, Data: "Failed to save chapter titles"})
		return
	}

	// 尝试解析章节标题并同步到 chapters 表
	go h.syncChapters(book.ID, titles)

	_ = sse.Send(c.Writer, sse.Message{Event: sse.EventEnd, Data: "Chapter titles generation completed"})
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
	if err := h.db.Where("book_id = ?", bookID).Order("\"order\" ASC").Find(&currentChapters).Error; err != nil {
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
		if line == "" || len(line) < 2 {
			continue
		}
		
		// 常见的干扰行：以“确定”、“好”、“根据”等开头的说明性文字
		if strings.HasPrefix(line, "好的") || strings.HasPrefix(line, "这里是") || strings.HasPrefix(line, "如下是") || strings.HasSuffix(line, "：") || strings.HasSuffix(line, ":") {
			if !strings.Contains(line, "第") {
				continue
			}
		}

		// 简单的清理：如果包含 "第x章："，尝试提取后面的内容
		title := line
		if idx := strings.Index(line, "："); idx != -1 {
			possibleTitle := strings.TrimSpace(line[idx+len("："):])
			if possibleTitle != "" {
				title = possibleTitle
			}
		} else if idx := strings.Index(line, ":"); idx != -1 {
			possibleTitle := strings.TrimSpace(line[idx+1:])
			if possibleTitle != "" {
				title = possibleTitle
			}
		} else if idx := strings.Index(line, " "); idx != -1 {
			// 如果有空格，尝试看前面是不是第x章
			firstPart := line[:idx]
			if strings.Contains(firstPart, "第") && strings.Contains(firstPart, "章") {
				title = strings.TrimSpace(line[idx+1:])
			}
		}

		// 如果提取后的标题依然太长（例如超过50个字符），可能不是标题而是描述，略过
		if len([]rune(title)) > 50 {
			continue
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
	chapters := h.parseBatchTitles(bookID, titlesText, 1)

	if len(chapters) > 0 {
		h.db.Transaction(func(tx *gorm.DB) error {
			// 同步策略：
			// 1) 已存在章节：仅在“未写正文”时更新标题，避免覆盖已写内容
			// 2) 不存在章节：创建缺失章节
			var existing []models.Chapter
			if err := tx.Where("book_id = ?", bookID).Find(&existing).Error; err != nil {
				return err
			}
			byOrder := make(map[int]models.Chapter, len(existing))
			for _, ch := range existing {
				byOrder[ch.Order] = ch
			}

			var toCreate []models.Chapter
			for _, target := range chapters {
				if old, ok := byOrder[target.Order]; ok {
					if old.Content == "" && old.Title != target.Title {
						if err := tx.Model(&models.Chapter{}).Where("id = ?", old.ID).Update("title", target.Title).Error; err != nil {
							return err
						}
					}
					continue
				}
				toCreate = append(toCreate, target)
			}
			if len(toCreate) > 0 {
				if err := tx.Create(&toCreate).Error; err != nil {
					return err
				}
			}
			return nil
		})
	}
}

func streamTextToSSE(c *gin.Context, streamChan <-chan core.StreamResponse) (string, error) {
	var sb strings.Builder
	for chunk := range streamChan {
		if chunk.Error != "" {
			return "", errors.New(chunk.Error)
		}
		if chunk.Content == "" {
			continue
		}
		if err := sse.SendText(c.Writer, chunk.Content); err != nil {
			return "", err
		}
		sb.WriteString(chunk.Content)
		c.Writer.Flush()
	}
	return sb.String(), nil
}

func parseCharactersFromText(raw string) ([]internalmodel.Character, error) {
	content := core.ParseJSON(raw)
	if strings.TrimSpace(content) == "" {
		// Log the original raw content for debugging
		logger.Error("Failed to parse characters from text. Raw content might not contain JSON: %s", raw)
		return nil, fmt.Errorf("failed to extract valid JSON from LLM output. Output was: %s", truncateString(raw, 100))
	}

	var characters []internalmodel.Character
	// 1. Try unmarshaling as object with "characters" key (Preferred)
	var wrapper struct {
		Characters []internalmodel.Character `json:"characters"`
	}
	if err := json.Unmarshal([]byte(content), &wrapper); err == nil && len(wrapper.Characters) > 0 {
		return wrapper.Characters, nil
	}

	// 2. Try to unmarshal directly as array
	if err := json.Unmarshal([]byte(content), &characters); err == nil && len(characters) > 0 {
		return characters, nil
	}

	// 3. If that also fails, try to find the first '[' and last ']' and unmarshal as array
	firstArray := strings.Index(content, "[")
	lastArray := strings.LastIndex(content, "]")
	if firstArray != -1 && lastArray != -1 && lastArray > firstArray {
		arrayContent := content[firstArray : lastArray+1]
		if err := json.Unmarshal([]byte(arrayContent), &characters); err == nil && len(characters) > 0 {
			return characters, nil
		}
	}

	return nil, fmt.Errorf("failed to parse characters from JSON. Content: %s", truncateString(content, 200))
}

func truncateString(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
