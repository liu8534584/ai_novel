package core

import "context"

type Provider interface {
	Chat(ctx context.Context, messages []Message, options Options) (Response, error)
	StreamChat(ctx context.Context, messages []Message, options Options) (<-chan StreamResponse, error)
	CreateEmbedding(ctx context.Context, input string, options Options) ([]float32, error)
}
