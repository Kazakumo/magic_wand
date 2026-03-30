package vectordb

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/qdrant/go-client/qdrant"
)

const (
	vectorDim      = 768
	batchSize      = 100
	fieldPath      = "path"
	fieldStatus    = "status"
	fieldFileHash  = "file_hash"
	fieldUpdatedAt = "updated_at"
)

// QdrantStore は VectorStore の Qdrant gRPC 実装。
type QdrantStore struct {
	client     *qdrant.Client
	collection string
}

// NewQdrantStore は Qdrant gRPC クライアントを作成してコレクションを初期化する。
func NewQdrantStore(ctx context.Context, host string, port int, collection string) (*QdrantStore, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:     host,
		Port:     port,
		PoolSize: 1,
	})
	if err != nil {
		return nil, fmt.Errorf("qdrant: new client: %w", err)
	}

	s := &QdrantStore{client: client, collection: collection}
	if err := s.ensureCollection(ctx); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *QdrantStore) ensureCollection(ctx context.Context) error {
	exists, err := s.client.CollectionExists(ctx, s.collection)
	if err != nil {
		return fmt.Errorf("qdrant: collection exists: %w", err)
	}
	if exists {
		return nil
	}

	if err := s.client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: s.collection,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     vectorDim,
			Distance: qdrant.Distance_Cosine,
		}),
	}); err != nil {
		return fmt.Errorf("qdrant: create collection: %w", err)
	}

	// インデックス作成: path, status, updated_at
	for _, field := range []struct {
		name      string
		fieldType qdrant.FieldType
	}{
		{fieldPath, qdrant.FieldType_FieldTypeKeyword},
		{fieldStatus, qdrant.FieldType_FieldTypeKeyword},
		{fieldUpdatedAt, qdrant.FieldType_FieldTypeDatetime},
	} {
		if _, err := s.client.CreateFieldIndex(ctx, &qdrant.CreateFieldIndexCollection{
			CollectionName: s.collection,
			FieldName:      field.name,
			FieldType:      qdrant.PtrOf(field.fieldType),
		}); err != nil {
			slog.Warn("qdrant: create field index failed", "field", field.name, "err", err)
		}
	}
	return nil
}

// chunkIDToUUID は chunk ID (sha256 hex 64文字) を UUID 形式にスライスする。
func chunkIDToUUID(id string) string {
	return id[0:8] + "-" + id[8:12] + "-" + id[12:16] + "-" + id[16:20] + "-" + id[20:32]
}

// TryValueMapFromPayload は Payload を Qdrant の Value マップに変換する。テスト用に公開。
func TryValueMapFromPayload(p *Payload) (map[string]*qdrant.Value, error) {
	vm, err := qdrant.TryValueMap(payloadToMap(p))
	if err != nil {
		return nil, fmt.Errorf("payload to value map: %w", err)
	}
	return vm, nil
}

func payloadToMap(p *Payload) map[string]any {
	return map[string]any{
		"path":          p.Path,
		"filename":      p.Filename,
		"ext":           p.Ext,
		"file_hash":     p.FileHash,
		"chunk_index":   p.ChunkIndex,
		"total_chunks":  p.TotalChunks,
		"content":       p.Content,
		"created_at":    p.CreatedAt.UTC().Format(time.RFC3339),
		"updated_at":    p.UpdatedAt.UTC().Format(time.RFC3339),
		"status":        p.Status,
		"mastery_level": p.MasteryLevel,
	}
}

func mapToPayload(m map[string]*qdrant.Value) Payload {
	getString := func(key string) string {
		v, ok := m[key]
		if !ok {
			return ""
		}
		return v.GetStringValue()
	}
	getInt := func(key string) int {
		v, ok := m[key]
		if !ok {
			return 0
		}
		return int(v.GetIntegerValue())
	}
	getFloat := func(key string) float64 {
		v, ok := m[key]
		if !ok {
			return 0
		}
		return v.GetDoubleValue()
	}
	parseTime := func(key string) time.Time {
		s := getString(key)
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return time.Time{}
		}
		return t
	}

	return Payload{
		Path:         getString("path"),
		Filename:     getString("filename"),
		Ext:          getString("ext"),
		FileHash:     getString("file_hash"),
		ChunkIndex:   getInt("chunk_index"),
		TotalChunks:  getInt("total_chunks"),
		Content:      getString("content"),
		CreatedAt:    parseTime("created_at"),
		UpdatedAt:    parseTime("updated_at"),
		Status:       getString("status"),
		MasteryLevel: getFloat("mastery_level"),
	}
}

func (s *QdrantStore) Upsert(ctx context.Context, chunks []Chunk) error {
	for i := 0; i < len(chunks); i += batchSize {
		end := i + batchSize
		if end > len(chunks) {
			end = len(chunks)
		}
		batch := chunks[i:end]

		points := make([]*qdrant.PointStruct, len(batch))
		for j := range batch {
			payload, err := qdrant.TryValueMap(payloadToMap(&batch[j].Payload))
			if err != nil {
				return fmt.Errorf("qdrant: payload map chunk %s: %w", batch[j].ID, err)
			}
			points[j] = &qdrant.PointStruct{
				Id:      qdrant.NewID(chunkIDToUUID(batch[j].ID)),
				Vectors: qdrant.NewVectorsDense(batch[j].Vector),
				Payload: payload,
			}
		}

		if _, err := s.client.Upsert(ctx, &qdrant.UpsertPoints{
			CollectionName: s.collection,
			Points:         points,
		}); err != nil {
			return fmt.Errorf("qdrant: upsert: %w", err)
		}
	}
	return nil
}

func (s *QdrantStore) Search(ctx context.Context, vector []float32, filter *Filter, limit uint64) ([]SearchResult, error) {
	req := &qdrant.QueryPoints{
		CollectionName: s.collection,
		Query:          qdrant.NewQuery(vector...),
		Limit:          qdrant.PtrOf(limit),
		WithPayload:    qdrant.NewWithPayload(true),
	}

	if filter != nil && filter.Status != nil {
		req.Filter = &qdrant.Filter{
			Must: []*qdrant.Condition{
				qdrant.NewMatchKeyword(fieldStatus, *filter.Status),
			},
		}
	}

	scored, err := s.client.Query(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("qdrant: search: %w", err)
	}

	results := make([]SearchResult, len(scored))
	for i, p := range scored {
		results[i] = SearchResult{
			ID:      p.GetId().GetUuid(),
			Score:   p.GetScore(),
			Payload: mapToPayload(p.GetPayload()),
		}
	}
	return results, nil
}

func (s *QdrantStore) Delete(ctx context.Context, ids []string) error {
	pointIDs := make([]*qdrant.PointId, len(ids))
	for i, id := range ids {
		pointIDs[i] = qdrant.NewID(chunkIDToUUID(id))
	}
	if _, err := s.client.Delete(ctx, &qdrant.DeletePoints{
		CollectionName: s.collection,
		Points:         qdrant.NewPointsSelectorIDs(pointIDs),
	}); err != nil {
		return fmt.Errorf("qdrant: delete: %w", err)
	}
	return nil
}

func (s *QdrantStore) UpdatePayload(ctx context.Context, id string, payload map[string]any) error {
	valueMap, err := qdrant.TryValueMap(payload)
	if err != nil {
		return fmt.Errorf("qdrant: update payload map: %w", err)
	}
	if _, err := s.client.SetPayload(ctx, &qdrant.SetPayloadPoints{
		CollectionName: s.collection,
		Payload:        valueMap,
		PointsSelector: qdrant.NewPointsSelector(qdrant.NewID(chunkIDToUUID(id))),
	}); err != nil {
		return fmt.Errorf("qdrant: set payload: %w", err)
	}
	return nil
}

func (s *QdrantStore) GetFileHash(ctx context.Context, path string) (string, error) {
	results, err := s.client.Scroll(ctx, &qdrant.ScrollPoints{
		CollectionName: s.collection,
		Filter: &qdrant.Filter{
			Must: []*qdrant.Condition{
				qdrant.NewMatchKeyword(fieldPath, path),
			},
		},
		Limit:       qdrant.PtrOf(uint32(1)),
		WithPayload: qdrant.NewWithPayloadInclude(fieldFileHash),
	})
	if err != nil {
		return "", fmt.Errorf("qdrant: get file hash: %w", err)
	}
	if len(results) == 0 {
		return "", nil
	}
	p := results[0].GetPayload()
	v, ok := p[fieldFileHash]
	if !ok {
		return "", nil
	}
	return v.GetStringValue(), nil
}
