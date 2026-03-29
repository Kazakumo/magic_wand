package vectordb

import (
	"testing"
	"time"
)

const (
	hexID1 = "a3f1c2d4e5b6a7f8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2"
	hexID2 = "b4e2d3c5f6a7b8e9d0c1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2c3"
)

func TestChunkIDToUUID_Deterministic(t *testing.T) {
	a := chunkIDToUUID(hexID1)
	b := chunkIDToUUID(hexID1)
	if a != b {
		t.Errorf("chunkIDToUUID not deterministic: %q != %q", a, b)
	}
}

func TestChunkIDToUUID_Unique(t *testing.T) {
	a := chunkIDToUUID(hexID1)
	b := chunkIDToUUID(hexID2)
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
