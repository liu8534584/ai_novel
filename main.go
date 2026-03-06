package main

import (
	"ai_novel/internal/config"
	"ai_novel/internal/handler"
	"ai_novel/internal/service"
	"ai_novel/internal/service/agent"
	svccontext "ai_novel/internal/service/context"
	"ai_novel/internal/service/llm"
	"ai_novel/internal/service/rag"
	"ai_novel/internal/static"
	"ai_novel/models"
	"ai_novel/pkg/logger"
	"ai_novel/pkg/vectorstore"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 0. Load Configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	// Initialize Logger
	if err := logger.Init(config.GlobalConfig.Log.Level, config.GlobalConfig.Log.Filename); err != nil {
		log.Printf("Warning: Failed to initialize file logger: %v", err)
	}
	defer logger.Close()

	// 1. Initialize Database
	dbConfig := config.GlobalConfig.Database
	var dialector gorm.Dialector
	switch dbConfig.Driver {
	case "mysql":
		dialector = mysql.Open(dbConfig.Source)
	case "postgres", "pgsql", "postgresql":
		dialector = postgres.Open(dbConfig.Source)
	case "sqlite":
		dialector = sqlite.Open(dbConfig.Source)
	default:
		log.Printf("Warning: Unknown driver %s, falling back to sqlite", dbConfig.Driver)
		dialector = sqlite.Open("novel.db")
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	// Auto Migrate the schema
	err = db.AutoMigrate(&models.Book{}, &models.Chapter{}, &models.ChapterVersion{}, &models.VectorRecord{}, &models.Character{}, &models.CharacterStateRecord{}, &models.StoryEvent{}, &models.OutlineVersion{}, &models.Foreshadowing{}, &models.CharacterAnchor{}, &models.OOCScore{}, &models.StoryContradiction{}, &models.ChapterHealthScore{}, &models.StoryArc{}, &models.ChapterBlueprint{})
	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}
	logger.Info("Database migration completed successfully.")

	// 2. Initialize LLM Provider
	llmConfig := config.GlobalConfig.LLM
	baseProvider := llm.NewProvider(llmConfig.Provider, llmConfig.APIKey, llmConfig.BaseURL)
	llmProvider := llm.NewDynamicProvider(db, baseProvider)

	// 3. Initialize Agents
	directorAgent := agent.NewDirectorAgent(llmProvider, db)
	outlinerAgent := agent.NewOutlinerAgent(llmProvider)
	writerAgent := agent.NewWriterAgent(llmProvider)
	stateAgent := agent.NewStateAgent(llmProvider)
	foresightAgent := agent.NewForesightAgent(db, llmProvider)
	summarizerAgent := agent.NewSummarizerAgent(llmProvider)
	characterAgent := agent.NewCharacterAgent(llmProvider)
	chapterTitleAgent := agent.NewChapterTitleAgent(llmProvider)
	consistencyAgent := agent.NewConsistencyAgent(llmProvider)

	// Initialize Vector Store
	var vStore vectorstore.VectorStore
	if config.GlobalConfig.Vector.Provider != "" && config.GlobalConfig.Vector.Provider != "sqlite" {
		var err error
		vStore, err = vectorstore.NewVectorStore(config.GlobalConfig.Vector)
		if err != nil {
			log.Printf("Warning: Failed to initialize vector store: %v", err)
		} else {
			// Initialize collections
			ctx := context.Background()
			dim := 1536 // Default for OpenAI/GLM, should ideally be configurable
			vStore.CreateCollection(ctx, rag.CollectionCharacters, dim)
			vStore.CreateCollection(ctx, rag.CollectionHistory, dim)
			vStore.CreateCollection(ctx, rag.CollectionOutlines, dim)
			vStore.CreateCollection(ctx, rag.CollectionWorldRules, dim)
		}
	}

	ragService := rag.NewMemoryRecallService(db, llmProvider, vStore)
	ctxMgr := svccontext.NewContextManager(db, summarizerAgent, writerAgent, ragService)

	// Post-write processor (统一后处理器)
	processor := &service.PostWriteProcessor{
		DB:          db,
		State:       stateAgent,
		Foresight:   foresightAgent,
		Consistency: consistencyAgent,
		Summarizer:  summarizerAgent,
		RAG:         ragService,
	}

	// 4. Initialize Handler
	novelHandler := handler.NewNovelHandler(db, directorAgent, outlinerAgent, writerAgent, stateAgent, foresightAgent, consistencyAgent, summarizerAgent, ctxMgr, ragService, processor)
	bookHandler := handler.NewBookHandler(db, llmProvider)
	planAgent := agent.NewPlanAgent(llmProvider)
	planHandler := handler.NewPlanHandler(db, planAgent, directorAgent, characterAgent, chapterTitleAgent, ragService)
	configHandler := handler.NewConfigHandler()
	outlineHandler := handler.NewOutlineHandler(db, llmProvider)

	// 5. Initialize Router
	gin.SetMode(config.GlobalConfig.Server.Mode)
	r := gin.New() // Use gin.New() to have full control over middlewares
	r.Use(gin.Recovery())

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Register Routes
	handler.RegisterRoutes(r, novelHandler, bookHandler, planHandler, configHandler, outlineHandler)

	// Root redirect to /ui
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/ui/")
	})

	// Static files from embed
	staticFS := static.GetStaticFS()
	staticServer := http.FileServer(http.FS(staticFS))

	// Serve UI at /ui/
	uiGroup := r.Group("/ui")
	{
		uiGroup.GET("/*any", func(c *gin.Context) {
			// Strip /ui prefix for the file server
			path := c.Param("any")

			// Check if file exists in staticFS
			// Note: fs.FS uses forward slashes and no leading slash usually
			trimPath := strings.TrimPrefix(path, "/")

			// If empty, serve index.html
			if trimPath == "" {
				c.Request.URL.Path = "/"
				staticServer.ServeHTTP(c.Writer, c.Request)
				return
			}

			_, err := fs.Stat(staticFS, trimPath)
			if errors.Is(err, fs.ErrNotExist) {
				// File not found, serve index.html for SPA
				c.Request.URL.Path = "/"
			} else {
				c.Request.URL.Path = path
			}
			staticServer.ServeHTTP(c.Writer, c.Request)
		})
	}

	// Handle SPA routes in NoRoute (fallback)
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// If it's an API request that wasn't found, let it be 404
		if len(path) >= 5 && path[:5] == "/api/" {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "API not found"})
			return
		}
		// Redirect other unknown routes to /ui/
		c.Redirect(http.StatusFound, "/ui/")
	})

	addr := fmt.Sprintf(":%s", config.GlobalConfig.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("failed to run server: ", err)
	}
}
