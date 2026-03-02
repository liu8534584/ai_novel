package vectorstore

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgVectorStore struct {
	pool *pgxpool.Pool
}

var collectionNamePattern = regexp.MustCompile(`[^a-z0-9_]+`)

func NewPgVectorStore(address string) (*PgVectorStore, error) {
	pool, err := pgxpool.New(context.Background(), address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to pgsql: %w", err)
	}
	return &PgVectorStore{pool: pool}, nil
}

func (s *PgVectorStore) AddDocuments(ctx context.Context, collectionName string, docs []Document) error {
	if len(docs) == 0 {
		return nil
	}
	if err := s.CreateCollection(ctx, collectionName, len(docs[0].Vector)); err != nil {
		return err
	}
	table := s.collectionName(collectionName)
	if table == "" {
		return fmt.Errorf("invalid collection name")
	}
	query := fmt.Sprintf("INSERT INTO %s (content, metadata, vector) VALUES ($1, $2::jsonb, $3::vector)", table)
	for _, doc := range docs {
		metadataJSON, err := json.Marshal(doc.Metadata)
		if err != nil {
			return err
		}
		if _, err := s.pool.Exec(ctx, query, doc.Content, string(metadataJSON), vectorLiteral(doc.Vector)); err != nil {
			return err
		}
	}
	return nil
}

func (s *PgVectorStore) Search(ctx context.Context, collectionName string, vector []float32, topK int) ([]Document, error) {
	if topK <= 0 {
		return []Document{}, nil
	}
	table := s.collectionName(collectionName)
	if table == "" {
		return []Document{}, nil
	}
	query := fmt.Sprintf("SELECT id, content, metadata, vector <-> $1::vector AS score FROM %s ORDER BY vector <-> $1::vector LIMIT $2", table)
	rows, err := s.pool.Query(ctx, query, vectorLiteral(vector), topK)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return []Document{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	var results []Document
	for rows.Next() {
		var id int64
		var content string
		var metadataRaw []byte
		var score float64
		if err := rows.Scan(&id, &content, &metadataRaw, &score); err != nil {
			return nil, err
		}
		metadata := map[string]interface{}{}
		if len(metadataRaw) > 0 {
			if err := json.Unmarshal(metadataRaw, &metadata); err != nil {
				return nil, err
			}
		}
		results = append(results, Document{
			ID:       strconv.FormatInt(id, 10),
			Content:  content,
			Metadata: metadata,
			Score:    float32(score),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *PgVectorStore) CreateCollection(ctx context.Context, collectionName string, dimension int) error {
	if dimension <= 0 {
		return fmt.Errorf("invalid vector dimension")
	}
	table := s.collectionName(collectionName)
	if table == "" {
		return fmt.Errorf("invalid collection name")
	}
	if _, err := s.pool.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS vector"); err != nil {
		return err
	}
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id bigserial primary key, content text, metadata jsonb, vector vector(%d))", table, dimension)
	_, err := s.pool.Exec(ctx, query)
	return err
}

func (s *PgVectorStore) DeleteCollection(ctx context.Context, collectionName string) error {
	table := s.collectionName(collectionName)
	if table == "" {
		return fmt.Errorf("invalid collection name")
	}
	_, err := s.pool.Exec(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
	return err
}

func (s *PgVectorStore) collectionName(name string) string {
	safe := strings.ToLower(name)
	safe = collectionNamePattern.ReplaceAllString(safe, "_")
	safe = strings.Trim(safe, "_")
	if safe == "" {
		return ""
	}
	if safe[0] >= '0' && safe[0] <= '9' {
		safe = "v_" + safe
	}
	return "vs_" + safe
}

func vectorLiteral(vector []float32) string {
	var builder strings.Builder
	builder.WriteByte('[')
	for i, v := range vector {
		if i > 0 {
			builder.WriteByte(',')
		}
		builder.WriteString(strconv.FormatFloat(float64(v), 'f', -1, 32))
	}
	builder.WriteByte(']')
	return builder.String()
}
