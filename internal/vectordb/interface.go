package vectordb

import "context"

// VectorStore はベクトルDBの操作インターフェース。
type VectorStore interface {
	// Upsert はチャンクをバッチでアップサートする。
	Upsert(ctx context.Context, chunks []Chunk) error

	// Search はクエリベクトルで類似検索し上位 limit 件を返す。
	Search(ctx context.Context, vector []float32, filter *Filter, limit uint64) ([]SearchResult, error)

	// Delete は指定 ID のベクトルを削除する。
	Delete(ctx context.Context, ids []string) error

	// UpdatePayload は指定 ID のペイロードを部分更新する。
	UpdatePayload(ctx context.Context, id string, payload map[string]any) error

	// GetFileHash は path に対応するファイルハッシュを返す。存在しない場合は空文字。
	GetFileHash(ctx context.Context, path string) (string, error)
}

// Filter は Search 時のペイロードフィルタ条件。
type Filter struct {
	Status *string // "unread" | "reading" | "mastered"
}

// PtrOf はポインタを生成するジェネリックヘルパー。
func PtrOf[T any](v T) *T { return &v }

// SearchResult は検索結果の 1 件。
type SearchResult struct {
	ID      string
	Score   float32
	Payload Payload
}
