//go:build integration

package vectordb_test

import (
	"context"
	"testing"
	"time"

	"github.com/Kazakumo/magic_wand/internal/vectordb"
)

// go test -tags integration -run TestQdrantStore ./internal/vectordb/
// 事前に make up (Qdrant on localhost:6334) が必要。

func TestQdrantStore_UpsertAndSearch(t *testing.T) {
	ctx := context.Background()
	store, err := vectordb.NewQdrantStore(ctx, "localhost", 6334, "magic_wand_test")
	if err != nil {
		t.Fatalf("NewQdrantStore: %v", err)
	}

	now := time.Now().UTC().Truncate(time.Second)
	chunks := []vectordb.Chunk{
		{
			ID:     "test-chunk-001",
			Vector: make([]float32, 768),
			Payload: vectordb.Payload{
				Path:        "/test/file.md",
				Filename:    "file.md",
				Ext:         ".md",
				FileHash:    "hash001",
				ChunkIndex:  0,
				TotalChunks: 1,
				Content:     "test content",
				CreatedAt:   now,
				UpdatedAt:   now,
				Status:      "unread",
			},
		},
	}

	if err := store.Upsert(ctx, chunks); err != nil {
		t.Fatalf("Upsert: %v", err)
	}

	hash, err := store.GetFileHash(ctx, "/test/file.md")
	if err != nil {
		t.Fatalf("GetFileHash: %v", err)
	}
	if hash != "hash001" {
		t.Errorf("GetFileHash: want %q, got %q", "hash001", hash)
	}

	results, err := store.Search(ctx, make([]float32, 768), nil, 10)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("Search: expected results, got none")
	}

	if err := store.UpdatePayload(ctx, "test-chunk-001", map[string]any{"status": "reading"}); err != nil {
		t.Fatalf("UpdatePayload: %v", err)
	}

	status := vectordb.PtrOf("reading")
	filtered, err := store.Search(ctx, make([]float32, 768), &vectordb.Filter{Status: status}, 10)
	if err != nil {
		t.Fatalf("Search with filter: %v", err)
	}
	if len(filtered) == 0 {
		t.Fatal("filtered search: expected results, got none")
	}

	if err := store.Delete(ctx, []string{"test-chunk-001"}); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}
