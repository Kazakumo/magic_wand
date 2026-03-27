# 基本設計 — magic_wand

最終更新: 2026-03-27 (v2 — TUI・進捗管理・インジェスト最適化・リランカー強化を追加)

---

## 1. システムアーキテクチャ全体図

```
┌──────────────────────────────────────────────────────────────────┐
│                    TUI (bubbletea)                               │
│  Dashboard │ Library │ Query │ Recommend │ Ingest               │
└──────┬──────────┬──────────┬───────────────┬────────────────────┘
       │          │          │               │
       │    ┌─────▼──────┐  ┌▼────────────┐  │
       │    │ RAG Engine │  │ Recommender │  │
       │    │ retriever  │  │ analyzer    │  │
       │    │ reranker   │  │ scorer      │  │
       │    │ generator  │  └──────┬──────┘  │
       │    └─────┬──────┘         │         │
       │          │                │         │
       │    ┌─────▼────────────────▼──────┐  │
       │    │       Vector DB Layer       │  │
       │    │    Qdrant (gRPC client)     │  │
       │    └─────────────┬──────────────┘  │
       │                  │                 │
       │    ┌─────────────▼─────────────────▼──────┐
       │    │        Ingestion Pipeline             │
       │    │  fsnotify → dedup(SHA256)             │
       │    │  → Chunker → Worker Pool              │
       │    │  → Embedding Client                   │
       │    └────────────────────┬─────────────────┘
       │                         │
       └──────────── ┌───────────▼──────────┐
                     │   Ollama (local)     │
                     │  embedding + chat    │
                     └──────────────────────┘
```

---

## 2. ディレクトリ構成

```
magic_wand/
├── cmd/
│   └── magic_wand/
│       └── main.go               # エントリポイント（TUI起動 or CLIサブコマンド）
├── internal/
│   ├── config/
│   │   └── config.go             # 設定構造体・viperロード
│   ├── ingestion/
│   │   ├── watcher.go            # fsnotifyラッパー + デバウンス
│   │   ├── chunker.go            # スライディングウィンドウチャンキング
│   │   ├── hasher.go             # SHA256ハッシュ計算・重複チェック
│   │   └── worker_pool.go        # goroutineプール・CPU動的制御
│   ├── embedding/
│   │   ├── interface.go          # Embedder interface
│   │   └── ollama.go             # Ollama embeddingクライアント
│   ├── vectordb/
│   │   ├── interface.go          # VectorStore interface
│   │   ├── qdrant.go             # Qdrant gRPCクライアント実装
│   │   └── schema.go             # コレクション・ペイロード定義
│   ├── rag/
│   │   ├── retriever.go          # クエリEmbedding → VectorStore.Search
│   │   ├── reranker.go           # 多段階リランキング（cosine + density + AST）
│   │   ├── generator.go          # Ollama chat streaming
│   │   └── pipeline.go           # RAGパイプライン統合
│   ├── recommender/
│   │   ├── scorer.go             # exponential decay + 習熟度フィルタ
│   │   ├── analyzer.go           # クエリ×履歴の類似度計算
│   │   └── recommender.go        # 推薦エンジン統合
│   └── tui/
│       ├── app.go                # tea.Program セットアップ・起動
│       ├── model.go              # ルートモデル（タブ管理・グローバル状態）
│       ├── keymap.go             # グローバルキーバインド定義
│       ├── styles/
│       │   └── theme.go          # lipglossテーマ・カラーパレット
│       ├── views/
│       │   ├── dashboard.go      # ダッシュボードビュー
│       │   ├── library.go        # ライブラリ（学習アイテム管理）ビュー
│       │   ├── query.go          # RAGチャットビュー
│       │   ├── recommend.go      # 推薦ビュー
│       │   └── ingest.go         # インジェスト状態ビュー
│       └── components/
│           ├── tabbar.go         # 上部タブナビゲーション
│           ├── statusbar.go      # 下部ステータスバー
│           ├── table.go          # 拡張テーブルコンポーネント
│           ├── chat.go           # チャットメッセージバブル
│           ├── gauge.go          # CPU/進捗ゲージ
│           └── badge.go          # スコア・ステータスバッジ
├── docker-compose.yml
├── config.yaml
├── Makefile
└── .claude/
    ├── requirements.md
    ├── design.md
    └── tasks.md
```

---

## 3. TUI設計

### 3.1 画面レイアウト

```
╔═══════════════════════════════════════════════════════════════╗
║  ✦ magic_wand                                                 ║
╠═══════════════════════════════════════════════════════════════╣
║  [1] Dashboard  [2] Library  [3] Query  [4] Recommend  [5] Ingest  ║
╠═══════════════════════════════════════════════════════════════╣
║                                                               ║
║              (アクティブビューのコンテンツ)                    ║
║                                                               ║
╠═══════════════════════════════════════════════════════════════╣
║  ● Watch:ON  ◈ Workers:3/4  ▸ CPU:42%  ◉ Indexed:1,234  [?] ║
╚═══════════════════════════════════════════════════════════════╝
```

### 3.2 カラーパレット（magic_wand テーマ）

| 用途 | 色名 | Hex |
|------|------|-----|
| 背景 | Midnight | `#1A1B26` |
| メインテキスト | StarDust | `#C0CAF5` |
| プライマリ | WandBlue | `#7AA2F7` |
| アクセント | SpellPurple | `#BB9AF7` |
| 成功・習得済 | ArcaneGreen | `#9ECE6A` |
| 警告・読中 | GlowAmber | `#E0AF68` |
| エラー | CrimsonFlare | `#F7768E` |
| ボーダー | NightBorder | `#3B4261` |
| サブテキスト | DimGlow | `#565F89` |

### 3.3 各ビュー詳細

#### Dashboard（`1`キー）
```
┌─ Summary ──────────────────────┐  ┌─ System Health ──────────────┐
│  Total Items      128           │  │  Ollama  ● Connected         │
│  Indexed Chunks   4,821         │  │  Qdrant  ● Connected         │
│  Unread    ████░░░░  42 (33%)   │  │  Watcher ● Running           │
│  Reading   ██░░░░░░  18 (14%)   │  └──────────────────────────────┘
│  Mastered  ██████░░  68 (53%)   │
└────────────────────────────────┘  ┌─ Recent Activity ────────────┐
                                    │  ● golang/channels.md  2m ago │
                                    │  ● notes/ddd.md        5m ago │
                                    │  ● src/main.go        12m ago │
                                    └──────────────────────────────┘
```

#### Library（`2`キー）
```
  Filter: [All] Unread  Reading  Mastered     Search: /

  Title                    Type   Status    Mastery  Updated
  ──────────────────────────────────────────────────────────
  ▶ golang/goroutines.md   .md   ◉ Reading  ████░░  2d ago
    design_patterns.md     .md   ○ Unread   ░░░░░░  1w ago
    src/rag/pipeline.go    .go   ✓ Mastered ██████  3h ago
  ──────────────────────────────────────────────────────────
  r: reading  m: mastered  u: unread  /: search  Enter: open
```

#### Query（`3`キー）
```
  ┌─ Chat ─────────────────────────────────────────────────┐
  │                                                         │
  │  You: Goのgoroutineでデッドロックが起きる原因は？        │
  │                                                         │
  │  ✦ AI: チャンネルの送受信が対応していない場合...        │
  │     ▌ (ストリーミング中...)                             │
  │                                                         │
  └─────────────────────────────────────────────────────────┘
  ┌─ Sources (3 found) ────────────────────────────────────┐
  │  [0.91] golang/goroutines.md:§チャンネル同期           │
  │  [0.84] notes/go_concurrency.md:§デッドロック回避      │
  └─────────────────────────────────────────────────────────┘
  > ▌                              Enter: 送信  Ctrl+C: クリア
```

#### Recommend（`4`キー）
```
  Current challenge: [goroutineのデッドロック                  ]
                                                     Enter: 推薦を取得

  Recommendations:
  ──────────────────────────────────────────────────────────────
  #1 ████ 0.91  golang/channels.md              ○ Unread
     理由: goroutineの同期パターンと現在の課題が強く一致します
     更新: 3日前

  #2 ███░ 0.78  notes/go_concurrency.md         ◉ Reading
     理由: デッドロック回避のパターンを詳しく説明しています
     更新: 1週間前
  ──────────────────────────────────────────────────────────────
  Enter: Queryビューで詳細検索  m: 習得済にマーク
```

#### Ingest（`5`キー）
```
  ┌─ Worker Pool ──────────────────────────────────────────┐
  │  Worker 1  [████████████████░░░░]  embedding...        │
  │  Worker 2  [████████████████████]  writing to DB       │
  │  Worker 3  [░░░░░░░░░░░░░░░░░░░░]  idle               │
  │  Worker 4  [░░░░░░░░░░░░░░░░░░░░]  idle (CPU 68%)     │
  └────────────────────────────────────────────────────────┘
  ┌─ CPU Usage ─┐  ┌─ Queue ─────┐  ┌─ Stats ────────────┐
  │ ████░░ 68%  │  │  Pending: 3 │  │  Done:   142 chunks │
  └─────────────┘  └─────────────┘  │  Skipped:  38 (dup) │
                                     └────────────────────┘
  ┌─ Event Log ────────────────────────────────────────────┐
  │  [10:42:01] ✓ golang/goroutines.md → 12 chunks        │
  │  [10:42:03] ⟳ notes/ddd.md (unchanged, skipped)      │
  │  [10:42:05] ✓ src/main.go → 8 chunks                 │
  └────────────────────────────────────────────────────────┘
```

### 3.4 グローバルキーバインド

| キー | アクション |
|------|-----------|
| `1`〜`5` | ビュー切り替え |
| `Tab` / `Shift+Tab` | 次/前のビュー |
| `?` | ヘルプオーバーレイ |
| `q` / `Ctrl+C` | 終了 |
| `Ctrl+W` | Watcherトグル |

---

## 4. コンポーネント詳細設計

### 4.1 Config

```go
type Config struct {
    Watch      WatchConfig
    Ollama     OllamaConfig
    Qdrant     QdrantConfig
    RAG        RAGConfig
    Worker     WorkerConfig
    Recommender RecommenderConfig
    Log        LogConfig
}

type WorkerConfig struct {
    MaxWorkers   int
    CPUThreshold float64       // 0.70
    PollInterval time.Duration // 1s
}
```

---

### 4.2 Ingestion Pipeline

#### ファイル監視フロー（SHA256デデュップ追加）
```
fsnotify.Event
    → イベントデバウンス（500ms）         ← fsnotify多重発火を吸収
    → ファイル種別フィルタ
    → Hasher.Compute(path) → sha256
    → VectorStore.GetHash(path)
    → [ハッシュ一致?] → スキップ（Ingestビューにskipped通知）
    → [ハッシュ不一致 or 新規]
    → Chunker.Split(content)
    → workerPool.Submit(chunks...)
        → Embedder.Embed(chunk.Text)
        → VectorStore.Upsert(vector, metadata{..., file_hash: sha256})
```

#### チャンキング戦略
- チャンクサイズ: 512 tokens（設定可能）
- オーバーラップ: 64 tokens
- チャンクID: `sha256(filepath + chunk_index)`
- `.go`ファイルは関数境界を優先してチャンク分割（AST parse）

#### Worker Pool & CPU制御
```
goroutine: CPU Monitor（1秒ごと）
  → 使用率 > threshold        → activeWorkers = max(1, activeWorkers-1)
  → 使用率 < threshold * 0.8  → activeWorkers = min(MaxWorkers, activeWorkers+1)
  → TUIのIngestビューにworker状態をチャンネル送信
```

---

### 4.3 Embedding Layer

```go
type Embedder interface {
    Embed(ctx context.Context, text string) ([]float32, error)
    Dimensions() int
}
```

- エンドポイント: `POST /api/embeddings`
- モデル: `nomic-embed-text`（768次元、多言語対応）
- タイムアウト: 30秒
- リトライ: exponential backoff（最大3回、1s/2s/4s）

---

### 4.4 Vector DB Layer

```go
type VectorStore interface {
    Upsert(ctx context.Context, points []Point) error
    Search(ctx context.Context, query []float32, topK int, filter *Filter) ([]SearchResult, error)
    Delete(ctx context.Context, ids []string) error
    UpdatePayload(ctx context.Context, id string, payload map[string]any) error
    GetFileHash(ctx context.Context, path string) (string, error)  // SHA256取得
}
```

#### Qdrantペイロード設計（拡張）
```
Collection: "learning_chunks"
  vector_size: 768
  distance:    Cosine
  payload:
    # ファイル情報
    path:           string     // ~/notes/golang/goroutines.md
    filename:       string     // goroutines.md
    ext:            string     // .md
    file_hash:      string     // SHA256（重複チェック・更新検知用）

    # チャンク情報
    chunk_index:    int
    total_chunks:   int
    content:        string     // 表示・プロンプト注入用テキスト

    # タイムスタンプ
    created_at:     int64      // unix timestamp
    updated_at:     int64

    # 進捗管理（NEW）
    status:         string     // "unread" | "reading" | "mastered"
    mastery_level:  float64    // 0.0 ~ 1.0（アウトプット類似度から自動計算）
```

**インデックス設計**
- `path` にKeywordインデックス（フィルタ用）
- `status` にKeywordインデックス（LibraryビューのStatusフィルタ用）
- `updated_at` にNumericインデックス（時刻フィルタ用）

---

### 4.5 RAG Engine — 多段階リランキング（強化）

#### パイプライン
```
Input: userQuery (string)
  → Embed(userQuery) → queryVec
  → VectorStore.Search(queryVec, topK=20)
  → Reranker.Rerank(queryVec, results)
      Stage 1: CosineSimilarity(queryVec, chunkVec)
      Stage 2: KeywordDensity(queryTokens, chunkText)
      Stage 3: ASTWeight(chunkText, ext)           ← .goなら関数シグネチャに加点
      CombinedScore = α·cosine + β·keyword + γ·ast
  → FilterByThreshold(score > 0.60)
  → BuildPrompt(query, top5chunks)
  → Generator.Stream(prompt) → chan string
Output: streamedAnswer + []SourceRef{Path, Score, Excerpt}
```

#### 多段階スコアリング式
```go
// Stage 1: コサイン類似度（ベクトル空間での意味的近さ）
cosine := CosineSimilarity(queryVec, chunkVec)

// Stage 2: キーワード密度（クエリトークンがチャンクに何割含まれるか）
density := KeywordDensity(queryTokens, chunkText)
// = len(intersection(queryTokens, chunkTokens)) / len(queryTokens)

// Stage 3: AST重み（.goファイルのみ）
astWeight := 0.0
if ext == ".go" {
    // 関数/型シグネチャが含まれていれば加点
    astWeight = ASTSignatureWeight(chunkText)
}

combined := α*cosine + β*density + γ*astWeight
// デフォルト: α=0.6, β=0.3, γ=0.1
```

---

### 4.6 Semantic Recommender（習熟度フィルタ追加）

#### スコアリング式
```
final_score = cosine_similarity × time_weight × status_boost

time_weight(t)   = exp(-λ × days_since_update)    λ=0.01
status_boost(s)  = { unread: 1.2, reading: 1.0, mastered: 0.5 }
                   ← 未習得を優先的に推薦
```

#### 推薦ロジック
1. クエリをEmbedding
2. VectorStore.Search(topK=30, filter={status: [unread, reading]})を優先
3. スコアリング適用
4. Top-N（デフォルト5）を返す
5. 各推薦にLLMで一文の理由を生成

---

## 5. 技術スタック

| 分類 | 採用技術 | 理由 |
|------|---------|------|
| 言語 | Go 1.22+ | goroutine・channelによる並列処理 |
| TUIフレームワーク | bubbletea | Elm-likeアーキテクチャ、最もGoらしいTUI |
| TUIスタイリング | lipgloss | 宣言的スタイル・カラーパレット |
| TUIコンポーネント | bubbles | list, table, textinput, viewport, spinner, progress |
| Markdownレンダリング | glamour | QueryビューのAI回答をマークダウン整形 |
| CLIフレームワーク | cobra | 非インタラクティブサブコマンド用 |
| ファイル監視 | fsnotify | クロスプラットフォーム対応 |
| ベクトルDB | Qdrant | gRPC対応・ローカル動作・高性能 |
| gRPC | google.golang.org/grpc | Qdrant公式クライアント |
| ローカルLLM | Ollama | embedding + chat両対応 |
| 設定管理 | viper | YAML + env override |
| 構造化ログ | slog（標準ライブラリ） | JSON出力、Go1.21+ |
| コンテナ | Docker Compose | Qdrant起動管理 |

---

## 6. TUIアーキテクチャ（Elm-likeモデル）

```go
// ルートモデル
type Model struct {
    activeTab   Tab
    tabs        []TabModel      // 各ビューのモデル
    statusbar   StatusBarModel
    keymap      KeyMap
    windowWidth int
    windowHeight int
    // グローバル共有状態
    ingestEvents chan IngestEvent   // IngestビューへのリアルタイムイベントchannelBridge
    workerState  WorkerState
}

// bubbletea Update関数
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // グローバルキー処理 → タブ切り替え
    case tea.WindowSizeMsg:
        // リサイズ → 全ビューに伝播
    case IngestEvent:
        // Ingestビューモデルに転送
    }
    // アクティブタブのUpdate委譲
    return m, nil
}
```

---

## 7. Docker Compose

```yaml
services:
  qdrant:
    image: qdrant/qdrant:latest
    ports:
      - "6333:6333"   # HTTP API
      - "6334:6334"   # gRPC
    volumes:
      - qdrant_data:/qdrant/storage
volumes:
  qdrant_data:
```

---

## 8. 設定ファイル（config.yaml）

```yaml
watch:
  dirs: [~/notes, ~/dev]
  extensions: [".md", ".txt", ".go", ".py", ".rs"]

ollama:
  base_url: "http://localhost:11434"
  embed_model: "nomic-embed-text"
  chat_model: "llama3.2"
  timeout: 30s

qdrant:
  host: "localhost"
  grpc_port: 6334
  collection: "learning_chunks"

worker:
  max_workers: 4
  cpu_threshold: 0.70
  poll_interval: 1s

rag:
  top_k: 20
  score_threshold: 0.60
  alpha: 0.6   # cosine weight
  beta:  0.3   # keyword density weight
  gamma: 0.1   # AST weight

recommender:
  top_n: 5
  time_decay_lambda: 0.01
  status_boost:
    unread:   1.2
    reading:  1.0
    mastered: 0.5

log:
  level: "info"
  format: "json"
```

---

## 9. 技術的判断ログ

| 日付 | 判断内容 | 理由 |
|------|---------|------|
| 2026-03-27 | Qdrant gRPCクライアントに公式`qdrant-go`採用 | proto自前管理よりメンテが楽 |
| 2026-03-27 | Embeddingモデルは`nomic-embed-text`をデフォルト | 768次元・多言語対応・ローカル動作 |
| 2026-03-27 | LangChainは不使用 | Go実装力・RAGの仕組み理解が目的 |
| 2026-03-27 | TUIはbubbletea + lipgloss + bubbles | Charm社のエコシステムが最も成熟・美しい |
| 2026-03-27 | 進捗管理はVectorDB payloadに持たせる | 別DBを増やさず、フィルタ検索も効く |
| 2026-03-27 | SHA256デデュップをwatcherに追加 | 無駄なEmbedding生成によるCPU消費を防ぐ |
| 2026-03-27 | Rerankerを3段階（cosine+density+AST）に強化 | .goファイルのシグネチャを重視した検索精度向上 |
