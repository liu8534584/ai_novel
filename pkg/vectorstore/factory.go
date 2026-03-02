package vectorstore

import (
	"ai_novel/internal/config"
	"fmt"
)

func NewVectorStore(cfg config.VectorConfig) (VectorStore, error) {
	switch cfg.Provider {
	case "pgsql", "postgres", "postgresql":
		return NewPgVectorStore(cfg.Address)
	case "qdrant":
		return NewQdrantStore(cfg.Address, cfg.APIKey)
	case "milvus":
		return NewMilvusStore(cfg.Address)
	// Add other providers here
	default:
		return nil, fmt.Errorf("unsupported vector provider: %s", cfg.Provider)
	}
}
