# サブディレクトリ CLAUDE.md 作成ガイド

## いつ作るか

モジュールの実装を開始するとき（Phaseが始まるとき）に作成する。
そのディレクトリ以下で作業するセッションでは自動でロードされる。

## どこに置くか・何を書くか

### `internal/tui/CLAUDE.md` — TUI実装時
```markdown
# TUI モジュール

## アーキテクチャ
Elm-likeアーキテクチャ（bubbletea）。Model/Update/View の3関数。
ルートモデル(model.go)がタブ管理。各ビューは独立したサブモデル。

## カラーパレット（theme.goから使う）
styles.Theme.Primary / Accent / Success / Warning / Error / Border / Subtle

## 注意点
- tea.Cmd は副作用の唯一の出口。goroutineを直接Updateで起動しない
- ウィンドウリサイズは tea.WindowSizeMsg で全ビューに伝播させる
- IngestEventはチャンネルをtea.Cmdでポーリングして受け取る（`waitForIngestEvent`パターン）
```

### `internal/rag/CLAUDE.md` — RAG実装時
```markdown
# RAG モジュール

## パイプライン順序
Retriever(top20) → Reranker(3段階) → top5 → Generator(streaming)

## リランキング係数（config.yaml で上書き可能）
α=0.6(cosine) β=0.3(keyword_density) γ=0.1(AST_weight)

## Generatorのストリーミング
`chan string` で返す。呼び出し元がrangeで受け取りTUIに送る。
```

### `internal/vectordb/CLAUDE.md` — VectorDB実装時
```markdown
# VectorDB モジュール

## ペイロードの必須フィールド
path / filename / ext / file_hash / chunk_index / total_chunks / content
created_at / updated_at / status / mastery_level

## インデックス
path(Keyword) / status(Keyword) / updated_at(Numeric) の3つが必須

## statusの値
"unread" | "reading" | "mastered"（それ以外は受け付けない）
```

## tasks.md の更新規則

セッション開始時: tasks.md を読んで `[~]` の作業を把握する
作業完了時: 必ず `[x]` に更新してコミットする
次フェーズ開始前: tasks.md の `実装順序の推奨` を確認する
