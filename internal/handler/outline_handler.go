package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"ai_novel/internal/service/llm/core"
	"ai_novel/models"
	"ai_novel/pkg/logger"
	"ai_novel/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OutlineHandler struct {
	db          *gorm.DB
	llmProvider core.Provider
}

func NewOutlineHandler(db *gorm.DB, llmProvider core.Provider) *OutlineHandler {
	return &OutlineHandler{db: db, llmProvider: llmProvider}
}

// GetMasterOutline 获取全部章节蓝图列表
func (h *OutlineHandler) GetMasterOutline(c *gin.Context) {
	bookID := c.Param("id")
	var blueprints []models.ChapterBlueprint
	if err := h.db.Where("book_id = ?", bookID).Order("chapter_index asc").Find(&blueprints).Error; err != nil {
		response.ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, "Failed to fetch blueprints")
		return
	}

	// 同时返回书籍信息
	var book models.Book
	h.db.First(&book, bookID)

	response.Success(c, gin.H{
		"book":       book,
		"blueprints": blueprints,
	})
}

// UpdateBlueprint 编辑单章蓝图（自动保存）
func (h *OutlineHandler) UpdateBlueprint(c *gin.Context) {
	bookID := c.Param("id")
	chapterIndex, err := strconv.Atoi(c.Param("chapterIndex"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid chapter index")
		return
	}

	var req struct {
		Title                 *string `json:"title"`
		Summary               *string `json:"summary"`
		ProtagonistMotivation *string `json:"protagonist_motivation"`
		KeyForeshadowing      *string `json:"key_foreshadowing"`
		AppearingCharacters   *string `json:"appearing_characters"`
		Highlight             *string `json:"highlight"`
		CoreEvents            *string `json:"core_events"`
		Challenges            *string `json:"challenges"`
		CharacterChanges      *string `json:"character_changes"`
		WorldChanges          *string `json:"world_changes"`
		NewForeshadowing      *string `json:"new_foreshadowing"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Summary != nil {
		updates["summary"] = *req.Summary
	}
	if req.ProtagonistMotivation != nil {
		updates["protagonist_motivation"] = *req.ProtagonistMotivation
	}
	if req.KeyForeshadowing != nil {
		updates["key_foreshadowing"] = *req.KeyForeshadowing
	}
	if req.AppearingCharacters != nil {
		updates["appearing_characters"] = *req.AppearingCharacters
	}
	if req.Highlight != nil {
		updates["highlight"] = *req.Highlight
	}
	if req.CoreEvents != nil {
		updates["core_events"] = *req.CoreEvents
	}
	if req.Challenges != nil {
		updates["challenges"] = *req.Challenges
	}
	if req.CharacterChanges != nil {
		updates["character_changes"] = *req.CharacterChanges
	}
	if req.WorldChanges != nil {
		updates["world_changes"] = *req.WorldChanges
	}
	if req.NewForeshadowing != nil {
		updates["new_foreshadowing"] = *req.NewForeshadowing
	}

	if len(updates) == 0 {
		response.Success(c, nil)
		return
	}

	if err := h.db.Model(&models.ChapterBlueprint{}).
		Where("book_id = ? AND chapter_index = ?", bookID, chapterIndex).
		Updates(updates).Error; err != nil {
		response.ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, "Failed to update blueprint")
		return
	}

	response.Success(c, gin.H{"chapter_index": chapterIndex})
}

// GenerateBatch 分批生成总纲（章节蓝图）
func (h *OutlineHandler) GenerateBatch(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid book id")
		return
	}

	var req struct {
		StartChapter int `json:"start_chapter"`
		BatchSize    int `json:"batch_size"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if req.BatchSize <= 0 {
		req.BatchSize = 5
	}
	if req.StartChapter <= 0 {
		req.StartChapter = 1
	}

	// 加载书籍 + 选中的方案
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

	// 获取已有蓝图作为上下文
	var existingBlueprints []models.ChapterBlueprint
	h.db.Where("book_id = ? AND chapter_index < ?", bookID, req.StartChapter).Order("chapter_index asc").Find(&existingBlueprints)

	var previousContext string
	for _, bp := range existingBlueprints {
		previousContext += fmt.Sprintf("第%d章 %s: %s\n", bp.ChapterIndex, bp.Title, bp.Summary)
	}

	ctx := core.WithBookID(c.Request.Context(), uint(bookID))

	// 构造 Prompt
	promptText := fmt.Sprintf(`你是一个专业的小说章节规划师。请为以下小说生成第%d章到第%d章的详细章节蓝图。

## 书籍信息
- 标题: %s
- 类型: %s
- 简介: %s

## 世界观设定
%s

## 剧情大纲
%s

## 角色设定
%s

## 前文蓝图 (已规划章节)
%s

请以纯 JSON 数组格式返回，每个元素包含以下字段：
- title: 章节标题（简短有力，如"青云试炼，初露锋芒"）
- summary: 章节大纲（详细描述本章剧情走向，3-5句话）
- protagonist_motivation: 主角动机（本章主角的核心驱动力）
- key_foreshadowing: 关键伏笔（本章埋下或需回收的伏笔）
- appearing_characters: 出场人物（章节中出场的角色及其作用，如"楚天阳（主角）、苏清月（背景提及）"）
- highlight: 审视亮点（本章的叙事风格与看点）
- core_events: 核心事件（本章最重要的情节转折）
- challenges: 面临挑战（主角本章面临的困难或冲突）

仅返回 JSON 数组。`,
		req.StartChapter, req.StartChapter+req.BatchSize-1,
		book.Title, book.Genre, book.Description,
		plan.WorldView,
		plan.Outline,
		plan.Characters,
		previousContext,
	)

	messages := []core.Message{
		{Role: core.RoleUser, Content: promptText},
	}

	options := core.Options{
		Model:    "",
		JSONMode: true,
	}
	core.GetStrategy(core.TaskPlanning).ApplyToOptions(&options)

	resp, err := h.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "LLM call failed: "+err.Error())
		return
	}

	// 解析 JSON 数组
	var blueprintResults []struct {
		Title                 string `json:"title"`
		Summary               string `json:"summary"`
		ProtagonistMotivation string `json:"protagonist_motivation"`
		KeyForeshadowing      string `json:"key_foreshadowing"`
		AppearingCharacters   string `json:"appearing_characters"`
		Highlight             string `json:"highlight"`
		CoreEvents            string `json:"core_events"`
		Challenges            string `json:"challenges"`
	}

	cleanJSON := core.ParseJSON(resp.Content)
	if err := json.Unmarshal([]byte(cleanJSON), &blueprintResults); err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to parse blueprint JSON: %v, raw: %s", err, resp.Content))
		return
	}

	// 持久化并同步到 Chapter 表
	var savedBlueprints []models.ChapterBlueprint
	err = h.db.Transaction(func(tx *gorm.DB) error {
		for i, bp := range blueprintResults {
			chapterIdx := req.StartChapter + i
			blueprint := models.ChapterBlueprint{
				BookID:                uint(bookID),
				ChapterIndex:          chapterIdx,
				Title:                 bp.Title,
				Summary:               bp.Summary,
				ProtagonistMotivation: bp.ProtagonistMotivation,
				KeyForeshadowing:      bp.KeyForeshadowing,
				AppearingCharacters:   bp.AppearingCharacters,
				Highlight:             bp.Highlight,
				CoreEvents:            bp.CoreEvents,
				Challenges:            bp.Challenges,
			}

			// Upsert: 如果已存在则更新，否则创建
			var existing models.ChapterBlueprint
			if err := tx.Where("book_id = ? AND chapter_index = ?", bookID, chapterIdx).First(&existing).Error; err == nil {
				blueprint.Model = existing.Model
				tx.Model(&existing).Updates(map[string]interface{}{
					"title":                  blueprint.Title,
					"summary":                blueprint.Summary,
					"protagonist_motivation": blueprint.ProtagonistMotivation,
					"key_foreshadowing":      blueprint.KeyForeshadowing,
					"appearing_characters":   blueprint.AppearingCharacters,
					"highlight":              blueprint.Highlight,
					"core_events":            blueprint.CoreEvents,
					"challenges":             blueprint.Challenges,
				})
				blueprint.ID = existing.ID
			} else {
				tx.Create(&blueprint)
			}
			savedBlueprints = append(savedBlueprints, blueprint)

			// 同步到 Chapter 表（仅创建不存在的章节）
			var chapter models.Chapter
			if err := tx.Where("book_id = ? AND \"order\" = ?", bookID, chapterIdx).First(&chapter).Error; err != nil {
				tx.Create(&models.Chapter{
					BookID: uint(bookID),
					Title:  bp.Title,
					Order:  chapterIdx,
				})
			} else if chapter.Title == "" || chapter.Title != bp.Title {
				tx.Model(&chapter).Update("title", bp.Title)
			}
		}
		return nil
	})

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to save blueprints: "+err.Error())
		return
	}

	logger.Info("GenerateBatch: generated %d blueprints for Book=%d (chapters %d-%d)",
		len(savedBlueprints), bookID, req.StartChapter, req.StartChapter+len(savedBlueprints)-1)

	response.Success(c, gin.H{
		"blueprints": savedBlueprints,
	})
}

// RegenerateChapter 重新生成单章蓝图
func (h *OutlineHandler) RegenerateChapter(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid book id")
		return
	}

	chapterIndex, err := strconv.Atoi(c.Param("chapterIndex"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid chapter index")
		return
	}

	// 临时修改 req 使其只生成一章
	c2 := c.Copy()
	c2.Set("_override_start", chapterIndex)
	c2.Set("_override_batch", 1)

	// 加载上下文
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

	// 获取相邻章节作为上下文
	var neighbors []models.ChapterBlueprint
	h.db.Where("book_id = ? AND chapter_index BETWEEN ? AND ?", bookID, chapterIndex-2, chapterIndex+2).
		Where("chapter_index != ?", chapterIndex).
		Order("chapter_index asc").Find(&neighbors)

	var neighborContext string
	for _, bp := range neighbors {
		neighborContext += fmt.Sprintf("第%d章 %s: %s\n", bp.ChapterIndex, bp.Title, bp.Summary)
	}

	ctx := core.WithBookID(c.Request.Context(), uint(bookID))

	promptText := fmt.Sprintf(`你是一个专业的小说章节规划师。请为第%d章重新生成详细的章节蓝图。

## 书籍信息
- 标题: %s
- 类型: %s

## 世界观设定
%s

## 剧情大纲
%s

## 角色设定
%s

## 相邻章节参考
%s

请以纯 JSON 对象格式返回，包含以下字段：
- title: 章节标题
- summary: 章节大纲
- protagonist_motivation: 主角动机
- key_foreshadowing: 关键伏笔
- appearing_characters: 出场人物
- highlight: 审视亮点
- core_events: 核心事件
- challenges: 面临挑战

仅返回 JSON。`,
		chapterIndex,
		book.Title, book.Genre,
		plan.WorldView, plan.Outline, plan.Characters,
		neighborContext,
	)

	messages := []core.Message{
		{Role: core.RoleUser, Content: promptText},
	}

	options := core.Options{
		Model:    "",
		JSONMode: true,
	}
	core.GetStrategy(core.TaskPlanning).ApplyToOptions(&options)

	resp, err := h.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "LLM call failed: "+err.Error())
		return
	}

	var result struct {
		Title                 string `json:"title"`
		Summary               string `json:"summary"`
		ProtagonistMotivation string `json:"protagonist_motivation"`
		KeyForeshadowing      string `json:"key_foreshadowing"`
		AppearingCharacters   string `json:"appearing_characters"`
		Highlight             string `json:"highlight"`
		CoreEvents            string `json:"core_events"`
		Challenges            string `json:"challenges"`
	}

	cleanJSON := core.ParseJSON(resp.Content)
	if err := json.Unmarshal([]byte(cleanJSON), &result); err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to parse blueprint JSON: %v", err))
		return
	}

	// Upsert
	var existing models.ChapterBlueprint
	updates := map[string]interface{}{
		"title":                  result.Title,
		"summary":                result.Summary,
		"protagonist_motivation": result.ProtagonistMotivation,
		"key_foreshadowing":      result.KeyForeshadowing,
		"appearing_characters":   result.AppearingCharacters,
		"highlight":              result.Highlight,
		"core_events":            result.CoreEvents,
		"challenges":             result.Challenges,
	}

	if err := h.db.Where("book_id = ? AND chapter_index = ?", bookID, chapterIndex).First(&existing).Error; err == nil {
		h.db.Model(&existing).Updates(updates)
	} else {
		bp := models.ChapterBlueprint{
			BookID:                uint(bookID),
			ChapterIndex:          chapterIndex,
			Title:                 result.Title,
			Summary:               result.Summary,
			ProtagonistMotivation: result.ProtagonistMotivation,
			KeyForeshadowing:      result.KeyForeshadowing,
			AppearingCharacters:   result.AppearingCharacters,
			Highlight:             result.Highlight,
			CoreEvents:            result.CoreEvents,
			Challenges:            result.Challenges,
		}
		h.db.Create(&bp)
	}

	// 同步标题到 Chapter 表
	h.db.Model(&models.Chapter{}).Where("book_id = ? AND \"order\" = ?", bookID, chapterIndex).Update("title", result.Title)

	response.Success(c, gin.H{
		"chapter_index": chapterIndex,
		"blueprint":     result,
	})
}
