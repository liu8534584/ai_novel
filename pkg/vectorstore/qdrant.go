package vectorstore

import (
	"context"
	"fmt"

	qdrant "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type QdrantStore struct {
	conn *grpc.ClientConn
}

func NewQdrantStore(address string, apiKey string) (*QdrantStore, error) {
	// For simplicity, we assume address is host:port
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial qdrant: %w", err)
	}

	return &QdrantStore{conn: conn}, nil
}

func (s *QdrantStore) CreateCollection(ctx context.Context, collectionName string, dimension int) error {
	collectionsClient := qdrant.NewCollectionsClient(s.conn)
	
	// Check if exists (this is a simplified check)
	_, err := collectionsClient.Get(ctx, &qdrant.GetCollectionInfoRequest{
		CollectionName: collectionName,
	})
	if err == nil {
		return nil // Already exists
	}

	_, err = collectionsClient.Create(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: &qdrant.VectorsConfig{
			Config: &qdrant.VectorsConfig_Params{
				Params: &qdrant.VectorParams{
					Size:     uint64(dimension),
					Distance: qdrant.Distance_Cosine,
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	return nil
}

func (s *QdrantStore) AddDocuments(ctx context.Context, collectionName string, docs []Document) error {
	pointsClient := qdrant.NewPointsClient(s.conn)

	points := make([]*qdrant.PointStruct, len(docs))
	for i, doc := range docs {
		payload := make(map[string]*qdrant.Value)
		for k, v := range doc.Metadata {
			payload[k] = &qdrant.Value{
				Kind: &qdrant.Value_StringValue{StringValue: fmt.Sprint(v)},
			}
		}
		// Also store content in payload
		payload["content"] = &qdrant.Value{
			Kind: &qdrant.Value_StringValue{StringValue: doc.Content},
		}

		points[i] = &qdrant.PointStruct{
			Id: &qdrant.PointId{
				PointIdOptions: &qdrant.PointId_Uuid{Uuid: doc.ID},
			},
			Vectors: &qdrant.Vectors{
				VectorsOptions: &qdrant.Vectors_Vector{
					Vector: &qdrant.Vector{
						Data: doc.Vector,
					},
				},
			},
			Payload: payload,
		}
	}

	_, err := pointsClient.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points:         points,
	})
	if err != nil {
		return fmt.Errorf("failed to upsert points: %w", err)
	}

	return nil
}

func (s *QdrantStore) Search(ctx context.Context, collectionName string, vector []float32, topK int) ([]Document, error) {
	pointsClient := qdrant.NewPointsClient(s.conn)

	resp, err := pointsClient.Search(ctx, &qdrant.SearchPoints{
		CollectionName: collectionName,
		Vector:         vector,
		Limit:          uint64(topK),
		WithPayload:    &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true}},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search points: %w", err)
	}

	docs := make([]Document, len(resp.Result))
	for i, hit := range resp.Result {
		doc := Document{
			ID:       hit.Id.GetUuid(),
			Score:    hit.Score,
			Metadata: make(map[string]interface{}),
		}

		for k, v := range hit.Payload {
			if k == "content" {
				doc.Content = v.GetStringValue()
			} else {
				doc.Metadata[k] = v.GetStringValue()
			}
		}
		docs[i] = doc
	}

	return docs, nil
}

func (s *QdrantStore) DeleteCollection(ctx context.Context, collectionName string) error {
	collectionsClient := qdrant.NewCollectionsClient(s.conn)
	_, err := collectionsClient.Delete(ctx, &qdrant.DeleteCollection{
		CollectionName: collectionName,
	})
	return err
}

func (s *QdrantStore) DeleteDocuments(ctx context.Context, collectionName string, filter map[string]interface{}) error {
	if len(filter) == 0 {
		return nil
	}

	pointsClient := qdrant.NewPointsClient(s.conn)
	
	var conditions []*qdrant.Condition
	for k, v := range filter {
		conditions = append(conditions, &qdrant.Condition{
			ConditionOneOf: &qdrant.Condition_Field{
				Field: &qdrant.FieldCondition{
					Key: k,
					Match: &qdrant.Match{
						MatchValue: &qdrant.Match_Text{Text: fmt.Sprint(v)},
					},
				},
			},
		})
	}

	_, err := pointsClient.Delete(ctx, &qdrant.DeletePoints{
		CollectionName: collectionName,
		Points: &qdrant.PointsSelector{
			PointsSelectorOneOf: &qdrant.PointsSelector_Filter{
				Filter: &qdrant.Filter{
					Must: conditions,
				},
			},
		},
	})
	return err
}
