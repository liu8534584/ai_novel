package handler

import (
	"context"
	"net/http"
	"strconv"

	"ai_novel/internal/service/agent"
	"ai_novel/internal/service/llm"
	"ai_novel/internal/service/llm/core"
	"ai_novel/models"
	"ai_novel/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookHandler struct {
	db          *gorm.DB
	llmProvider *llm.DynamicProvider
}

func NewBookHandler(db *gorm.DB, llmProvider *llm.DynamicProvider) *BookHandler {
	return &BookHandler{
		db:          db,
		llmProvider: llmProvider,
	}
}

func (h *BookHandler) ChatInspiration(c *gin.Context) {
	var req struct {
		Messages []core.Message `json:"messages" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	director := agent.NewDirectorAgent(h.llmProvider, h.db)
	resp, err := director.ChatForInspiration(context.Background(), req.Messages)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, resp)
}

func (h *BookHandler) FinalizeInspiration(c *gin.Context) {
	var req struct {
		Conversation string `json:"conversation" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	director := agent.NewDirectorAgent(h.llmProvider, h.db)
	inspiration, err := director.FinalizeInspiration(context.Background(), req.Conversation)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, inspiration)
}

func (h *BookHandler) GetInspirationChat(c *gin.Context) {
	var chat models.InspirationChat
	if err := h.db.Order("created_at desc").First(&chat).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Success(c, nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, chat.Messages)
}

func (h *BookHandler) SaveInspirationChat(c *gin.Context) {
	var req struct {
		Messages string `json:"messages" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var chat models.InspirationChat
	if err := h.db.Order("created_at desc").First(&chat).Error; err == nil {
		chat.Messages = req.Messages
		h.db.Save(&chat)
	} else {
		chat = models.InspirationChat{Messages: req.Messages}
		h.db.Create(&chat)
	}

	response.Success(c, nil)
}

func (h *BookHandler) ListBooks(c *gin.Context) {
	var books []models.Book
	if err := h.db.Order("created_at desc").Find(&books).Error; err != nil {
		response.ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, "Failed to fetch books")
		return
	}
	response.Success(c, books)
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var req struct {
		Title          string                `json:"title" binding:"required"`
		Author         string                `json:"author"`
		Genre          string                `json:"genre"`
		Tags           string                `json:"tags"`
		Language       string                `json:"language"`
		Description    string                `json:"description"`
		TotalChapters  int                   `json:"total_chapters"`
		Status         string                `json:"status"`
		LLMConfig      models.BookLLMConfig  `json:"llm_config"`
		PromptBindings models.PromptBindings `json:"prompt_bindings"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.TotalChapters <= 0 {
		req.TotalChapters = 1
	}
	if req.Status == "" {
		req.Status = "draft"
	}
	if isPromptBindingsEmpty(req.PromptBindings) {
		req.PromptBindings = models.PromptBindings{
			WorldView:      []string{"director"},
			Plan:           []string{"planner"},
			Character:      []string{"character"},
			ChapterTitle:   []string{"chapter_title"},
			ChapterOutline: []string{"outliner"},
			Writing:        []string{"writer_layered"},
			Review:         []string{"state_audit", "event_extraction", "foreshadowing_resolution", "ooc_evaluation", "contradiction_detection"},
		}
	}

	book := models.Book{
		Title:          req.Title,
		Author:         req.Author,
		Genre:          req.Genre,
		Tags:           req.Tags,
		Language:       req.Language,
		Description:    req.Description,
		TotalChapters:  req.TotalChapters,
		Status:         req.Status,
		LLMConfig:      req.LLMConfig,
		PromptBindings: req.PromptBindings,
	}

	if err := h.db.Create(&book).Error; err != nil {
		response.ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, "Failed to create book")
		return
	}
	response.Success(c, book)
}

func isPromptBindingsEmpty(bindings models.PromptBindings) bool {
	return len(bindings.WorldView) == 0 &&
		len(bindings.Plan) == 0 &&
		len(bindings.Character) == 0 &&
		len(bindings.ChapterTitle) == 0 &&
		len(bindings.ChapterOutline) == 0 &&
		len(bindings.Writing) == 0 &&
		len(bindings.Review) == 0
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Title          string                `json:"title"`
		Author         string                `json:"author"`
		Genre          string                `json:"genre"`
		Tags           string                `json:"tags"`
		Language       string                `json:"language"`
		Description    string                `json:"description"`
		TotalChapters  int                   `json:"total_chapters"`
		Status         string                `json:"status"`
		LLMConfig      models.BookLLMConfig  `json:"llm_config"`
		PromptBindings models.PromptBindings `json:"prompt_bindings"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	updates := models.Book{}
	selectFields := make([]string, 0, 10)
	if req.Title != "" {
		updates.Title = req.Title
		selectFields = append(selectFields, "title")
	}
	if req.Author != "" {
		updates.Author = req.Author
		selectFields = append(selectFields, "author")
	}
	if req.Genre != "" {
		updates.Genre = req.Genre
		selectFields = append(selectFields, "genre")
	}
	if req.Tags != "" {
		updates.Tags = req.Tags
		selectFields = append(selectFields, "tags")
	}
	if req.Language != "" {
		updates.Language = req.Language
		selectFields = append(selectFields, "language")
	}
	if req.Description != "" {
		updates.Description = req.Description
		selectFields = append(selectFields, "description")
	}
	if req.TotalChapters > 0 {
		updates.TotalChapters = req.TotalChapters
		selectFields = append(selectFields, "total_chapters")
	}
	if req.Status != "" {
		updates.Status = req.Status
		selectFields = append(selectFields, "status")
	}
	if req.LLMConfig.Provider != "" || req.LLMConfig.APIKey != "" || req.LLMConfig.BaseURL != "" || req.LLMConfig.Model != "" {
		updates.LLMConfig = req.LLMConfig
		selectFields = append(selectFields, "llm_config")
	}
	if len(req.PromptBindings.WorldView) > 0 ||
		len(req.PromptBindings.Plan) > 0 ||
		len(req.PromptBindings.Character) > 0 ||
		len(req.PromptBindings.ChapterTitle) > 0 ||
		len(req.PromptBindings.ChapterOutline) > 0 ||
		len(req.PromptBindings.Writing) > 0 ||
		len(req.PromptBindings.Review) > 0 {
		updates.PromptBindings = req.PromptBindings
		selectFields = append(selectFields, "prompt_bindings")
	}

	var book models.Book
	if err := h.db.First(&book, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.ErrorWithStatus(c, http.StatusNotFound, http.StatusNotFound, "Book not found")
			return
		}
		response.ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, "Failed to fetch book")
		return
	}

	if len(selectFields) == 0 {
		response.Success(c, book)
		return
	}

	if err := h.db.Model(&book).Select(selectFields).Updates(updates).Error; err != nil {
		response.ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, "Failed to update book")
		return
	}

	// Clear provider cache if config was updated
	bookIDUint, _ := strconv.ParseUint(id, 10, 32)
	h.llmProvider.ClearCache(uint(bookIDUint))

	if err := h.db.First(&book, id).Error; err != nil {
		response.ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, "Failed to load book")
		return
	}
	response.Success(c, book)
}

func (h *BookHandler) GetBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	if err := h.db.First(&book, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.ErrorWithStatus(c, http.StatusNotFound, http.StatusNotFound, "Book not found")
			return
		}
		response.ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, "Failed to fetch book")
		return
	}
	response.Success(c, book)
}

func (h *BookHandler) ListChapters(c *gin.Context) {
	bookID := c.Param("id")
	var chapters []models.Chapter
	if err := h.db.Where("book_id = ?", bookID).Order("`order` asc").Find(&chapters).Error; err != nil {
		response.ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, "Failed to fetch chapters")
		return
	}
	response.Success(c, chapters)
}

func (h *BookHandler) ListCharacters(c *gin.Context) {
	bookID := c.Param("id")
	var characters []models.Character
	if err := h.db.Where("book_id = ?", bookID).Find(&characters).Error; err != nil {
		response.ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, "Failed to fetch characters")
		return
	}
	response.Success(c, characters)
}

func (h *BookHandler) UpdateChapter(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Outline string `json:"outline"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Outline != "" {
		updates["outline"] = req.Outline
	}

	if err := h.db.Model(&models.Chapter{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		response.ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, "Failed to update chapter")
		return
	}
	response.Success(c, gin.H{"id": id})
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" || idStr == "undefined" || idStr == "null" {
		response.Error(c, http.StatusBadRequest, "Invalid book ID")
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		// 删除关联的章节
		if err := tx.Where("book_id = ?", idStr).Delete(&models.Chapter{}).Error; err != nil {
			return err
		}
		// 删除关联的角色
		if err := tx.Where("book_id = ?", idStr).Delete(&models.Character{}).Error; err != nil {
			return err
		}
		// 删除书籍本身
		if err := tx.Where("id = ?", idStr).Delete(&models.Book{}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		response.ErrorWithStatus(c, http.StatusInternalServerError, http.StatusInternalServerError, "Failed to delete book: "+err.Error())
		return
	}
	response.Success(c, gin.H{"id": idStr})
}
