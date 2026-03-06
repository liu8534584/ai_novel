package vectorstore

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type MilvusStore struct {
	client client.Client
}

func NewMilvusStore(address string) (*MilvusStore, error) {
	c, err := client.NewClient(context.Background(), client.Config{
		Address: address,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to milvus: %w", err)
	}
	return &MilvusStore{client: c}, nil
}

func (s *MilvusStore) CreateCollection(ctx context.Context, collectionName string, dimension int) error {
	has, err := s.client.HasCollection(ctx, collectionName)
	if err != nil {
		return err
	}
	if has {
		return nil
	}

	schema := entity.NewSchema().WithName(collectionName).WithDescription("Novel context storage").
		WithField(entity.NewField().WithName("id").WithDataType(entity.FieldTypeInt64).WithIsPrimaryKey(true).WithIsAutoID(true)).
		WithField(entity.NewField().WithName("content").WithDataType(entity.FieldTypeVarChar).WithMaxLength(65535)).
		// Metadata is tricky in Milvus, often handled via separate fields or JSON
		WithField(entity.NewField().WithName("metadata").WithDataType(entity.FieldTypeJSON)).
		WithField(entity.NewField().WithName("vector").WithDataType(entity.FieldTypeFloatVector).WithDim(int64(dimension)))

	err = s.client.CreateCollection(ctx, schema, entity.DefaultShardNumber)
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	// Create index
	idx, err := entity.NewIndexIvfFlat(entity.L2, 1024)
	if err != nil {
		return err
	}
	err = s.client.CreateIndex(ctx, collectionName, "vector", idx, false)
	if err != nil {
		return err
	}

	// Load collection
	return s.client.LoadCollection(ctx, collectionName, false)
}

func (s *MilvusStore) AddDocuments(ctx context.Context, collectionName string, docs []Document) error {
	contents := make([]string, len(docs))
	metadatas := make([][]byte, len(docs))
	vectors := make([][]float32, len(docs))

	for i, doc := range docs {
		contents[i] = doc.Content
		vectors[i] = doc.Vector
		metaJSON, err := json.Marshal(doc.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
		metadatas[i] = metaJSON
	}

	contentColumn := entity.NewColumnVarChar("content", contents)
	metadataColumn := entity.NewColumnJSONBytes("metadata", metadatas)
	vectorColumn := entity.NewColumnFloatVector("vector", int(len(vectors[0])), vectors)

	_, err := s.client.Insert(ctx, collectionName, "", contentColumn, metadataColumn, vectorColumn)
	return err
}

func (s *MilvusStore) Search(ctx context.Context, collectionName string, vector []float32, topK int) ([]Document, error) {
	searchParam, _ := entity.NewIndexIvfFlatSearchParam(10)

	res, err := s.client.Search(ctx, collectionName, nil, "", []string{"content", "metadata"}, []entity.Vector{entity.FloatVector(vector)}, "vector", entity.L2, topK, searchParam)
	if err != nil {
		return nil, err
	}

	var docs []Document
	for _, result := range res {
		for i := 0; i < result.ResultCount; i++ {
			content, _ := result.Fields.GetColumn("content").GetAsString(i)
			metadataBytes, _ := result.Fields.GetColumn("metadata").Get(i)
			
			var metadata map[string]interface{}
			if bytes, ok := metadataBytes.([]byte); ok {
				json.Unmarshal(bytes, &metadata)
			}
			
			docs = append(docs, Document{
				ID:       strconv.FormatInt(result.IDs.(*entity.ColumnInt64).Data()[i], 10),
				Content:  content,
				Metadata: metadata,
				Score:    result.Scores[i],
			})
		}
	}
	return docs, nil
}

func (s *MilvusStore) DeleteCollection(ctx context.Context, collectionName string) error {
	return s.client.DropCollection(ctx, collectionName)
}

func (s *MilvusStore) DeleteDocuments(ctx context.Context, collectionName string, filter map[string]interface{}) error {
	if len(filter) == 0 {
		return nil
	}
	// Build expression: e.g. metadata["chapter_id"] == 123
	// Note: Milvus JSON field querying syntax might vary by version.
	// Assuming standard Milvus 2.3+ JSON support.
	
	var exprs []string
	for k, v := range filter {
		switch val := v.(type) {
		case string:
			exprs = append(exprs, fmt.Sprintf("metadata[\"%s\"] == \"%s\"", k, val))
		case int, int32, int64, uint, uint32, uint64:
			exprs = append(exprs, fmt.Sprintf("metadata[\"%s\"] == %d", k, val))
		case float32, float64:
			exprs = append(exprs, fmt.Sprintf("metadata[\"%s\"] == %f", k, val))
		default:
			// Fallback
			exprs = append(exprs, fmt.Sprintf("metadata[\"%s\"] == \"%v\"", k, val))
		}
	}
	expr := strings.Join(exprs, " && ")
	return s.client.Delete(ctx, collectionName, "", expr)
}
