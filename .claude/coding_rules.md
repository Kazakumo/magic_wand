# コーディング規則 — magic_wand

ベース: [Effective Go](https://go.dev/doc/effective_go)
最終更新: 2026-03-27

---

## 1. フォーマット

- `gofmt` / `goimports` を必ず通す（CIでも検証）
- タブインデント（スペース禁止）
- 行長の制限なし（gofmtに従う）
- `if` / `for` / `switch` の開き括弧は同じ行

---

## 2. 命名規則

### パッケージ名
```go
// Good: 短く・小文字・単数形
package embedding
package vectordb
package rag

// Bad
package embeddingClient
package vector_db
```

### 変数・関数・型
```go
// MixedCaps（アンダースコア禁止）
workerPool    // private
WorkerPool    // public

// 略語は全大文字か全小文字で統一
urlParser     // private
URLParser     // public（Url は NG）
grpcClient    // private（Grpc は NG）
```

### インターフェース
```go
// 単一メソッドは -er 接尾辞
type Embedder interface { Embed(...) }
type Retriever interface { Retrieve(...) }
type Reranker interface { Rerank(...) }

// 複数メソッドは責務を表す名詞
type VectorStore interface { ... }
```

### Getter
```go
// Get接頭辞は不要
func (c *Client) Status() Status { ... }  // Good
func (c *Client) GetStatus() Status { ... } // Bad
```

### エラー変数
```go
// エラー型は Error 接尾辞
type EmbedError struct { ... }

// センチネルエラーは Err 接頭辞
var ErrCollectionNotFound = errors.New("collection not found")
var ErrHashMismatch      = errors.New("hash mismatch")
```

---

## 3. エラーハンドリング

```go
// Early return で深いネストを避ける
func (p *Pipeline) Run(ctx context.Context, query string) (*Result, error) {
    vec, err := p.embedder.Embed(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("embed query: %w", err)
    }

    results, err := p.store.Search(ctx, vec, p.cfg.TopK, nil)
    if err != nil {
        return nil, fmt.Errorf("search: %w", err)
    }
    // ...
}

// エラーは必ずラップして文脈を付ける（%w）
// エラーを無視しない（_ は原則禁止）
// panic は本当に回復不可能な場合のみ（プログラムバグ）
```

---

## 4. インターフェース設計

```go
// 小さく保つ（1〜2メソッド推奨）
// 実装側ではなく利用側で定義する

// Good: internal/rag/retriever.go が利用側で定義
type VectorSearcher interface {
    Search(ctx context.Context, vec []float32, topK int, f *Filter) ([]SearchResult, error)
}

// Bad: 20メソッドの巨大インターフェース
```

- **インターフェースの実装チェックはコンパイル時に行う**
```go
var _ Embedder = (*OllamaEmbedder)(nil)
```

---

## 5. 並行処理（このプロジェクト特有）

```go
// 原則: 「通信によってメモリを共有する」
// チャンネルで goroutine 間のデータをやり取りする

// Worker Pool パターン（internal/ingestion/worker_pool.go）
sem := make(chan struct{}, maxWorkers)
for _, chunk := range chunks {
    sem <- struct{}{}
    go func(c Chunk) {
        defer func() { <-sem }()
        process(c)
    }(chunk)
}

// context.Context を必ず引数の第一引数に
func (e *OllamaEmbedder) Embed(ctx context.Context, text string) ([]float32, error)

// goroutine のリークを防ぐ: context キャンセルを必ず監視
select {
case <-ctx.Done():
    return ctx.Err()
case result := <-ch:
    return result, nil
}
```

---

## 6. defer の使い方

```go
// リソース管理に使う（Close / Unlock など）
conn, err := grpc.Dial(addr, opts...)
if err != nil { return err }
defer conn.Close()

// ループ内では使わない（メモリリークの原因）
for _, path := range paths {
    f, _ := os.Open(path)
    defer f.Close() // Bad: ループが終わるまでCloseされない
}
```

---

## 7. 構造体・初期化

```go
// フィールド名を明示する（順序依存にしない）
cfg := WorkerConfig{
    MaxWorkers:   runtime.NumCPU(),
    CPUThreshold: 0.70,
    PollInterval: time.Second,
}

// ゼロ値を活用する（不要な初期化をしない）
var mu sync.Mutex  // そのまま使える
```

---

## 8. パッケージ設計

```go
// internal/ 以下は外部公開しない
// 循環インポート禁止（依存方向: tui → rag/recommender → vectordb → embedding）
// パッケージ間の通信はインターフェース経由

// ファイル分割の基準
// 1ファイル = 1つの責務（rag/retriever.go, rag/reranker.go, rag/generator.go）
```

---

## 9. テスト

```go
// テストファイルは同パッケージ（ブラックボックス非推奨）
// ファイル名: xxx_test.go
// 関数名: Test<対象関数名>_<シナリオ>
func TestCosineSimilarity_Identical(t *testing.T) { ... }
func TestCosineSimilarity_Orthogonal(t *testing.T) { ... }

// テーブルドリブンテスト
tests := []struct {
    name  string
    input []float32
    want  float32
}{
    {"identical vectors", []float32{1, 0}, 1.0},
    {"orthogonal",        []float32{0, 1}, 0.0},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) { ... })
}

// 外部依存はインターフェースでモック
// Qdrant/Ollama は testcontainer か httptest.Server でモック
```

---

## 10. このプロジェクト固有のガードレール

| ルール | 理由 |
|--------|------|
| `context.Context` は常に第一引数 | タイムアウト・キャンセルを全レイヤーで制御 |
| エラーは必ず `%w` でラップ | スタック全体のエラー文脈を保持 |
| goroutine を起動したら必ず終了を保証 | TUI の長時間稼働でリークを防ぐ |
| ファイルパスは `path/filepath` を使う | クロスプラットフォーム対応 |
| `log` パッケージを直接使わない | `slog` を使う（構造化・レベル制御） |
| `panic` は main / init 以外では書かない | ライブラリコードは error を返す |
| マジックナンバーは定数か config に | 0.70 / 512 / 64 などは全て config 経由 |
