# タスク管理 — magic_wand

> **⚠️ このファイルは参照用マスターリストです。**
> 進捗の正は **GitHub Issues** です。実作業は Issue を見てください。
> 現在状態の確認: `/status` コマンドを実行してください。

最終更新: 2026-03-27 (v3 — Phase 1完了)

凡例: `[ ]` 未着手 / `[~]` 進行中 / `[x]` 完了

---

## Phase 1: プロジェクト基盤

- [x] **T-01** Go module初期化 (`go mod init github.com/Kazakumo/magic_wand`)
- [x] **T-02** cobraによるCLI骨格作成（引数なし→TUI起動 / `ingest` / `query` / `watch`サブコマンド）
- [x] **T-03** `internal/config/` — viper設定ロード + `config.yaml`デフォルト
- [x] **T-04** `slog`ベースのロガー初期化（JSON/text切り替え対応）
- [x] **T-05** `docker-compose.yml` — Qdrant起動設定
- [x] **T-06** Makefile — `make run` / `make test` / `make lint` / `make up` / `make down` / `make tui`

---

## Phase 2: Embedding Layer

- [ ] **T-07** `internal/embedding/interface.go` — `Embedder` interface定義
- [ ] **T-08** `internal/embedding/ollama.go` — Ollama `/api/embeddings` クライアント実装
  - context.WithTimeout（30秒）
  - exponential backoffリトライ（最大3回: 1s/2s/4s）
- [ ] **T-09** Embedderの単体テスト（HTTPモックサーバー使用）

---

## Phase 3: Vector DB Layer

- [ ] **T-10** Qdrant公式クライアント(`qdrant/go-client`)追加 + gRPC接続確認
- [ ] **T-11** `internal/vectordb/interface.go` — `VectorStore` interface定義
  - `Upsert / Search / Delete / UpdatePayload / GetFileHash`
- [ ] **T-12** `internal/vectordb/schema.go` — ペイロード定義（拡張版）
  - `path / filename / ext / file_hash / chunk_index / total_chunks / content`
  - `created_at / updated_at / status / mastery_level`
- [ ] **T-13** `internal/vectordb/qdrant.go` — gRPCクライアント実装
  - コレクション自動作成（インデックス: path, status, updated_at）
  - Upsert（バッチアップサート）
  - Search（フィルタ付き）
  - Delete / UpdatePayload / GetFileHash
- [ ] **T-14** VectorStoreの統合テスト（Qdrant testcontainer or モック）

---

## Phase 4: Ingestion Pipeline

- [ ] **T-15** `internal/ingestion/hasher.go` — SHA256ハッシュ計算 + VectorStore照合
  - `ComputeHash(path string) (string, error)`
  - `IsChanged(path string, store VectorStore) (bool, error)`
- [ ] **T-16** `internal/ingestion/chunker.go` — スライディングウィンドウチャンキング
  - チャンクサイズ・オーバーラップ設定可能
  - `.go`ファイルはAST parseで関数境界を優先してチャンク分割
  - チャンクID生成: `sha256(filepath + chunk_index)`
- [ ] **T-17** `internal/ingestion/worker_pool.go` — goroutineプール実装
  - semaphoreによる並列数制御
  - CPU使用率モニタリングgoroutine（1秒ごと）
  - 動的ワーカー数増減ロジック
  - TUIへのイベント通知チャンネル（`chan IngestEvent`）
- [ ] **T-18** `internal/ingestion/watcher.go` — fsnotifyラッパー
  - 再帰ディレクトリ監視
  - イベントデバウンス（500ms）
  - Create / Write / Remove / Rename ハンドリング
  - SHA256デデュップ呼び出し（未変更ならスキップ）
- [ ] **T-19** インジェスション統合テスト（一時ディレクトリ使用）

---

## Phase 5: RAG Engine

- [ ] **T-20** `internal/rag/retriever.go` — クエリEmbedding → VectorStore.Search
- [ ] **T-21** `internal/rag/reranker.go` — 多段階リランキング
  - `CosineSimilarity(a, b []float32) float32`
  - `KeywordDensity(queryTokens []string, text string) float64`
  - `ASTSignatureWeight(text, ext string) float64`（.goのみAST解析）
  - `CombinedScore(cosine, density, ast, α, β, γ float64) float64`
  - スコア閾値フィルタリング
- [ ] **T-22** `internal/rag/generator.go` — Ollama chat streamingクライアント
  - プロンプトテンプレート（コンテキスト注入）
  - `chan string`でストリーミング返却
- [ ] **T-23** `internal/rag/pipeline.go` — RAGパイプライン統合
  - Retriever(top20) → Reranker → top5 → Generator
  - `[]SourceRef{Path, Score, Excerpt}`を合わせて返す
- [ ] **T-24** RAGパイプラインの単体テスト（各コンポーネントモック）

---

## Phase 6: Semantic Recommender

- [ ] **T-25** `internal/recommender/scorer.go`
  - `TimeWeight(updatedAt time.Time, lambda float64) float64` (exponential decay)
  - `StatusBoost(status string) float64` (unread=1.2, reading=1.0, mastered=0.5)
  - `FinalScore(cosine, timeWeight, statusBoost float64) float64`
- [ ] **T-26** `internal/recommender/analyzer.go`
  - クエリVec × 全チャンク類似度計算
  - 未習得フィルタ優先のVectorStore.Search
  - スコアリング・ランキング生成
- [ ] **T-27** `internal/recommender/recommender.go`
  - 推薦エンジン統合
  - 各推薦にLLMで一文の理由を生成
- [ ] **T-28** Recommenderの単体テスト

---

## Phase 7: TUI実装

### 7.1 スタイル・基盤

- [ ] **T-29** 依存パッケージ追加
  - `github.com/charmbracelet/bubbletea`
  - `github.com/charmbracelet/lipgloss`
  - `github.com/charmbracelet/bubbles`
  - `github.com/charmbracelet/glamour`
- [ ] **T-30** `internal/tui/styles/theme.go` — magic_wandカラーパレット定義
  - Midnight背景 / WandBlue / SpellPurple / ArcaneGreen / GlowAmber / CrimsonFlare
  - 共通スタイル（ボーダー・タイトル・バッジ・テーブル行）
- [ ] **T-31** `internal/tui/keymap.go` — グローバルキーバインド定義

### 7.2 共通コンポーネント

- [ ] **T-32** `internal/tui/components/tabbar.go` — タブナビゲーション（1〜5キー対応）
- [ ] **T-33** `internal/tui/components/statusbar.go` — ボトムバー（Watcher/Workers/CPU/Indexed）
- [ ] **T-34** `internal/tui/components/gauge.go` — CPU使用率ゲージ・プログレスバー
- [ ] **T-35** `internal/tui/components/badge.go` — ステータスバッジ・スコアバッジ（色分け）
- [ ] **T-36** `internal/tui/components/table.go` — カスタムテーブル（選択・フィルタ・ソート対応）
- [ ] **T-37** `internal/tui/components/chat.go` — チャットバブル（user/ai色分け）

### 7.3 各ビュー実装

- [ ] **T-38** `internal/tui/views/dashboard.go`
  - 統計カード（Total/Indexed/Status分布）
  - システムヘルス（Ollama/Qdrant/Watcher接続状態）
  - 最近のアクティビティリスト
- [ ] **T-39** `internal/tui/views/library.go`
  - 学習アイテムテーブル（Title/Type/Status/Mastery/Updated）
  - フィルタタブ（All/Unread/Reading/Mastered）
  - キー操作（`r`/`m`/`u`でステータス更新、`/`で検索）
  - VectorStore.UpdatePayload呼び出し
- [ ] **T-40** `internal/tui/views/query.go`
  - チャットUI（bubbles/viewport + テキスト入力）
  - ストリーミング回答表示（glamourでMarkdownレンダリング）
  - Sourcesパネル（スコア付きソース一覧）
  - セッション内のクエリ履歴保持
- [ ] **T-41** `internal/tui/views/recommend.go`
  - クエリ入力フィールド
  - 推薦ランキングリスト（スコアバー・ステータスバッジ・理由・更新日）
  - `Enter`でQueryビューに遷移して詳細検索
- [ ] **T-42** `internal/tui/views/ingest.go`
  - ワーカープール可視化（各ワーカーの状態アニメーション）
  - CPU使用率ゲージ（リアルタイム更新）
  - キュー・完了数・スキップ数カウンター
  - イベントログ（スクロール可能なviewport）
  - IngestEventチャンネルの購読

### 7.4 ルートモデル・統合

- [ ] **T-43** `internal/tui/model.go` — ルートモデル（タブ管理・グローバル状態・リサイズ処理）
- [ ] **T-44** `internal/tui/app.go` — tea.Program初期化・全コンポーネント配線
- [ ] **T-45** ヘルプオーバーレイ実装（`?`キー → キーバインド一覧表示）
- [ ] **T-46** ウィンドウリサイズ対応（各ビューへの伝播・レイアウト再計算）

---

## Phase 8: CLI非インタラクティブモード

- [ ] **T-47** `magic_wand watch` — バックグラウンドウォッチャー起動（TUI不要）
- [ ] **T-48** `magic_wand ingest <path>` — 指定パスを手動インジェスト + 進捗表示
- [ ] **T-49** `magic_wand query "<text>"` — 非インタラクティブRAG検索（stdout出力）

---

## Phase 9: 品質・仕上げ

- [ ] **T-50** E2Eテスト（`watch`起動 → ファイル追加 → TUI Libraryで確認 → `query`）
- [ ] **T-51** README.md 更新（インストール・Ollama/Docker起動手順・TUIスクリーンショット）
- [x] **T-52** golangci-lint設定 (`.golangci.yml`)
- [x] **T-53** GitHub Actions CI（lint + test）

---

## 実装順序の推奨

```
Phase 1  (T-01〜06)   基盤・プロジェクト構造
  ↓
Phase 2  (T-07〜09)   Embedding（Ollama疎通確認）
  ↓
Phase 3  (T-10〜14)   VectorDB（Qdrant疎通確認 ← ここで拡張ペイロード含め動作確認）
  ↓
Phase 4  (T-15〜19)   Ingestion Pipeline（SHA256デデュップ含め完成）
  ↓
Phase 5  (T-20〜24)   RAG Engine（多段階リランカー含め完成）
  ↓
Phase 6  (T-25〜28)   Recommender（習熟度フィルタ含め完成）
  ↓
Phase 7  (T-29〜46)   TUI（スタイル→コンポーネント→各ビュー→統合）
  ↓
Phase 8  (T-47〜49)   CLI非インタラクティブモード
  ↓
Phase 9  (T-50〜53)   品質・仕上げ
```

---

## 技術的判断ログ

| 日付 | 判断内容 | 理由 |
|------|---------|------|
| 2026-03-27 | Qdrant gRPCクライアントに公式`qdrant-go`採用 | proto自前管理よりメンテが楽 |
| 2026-03-27 | Embeddingモデルは`nomic-embed-text`をデフォルト | 768次元・多言語対応・ローカル動作 |
| 2026-03-27 | LangChainは不使用 | Go実装力・RAGの仕組み理解が目的 |
| 2026-03-27 | TUIはbubbletea + lipgloss + bubbles | Charm社エコシステムが最も成熟・美しい |
| 2026-03-27 | 進捗管理はVectorDB payloadに持たせる | 別DBを増やさず、フィルタ検索も効く |
| 2026-03-27 | SHA256デデュップをwatcherに追加 | 無駄なEmbedding生成によるCPU消費を防ぐ |
| 2026-03-27 | Rerankerを3段階（cosine+density+AST）に強化 | .goのシグネチャを重視した検索精度向上 |
| 2026-03-27 | TUIを主インターフェース（引数なし起動）に昇格 | 視認性・操作性を最優先。CLI は補助 |
