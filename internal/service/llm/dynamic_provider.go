package llm

import (
	"ai_novel/internal/config"
	"ai_novel/internal/service/llm/core"
	"ai_novel/models"
	"ai_novel/pkg/logger"
	"context"
	"encoding/json"
	"sync"
	"time"

	"gorm.io/gorm"
)

const (
	DefaultLLMTimeout = 10 * time.Minute
)

type DynamicProvider struct {
	db      *gorm.DB
	cache   map[uint]core.Provider
	cacheMu sync.RWMutex
	global  core.Provider
}

func NewDynamicProvider(db *gorm.DB, global core.Provider) *DynamicProvider {
	return &DynamicProvider{
		db:    db,
		cache: make(map[uint]core.Provider),
		global: global,
	}
}

func (p *DynamicProvider) resolveProvider(ctx context.Context, options *core.Options) core.Provider {
	bookID, ok := ctx.Value(core.ContextKeyBookID).(uint)
	if !ok || bookID == 0 {
		if options.Model == "" {
			options.Model = config.GlobalConfig.LLM.Model
		}
		return p.global
	}

	p.cacheMu.RLock()
	provider, exists := p.cache[bookID]
	p.cacheMu.RUnlock()

	if exists {
		// Even if cached, we might need to update the model in options
		p.updateOptions(bookID, options)
		return provider
	}

	// Load from DB
	var book models.Book
	if err := p.db.Select("llm_config").First(&book, bookID).Error; err != nil {
		if options.Model == "" {
			options.Model = config.GlobalConfig.LLM.Model
		}
		return p.global
	}

	if book.LLMConfig.APIKey != "" {
		provider = NewProvider(book.LLMConfig.Provider, book.LLMConfig.APIKey, book.LLMConfig.BaseURL)
		p.cacheMu.Lock()
		p.cache[bookID] = provider
		p.cacheMu.Unlock()
		
		p.updateOptionsFromBook(&book, options)
		return provider
	}

	if options.Model == "" {
		options.Model = config.GlobalConfig.LLM.Model
	}
	return p.global
}

func (p *DynamicProvider) updateOptions(bookID uint, options *core.Options) {
	if options.Model != "" {
		return
	}
	var book models.Book
	if err := p.db.Select("llm_config").First(&book, bookID).Error; err == nil {
		p.updateOptionsFromBook(&book, options)
	}
}

func (p *DynamicProvider) updateOptionsFromBook(book *models.Book, options *core.Options) {
	if options.Model == "" {
		if book.LLMConfig.Model != "" {
			options.Model = book.LLMConfig.Model
		} else {
			options.Model = config.GlobalConfig.LLM.Model
		}
	}
}

func (p *DynamicProvider) Chat(ctx context.Context, messages []core.Message, options core.Options) (core.Response, error) {
	// Create a timeout context based on the provided one
	chatCtx, cancel := context.WithTimeout(ctx, DefaultLLMTimeout)
	defer cancel()

	provider := p.resolveProvider(chatCtx, &options)
	resp, err := provider.Chat(chatCtx, messages, options)

	// Log request and response
	bookID, _ := ctx.Value(core.ContextKeyBookID).(uint)
	msgJSON, _ := json.MarshalIndent(messages, "", "  ")
	respJSON, _ := json.MarshalIndent(resp, "", "  ")
	logger.LogLLMRequest(bookID, options.Model, string(msgJSON), string(respJSON), err)

	return resp, err
}

func (p *DynamicProvider) StreamChat(ctx context.Context, messages []core.Message, options core.Options) (<-chan core.StreamResponse, error) {
	// For streaming, we also apply a timeout to the initial connection and the overall stream
	streamCtx, cancel := context.WithTimeout(ctx, DefaultLLMTimeout)
	// Note: We don't defer cancel() here because the stream needs the context to stay alive.
	// The provider or the caller should manage the lifecycle.
	// Actually, in many implementations, the stream is tied to the context.

	provider := p.resolveProvider(streamCtx, &options)
	
	// Log request
	bookID, _ := ctx.Value(core.ContextKeyBookID).(uint)
	msgJSON, _ := json.MarshalIndent(messages, "", "  ")
	logger.LLM("BookID: %d | Model: %s | Status: STREAM_STARTED\n--- [REQUEST] ---\n%s", 
		bookID, options.Model, string(msgJSON))

	respChan, err := provider.StreamChat(streamCtx, messages, options)
	if err != nil {
		cancel()
		return nil, err
	}

	// Create a wrapper channel to ensure cancel() is called when the stream ends
	wrappedChan := make(chan core.StreamResponse)
	go func() {
		defer cancel()
		defer close(wrappedChan)
		for resp := range respChan {
			wrappedChan <- resp
		}
	}()

	return wrappedChan, nil
}

func (p *DynamicProvider) CreateEmbedding(ctx context.Context, input string, options core.Options) ([]float32, error) {
	provider := p.resolveProvider(ctx, &options)
	
	// If model is still empty, use a default embedding model if configured
	if options.Model == "" || options.Model == config.GlobalConfig.LLM.Model {
		if config.GlobalConfig.LLM.EmbeddingModel != "" {
			options.Model = config.GlobalConfig.LLM.EmbeddingModel
		}
	}

	return provider.CreateEmbedding(ctx, input, options)
}

// ClearCache allows clearing the provider cache when a book's config is updated.
func (p *DynamicProvider) ClearCache(bookID uint) {
	p.cacheMu.Lock()
	delete(p.cache, bookID)
	p.cacheMu.Unlock()
}
