package handler

import (
	"net/http"

	"ai_novel/internal/config"
	"ai_novel/internal/service/llm"
	"ai_novel/internal/service/llm/core"
	"ai_novel/pkg/response"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct{}

func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{}
}

func (h *ConfigHandler) GetLLMConfig(c *gin.Context) {
	response.Success(c, config.GlobalConfig.LLM)
}

func (h *ConfigHandler) TestLLMConnection(c *gin.Context) {
	var req config.LLMConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Create temporary provider for testing
	provider := llm.NewProvider(req.Provider, req.APIKey, req.BaseURL)
	
	// Try a simple chat request
	ctx := c.Request.Context()
	messages := []core.Message{
		{Role: "user", Content: "Hello, this is a connection test. Please reply with 'OK'."},
	}
	
	options := core.Options{
		Model:       req.Model,
		Temperature: 0.1,
		MaxTokens:   10,
	}

	resp, err := provider.Chat(ctx, messages, options)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Connection failed: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "Connection successful",
		"reply":   resp.Content,
	})
}

func (h *ConfigHandler) UpdateLLMConfig(c *gin.Context) {
	var req config.LLMConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update Global Config
	config.GlobalConfig.LLM = req

	// Save to file
	if err := config.SaveConfig(); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to save config: "+err.Error())
		return
	}

	response.Success(c, config.GlobalConfig.LLM)
}
