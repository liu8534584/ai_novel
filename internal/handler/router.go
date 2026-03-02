package handler

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *NovelHandler, bookHandler *BookHandler, planHandler *PlanHandler, configHandler *ConfigHandler) {
	api := r.Group("/api")
	{
		api.GET("/books", bookHandler.ListBooks)
		api.POST("/books", bookHandler.CreateBook)
		api.POST("/books/inspiration/chat", bookHandler.ChatInspiration)
		api.POST("/books/inspiration/finalize", bookHandler.FinalizeInspiration)
		api.GET("/books/inspiration/chat", bookHandler.GetInspirationChat)
		api.POST("/books/inspiration/chat/save", bookHandler.SaveInspirationChat)
		api.GET("/books/:id", bookHandler.GetBook)
		api.PUT("/books/:id", bookHandler.UpdateBook)
		api.DELETE("/books/:id", bookHandler.DeleteBook)
		api.GET("/books/:id/chapters", bookHandler.ListChapters)
		api.GET("/books/:id/characters", bookHandler.ListCharacters)
		api.PUT("/chapters/:id", bookHandler.UpdateChapter)

		api.POST("/books/init", h.CreateBook)
		api.GET("/books/:id/state", h.GetState)

		api.GET("/books/:id/plans", planHandler.ListPlans)
		api.POST("/books/:id/plans/generate", planHandler.GeneratePlans)
		api.PUT("/books/:id/plans/:planId/select", planHandler.SelectPlan)
		api.PUT("/books/:id/plans/:planId/lock", planHandler.LockPlan)
		api.PUT("/books/:id/plans/:planId", planHandler.UpdatePlan)
		api.POST("/books/:id/plans/characters", planHandler.GenerateCharacters)
		api.POST("/books/:id/plans/chapters", planHandler.GenerateChapterTitles)
		api.POST("/books/:id/plans/chapters/batch", planHandler.GenerateChapterTitlesBatch)

		api.POST("/chapters/:id/outline", h.GenerateOutline)
		api.POST("/chapters/:id/outline/confirm", h.ConfirmOutline)
		api.POST("/chapters/:id/write", h.WriteChapter)
		api.GET("/chapters/:id/versions", h.ListChapterVersions)
		api.PUT("/chapters/:id/versions/:versionId/rollback", h.RollbackChapter)
		api.GET("/chapters/:id/ooc-scores", h.GetOOCScores)
		api.GET("/chapters/:id/contradictions", h.GetContradictions)
		api.GET("/chapters/:id/health", h.GetChapterHealth)
		api.GET("/books/:id/foreshadowing-alerts", h.GetForeshadowingAlerts)

		// Config Routes
		api.GET("/config/llm", configHandler.GetLLMConfig)
		api.PUT("/config/llm", configHandler.UpdateLLMConfig)
		api.POST("/config/llm/test", configHandler.TestLLMConnection)
	}
}
