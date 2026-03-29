package vectordb

import (
	"testing"
	"time"
)

func TestChunkIDToUUID_Deterministic(t *testing.T) {
	a := chunkIDToUUID("some-chunk-id")
	b := chunkIDToUUID("some-chunk-id")
	if a != b {
		t.Errorf("chunkIDToUUID not deterministic: %q != %q", a, b)
	}
}

func TestChunkIDToUUID_Unique(t *testing.T) {
	a := chunkIDToUUID("chunk-1")
	b := chunkIDToUUID("chunk-2")
	if a == b {
		t.Errorf("different IDs produced same UUID: %q", a)
	}
}

func TestPayloadRoundtrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	original := Payload{
		Path:         "/home/user/notes.md",
		Filename:     "notes.md",
		Ext:          ".md",
		FileHash:     "abc123",
		ChunkIndex:   2,
		TotalChunks:  5,
		Content:      "some content",
		CreatedAt:    now,
		UpdatedAt:    now,
		Status:       "reading",
		MasteryLevel: 0.42,
	}

	valueMap, err := TryValueMapFromPayload(&original)
	if err != nil {
		t.Fatalf("TryValueMapFromPayload: %v", err)
	}
	got := mapToPayload(valueMap)

	if got.Path != original.Path {
		t.Errorf("Path: want %q, got %q", original.Path, got.Path)
	}
	if got.FileHash != original.FileHash {
		t.Errorf("FileHash: want %q, got %q", original.FileHash, got.FileHash)
	}
	if got.ChunkIndex != original.ChunkIndex {
		t.Errorf("ChunkIndex: want %d, got %d", original.ChunkIndex, got.ChunkIndex)
	}
	if got.Status != original.Status {
		t.Errorf("Status: want %q, got %q", original.Status, got.Status)
	}
	if got.MasteryLevel != original.MasteryLevel {
		t.Errorf("MasteryLevel: want %v, got %v", original.MasteryLevel, got.MasteryLevel)
	}
	if !got.UpdatedAt.Equal(original.UpdatedAt) {
		t.Errorf("UpdatedAt: want %v, got %v", original.UpdatedAt, got.UpdatedAt)
	}
}
