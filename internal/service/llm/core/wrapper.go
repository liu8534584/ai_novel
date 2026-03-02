package core

import (
	"context"
	"fmt"
	"time"
)

// Wrapper provides a base for LLM provider decorators
type Wrapper struct {
	Provider Provider
}

// Chat implements the Provider interface with retry logic
func (w *Wrapper) Chat(ctx context.Context, messages []Message, options Options) (Response, error) {
	var lastErr error
	maxRetries := 3
	if options.MaxRetries > 0 {
		maxRetries = options.MaxRetries
	}

	for i := 0; i < maxRetries; i++ {
		resp, err := w.Provider.Chat(ctx, messages, options)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		// Exponential backoff
		time.Sleep(time.Duration(i+1) * 500 * time.Millisecond)
	}

	return Response{}, fmt.Errorf("all retries failed: %w", lastErr)
}

// StreamChat implements the Provider interface
func (w *Wrapper) StreamChat(ctx context.Context, messages []Message, options Options) (<-chan StreamResponse, error) {
	// For streaming, retrying mid-stream is complex, so we just pass through for now
	// or retry the initial connection.
	return w.Provider.StreamChat(ctx, messages, options)
}

// CreateEmbedding implements the Provider interface with retry logic
func (w *Wrapper) CreateEmbedding(ctx context.Context, input string, options Options) ([]float32, error) {
	var lastErr error
	for i := 0; i < 3; i++ {
		resp, err := w.Provider.CreateEmbedding(ctx, input, options)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		time.Sleep(time.Duration(i+1) * 500 * time.Millisecond)
	}
	return nil, fmt.Errorf("all retries failed: %w", lastErr)
}

// MetricsWrapper tracks token usage
type MetricsWrapper struct {
	Provider Provider
	OnUsage  func(Usage)
}

func (w *MetricsWrapper) Chat(ctx context.Context, messages []Message, options Options) (Response, error) {
	resp, err := w.Provider.Chat(ctx, messages, options)
	if err == nil && w.OnUsage != nil {
		w.OnUsage(resp.Usage)
	}
	return resp, err
}

func (w *MetricsWrapper) StreamChat(ctx context.Context, messages []Message, options Options) (<-chan StreamResponse, error) {
	return w.Provider.StreamChat(ctx, messages, options)
}

func (w *MetricsWrapper) CreateEmbedding(ctx context.Context, input string, options Options) ([]float32, error) {
	return w.Provider.CreateEmbedding(ctx, input, options)
}
