package vectordb

import (
	"context"
	"sync"
)

// MockVectorStore は VectorStore のインメモリ実装。他パッケージのテスト用。
type MockVectorStore struct {
	mu     sync.RWMutex
	chunks map[string]Chunk
	hashes map[string]string // path -> file_hash
}

func NewMockVectorStore() *MockVectorStore {
	return &MockVectorStore{
		chunks: make(map[string]Chunk),
		hashes: make(map[string]string),
	}
}

func (m *MockVectorStore) Upsert(_ context.Context, chunks []Chunk) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range chunks {
		m.chunks[chunks[i].ID] = chunks[i]
		m.hashes[chunks[i].Payload.Path] = chunks[i].Payload.FileHash
	}
	return nil
}

func (m *MockVectorStore) Search(_ context.Context, _ []float32, filter *Filter, limit uint64) ([]SearchResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var results []SearchResult
	for id := range m.chunks {
		if filter != nil && filter.Status != nil && m.chunks[id].Payload.Status != *filter.Status {
			continue
		}
		results = append(results, SearchResult{
			ID:      m.chunks[id].ID,
			Score:   1.0,
			Payload: m.chunks[id].Payload,
		})
		if uint64(len(results)) >= limit {
			break
		}
	}
	return results, nil
}

func (m *MockVectorStore) Delete(_ context.Context, ids []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, id := range ids {
		delete(m.chunks, id)
	}
	return nil
}

func (m *MockVectorStore) UpdatePayload(_ context.Context, id string, payload map[string]any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.chunks[id]
	if !ok {
		return nil
	}
	if s, ok := payload["status"].(string); ok {
		c.Payload.Status = s
	}
	if v, ok := payload["mastery_level"].(float64); ok {
		c.Payload.MasteryLevel = v
	}
	m.chunks[id] = c
	return nil
}

func (m *MockVectorStore) GetFileHash(_ context.Context, path string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.hashes[path], nil
}
