# Git戦略 — magic_wand

採用: **GitHub Flow** + **Conventional Commits**
最終更新: 2026-03-27

---

## 1. GitHub Flow の運用ルール

### ブランチ構成

```
main          ← 常にデプロイ可能（動くコードのみ）
  └─ feat/tui-dashboard
  └─ feat/rag-reranker
  └─ fix/worker-pool-leak
  └─ refactor/vectordb-interface
```

### フロー

```
1. mainから作業ブランチを切る
2. 小さく・頻繁にコミットする
3. PRを作成してセルフレビュー
4. CI（lint + test）がグリーンになったらmainにマージ
5. 作業ブランチを削除する
```

### ブランチ命名規則

```
<type>/<短い説明>

feat/tui-library-view
feat/ingestion-sha256-dedup
fix/goroutine-context-cancel
refactor/embedder-interface
test/rag-reranker-unit
docs/readme-setup
chore/golangci-lint-config
```

- 小文字・ハイフン区切り
- スラッシュの前は Conventional Commits の type と一致させる
- 20文字以内を目安

### マージ戦略

- **Squash merge** を基本とする（PRの複数コミットを1つに整理）
- ただしPhaseの完了など履歴に残したい節目は **Merge commit** でもよい
- `main` への直接 push は禁止（必ずブランチ経由）

---

## 2. Conventional Commits

### フォーマット

```
<type>[(<scope>)][!]: <タイトル>
<空行>
[本文]
<空行>
[フッター]
```

### Type 一覧

| type | 用途 | SemVer |
|------|------|--------|
| `feat` | 新機能追加 | MINOR |
| `fix` | バグ修正 | PATCH |
| `refactor` | 動作を変えないコード改善 | — |
| `perf` | パフォーマンス改善 | PATCH |
| `test` | テストの追加・修正 | — |
| `docs` | ドキュメントのみの変更 | — |
| `build` | ビルドシステム・依存関係の変更 | — |
| `ci` | CI設定の変更 | — |
| `chore` | その他（`.gitignore`・Makefile等） | — |

### このプロジェクトの Scope 一覧

| scope | 対象 |
|-------|------|
| `config` | `internal/config/` |
| `ingestion` | `internal/ingestion/` |
| `embedding` | `internal/embedding/` |
| `vectordb` | `internal/vectordb/` |
| `rag` | `internal/rag/` |
| `recommender` | `internal/recommender/` |
| `tui` | `internal/tui/` 全体 |
| `tui/dashboard` | `internal/tui/views/dashboard.go` |
| `tui/library` | `internal/tui/views/library.go` |
| `tui/query` | `internal/tui/views/query.go` |
| `tui/recommend` | `internal/tui/views/recommend.go` |
| `tui/ingest` | `internal/tui/views/ingest.go` |
| `tui/components` | `internal/tui/components/` |
| `cli` | `cmd/magic_wand/` |
| `deps` | `go.mod` / `go.sum` |

### 破壊的変更（BREAKING CHANGE）

```
feat(vectordb)!: add status and mastery_level to payload schema

既存のQdrantコレクションはスキーマ変更のため再作成が必要。

BREAKING CHANGE: collection "learning_chunks" must be recreated.
Run `make reset-db` before starting.
```

---

## 3. コミットメッセージの書き方

### タイトル（1行目）のルール

- 50文字以内を目安（英語）/ 日本語も可
- 命令形（現在形）で書く: `add`, `fix`, `remove`（`added`, `adding` は NG）
- 先頭を大文字にしない（`feat: add` ✅ `feat: Add` ❌）
- 末尾にピリオドを付けない

### 本文のルール

- **何を** より **なぜ** を書く
- 72文字で折り返す
- 箇条書き（`-` か `*`）でも可

### フッターのルール

```
Closes #12                                                  ← Issue参照
Refs #8                                                     ← 関連Issue
BREAKING CHANGE: <説明>                                      ← 破壊的変更
```

---

## 4. 実例集

### 基本的なfeat

```
feat(ingestion): add SHA256 dedup to skip unchanged files

ファイル内容のSHA256ハッシュをQdrantのpayloadに保存し、
再インジェスト時に変化がなければEmbeddingをスキップする。

Closes #3
```

### バグ修正

```
fix(ingestion): prevent goroutine leak on watcher shutdown

context.Done()のシグナルを監視せずgoroutineが終了しなかった。
watcherのStop()でcontextをキャンセルするよう修正。
```

### リファクタリング

```
refactor(rag): extract reranker into separate type

pipeline.goに混在していたリランキングロジックを
Rerankerインターフェースとして分離した。動作は変更しない。
```

### TUIコンポーネント

```
feat(tui/library): implement status filter tabs

All / Unread / Reading / Mastered のフィルタタブを追加。
タブ切り替えでQdrantのstatusフィールドを使ってフィルタする。
```

### 依存関係追加

```
build(deps): add charmbracelet/bubbletea and lipgloss

TUI実装のため以下を追加:
- github.com/charmbracelet/bubbletea v1.x
- github.com/charmbracelet/lipgloss v1.x
- github.com/charmbracelet/bubbles v0.x
- github.com/charmbracelet/glamour v0.x
```

### Phase完了のマージコミット

```
feat(tui): complete Phase 7 TUI implementation

5ビュー（Dashboard/Library/Query/Recommend/Ingest）と
共通コンポーネント群の実装を完了。

- T-29〜T-46 すべて完了
- magic_wandテーマ適用
- ウィンドウリサイズ対応
```

---

## 5. PR（プルリクエスト）の運用

### PRテンプレート（`.github/pull_request_template.md`）

```markdown
## 概要
<!-- 何をなぜ変更したか -->

## 変更内容
-

## タスク
- Closes #

## テスト
- [ ] 単体テスト追加・更新
- [ ] `make test` がグリーン
- [ ] `make lint` がグリーン
```

### PRのルール

- タイトルも Conventional Commits 形式にする
- 1PRは1つの論理的変更（大きすぎるPRは分割）
- Draftで作業中を示す
- Squash merge 後はブランチ削除

---

## 6. タグ・バージョニング

リリース時（v1.0.0到達時）に SemVer タグを打つ。

```bash
git tag -a v0.1.0 -m "feat: Phase 1-3 complete (foundation + embedding + vectordb)"
git tag -a v0.2.0 -m "feat: Phase 4-6 complete (ingestion + rag + recommender)"
git tag -a v1.0.0 -m "feat: Phase 7-9 complete (tui + cli + quality)"
```
