package vectordb

import "time"

// Chunk はベクトルDB に格納する 1 チャンク。
type Chunk struct {
	ID      string    // sha256(filepath + chunk_index)
	Vector  []float32 // embedding vector
	Payload Payload
}

// Payload はベクトルDBのメタデータ。
type Payload struct {
	Path         string    `json:"path"`
	Filename     string    `json:"filename"`
	Ext          string    `json:"ext"`
	FileHash     string    `json:"file_hash"`
	ChunkIndex   int       `json:"chunk_index"`
	TotalChunks  int       `json:"total_chunks"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Status       string    `json:"status"`        // unread | reading | mastered
	MasteryLevel float64   `json:"mastery_level"` // 0.0 〜 1.0
}
