# magic_wand — CLAUDE.md

ローカルファーストの分散型AI学習オーケストレーター。
学習ログ・本・コードをベクトル化してRAG検索・推薦・進捗管理をするGoアプリ。
**TUI（bubbletea）が主インターフェース。**

---

## セッション開始時に必ずやること

**新しいセッションを始めるときは、最初に `/status` を実行してから作業を開始すること。**

これにより現在のフェーズ・残タスク・ブランチ状態が確認できる。
作業フローの詳細は [`.claude/issue_management.md`](.claude/issue_management.md) を参照。

---

## 絶対ルール

**プロジェクト方針**
- LangChainなど高レベルAIフレームワークは使わない（自力実装が目的）
- 外部ネットワーク不要・ローカルファースト
- 各コンポーネントはinterfaceで抽象化してモック差し替え可能に保つ

**コーディング（Effective Go準拠 / 詳細は `.claude/coding_rules.md`）**
- `context.Context` は全関数の第一引数
- エラーは必ず `%w` でラップして文脈を付ける
- goroutine を起動したら context キャンセルで終了を保証する
- `log` パッケージ禁止 → `slog` を使う
- マジックナンバー禁止 → 定数か config 経由
- `panic` は main/init 以外で書かない

**Git（詳細は `.claude/git_strategy.md`）**
- ブランチ: `<type>/<説明>`（例: `feat/tui-dashboard`）
- コミット: Conventional Commits形式（例: `feat(tui): add dashboard view`）
- `main` への直接 push 禁止・必ずブランチ経由PR
- **ファイルを1行も書く前にブランチを作ること**

**フォーマット**
- Goファイルを作成・編集したら、コミット前に `gofmt -w <file>` を実行すること

---

## 技術スタック（確定）

| 分類 | 技術 |
|------|------|
| 言語 | Go 1.24（golangci-lint v2互換のため） |
| TUI | bubbletea + lipgloss + bubbles + glamour |
| ベクトルDB | Qdrant（gRPC、公式`qdrant-go`クライアント） |
| ローカルLLM | Ollama（`nomic-embed-text` embedding / `llama3.2` chat） |
| ファイル監視 | fsnotify |
| 設定 | viper（`config.yaml` + env override） |
| ロガー | slog（JSON/text切り替え） |
| CLI補助 | cobra（非インタラクティブサブコマンド用） |
| CI | golangci-lint v2.1.6 / golangci-lint-action@v7 |

---

## ディレクトリ構成（重要部分）

```
cmd/magic_wand/main.go     # エントリポイント（引数なし→TUI起動）
internal/
  config/                  # viper設定
  ingestion/               # fsnotify + SHA256dedup + WorkerPool
  embedding/               # Ollama embeddingクライアント
  vectordb/                # Qdrant gRPCクライアント
  rag/                     # Retriever + Reranker(3段階) + Generator
  recommender/             # 習熟度フィルタ + exponential decay
  tui/                     # bubbletea TUI（5ビュー）
    styles/theme.go        # magic_wandカラーパレット
    views/                 # dashboard / library / query / recommend / ingest
    components/            # tabbar / statusbar / gauge / badge / table / chat
.claude/
  commands/                # Claudeスラッシュコマンド（/status, /start-phase, /close-task）
```

---

## 主要設計決定（変更時はdesign.mdの判断ログも更新）

- **進捗管理**: VectorDB payloadに`status`(unread/reading/mastered)と`mastery_level`を持たせる
- **SHA256デデュップ**: 未変更ファイルはEmbeddingをスキップ（CPU節約）
- **3段階リランカー**: `α×cosine + β×keyword_density + γ×AST_weight`（デフォルト0.6/0.3/0.1）
- **推薦フィルタ**: `status_boost`でunread=1.2 / reading=1.0 / mastered=0.5
- **TUIカラー**: Midnight背景(`#1A1B26`) / WandBlue(`#7AA2F7`) / SpellPurple(`#BB9AF7`)

---

## よく使うコマンド

```bash
make up      # Qdrant起動（docker-compose）
make down    # Qdrant停止
make run     # TUI起動
make test    # テスト実行
make lint    # golangci-lint
```

---

## ドキュメントマップ（詳細が必要なとき）

| 知りたいこと | 読むファイル |
|------------|-------------|
| **現在の進捗・次のタスク** | **GitHub Issues（`/status` で確認）** |
| **Issue管理ルール・フェーズ番号マッピング** | **`.claude/issue_management.md`** |
| ユーザーストーリー・機能要件 | `.claude/requirements.md` |
| アーキテクチャ・各コンポーネント詳細設計 | `.claude/design.md` |
| TUI画面レイアウト・カラーパレット詳細 | `.claude/design.md` §3 |
| VectorDBペイロード定義 | `.claude/design.md` §4.4 |
| RAGリランキング詳細 | `.claude/design.md` §4.5 |
| **コーディング規則（命名・エラー・goroutine等）** | `.claude/coding_rules.md` |
| **Gitブランチ戦略・コミットメッセージ** | `.claude/git_strategy.md` |
| サブCLAUDE.md の作り方 | `.claude/claude_md_guide.md` |

---

## タスク進捗（GitHub Issuesが正）

タスクの進捗は `.claude/tasks.md` ではなく **GitHub Issues** で管理する。

- 各フェーズの Epic Issue にタスク一覧がある
- 実装完了したら `/close-task <issue番号>` で閉じる
- フェーズ開始時は `/start-phase <フェーズ番号>` でブランチ作成からタスク一覧確認まで一括でできる

@.claude/tasks.md
