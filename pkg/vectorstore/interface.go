package vectorstore

import (
	"context"
)

// Document represents a document to be stored in the vector store
type Document struct {
	ID        string                 `json:"id"`
	Content   string                 `json:"content"`
	Vector    []float32              `json:"vector"`
	Metadata  map[string]interface{} `json:"metadata"`
	Score     float32                `json:"score"` // Only used for search results
}

// SearchResult represents the result of a similarity search
type SearchResult struct {
	Documents []Document
}

// VectorStore is the interface for vector database operations
type VectorStore interface {
	// AddDocuments adds documents to the vector store
	AddDocuments(ctx context.Context, collectionName string, docs []Document) error
	
	// Search searches for similar documents in the vector store
	Search(ctx context.Context, collectionName string, vector []float32, topK int) ([]Document, error)
	
	// CreateCollection creates a new collection/index if it doesn't exist
	CreateCollection(ctx context.Context, collectionName string, dimension int) error
	
	// DeleteCollection deletes a collection/index
	DeleteCollection(ctx context.Context, collectionName string) error
}
