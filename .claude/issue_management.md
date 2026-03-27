# Issue管理ルール — magic_wand

## 構造

タスクは **GitHub Issues** で管理する（`.claude/tasks.md` はマスターリストとして参照用に保持するが、進捗管理はIssuesが正）。

```
Epic Issue（フェーズ単位）
  └─ Task Issue（個別タスク）
  └─ Task Issue
  └─ Task Issue
```

### Epic Issue
- 1フェーズ = 1 Epic Issue
- ラベル: `epic` + `phaseN`（例: `epic,phase2`）
- 本文にタスクのチェックリストを記載（子Issueリンクではなくテキストチェックボックス）

### Task Issue
- 1タスク = 1 Issue
- ラベル: `task` + `phaseN`
- 本文の先頭に `Parent: #<epic番号>` を記載
- 関連する他タスクがある場合は `Requires: #<issue番号>` を記載

## ラベル一覧

| ラベル | 用途 |
|-------|------|
| `epic` | フェーズ単位の親チケット |
| `task` | 個別実装タスク |
| `phase1`〜`phase9` | フェーズ識別 |

## Issue番号マッピング（初期）

| Issue | 内容 |
|-------|------|
| #2 | Epic: Phase 1（完了済み・closed） |
| #9 | Epic: Phase 2 — Embedding Layer |
| #13 | Epic: Phase 3 — Vector DB Layer |
| #19 | Epic: Phase 4 — Ingestion Pipeline |
| #20 | Epic: Phase 5 — RAG Engine |
| #21 | Epic: Phase 6 — Semantic Recommender |
| #22 | Epic: Phase 7 — TUI実装 |
| #23 | Epic: Phase 8 — CLI非インタラクティブ |
| #24 | Epic: Phase 9 — 品質・仕上げ |

## Claudeスラッシュコマンド（`.claude/commands/`）

| コマンド | 説明 |
|---------|------|
| `/status` | 現在のフェーズ・残タスク・ブランチ状態を一覧表示 |
| `/start-phase <N>` | フェーズNの作業開始（ブランチ作成・Issue一覧表示） |
| `/close-task <Issue番号>` | タスク完了処理（Issueclose・Epic進捗更新） |

## 実装フロー（セッションをまたぐ場合）

```
セッション開始
  └─ /status                          # 現在状態を確認
  └─ /start-phase N                   # ブランチ作成 + タスク確認
      └─ 実装（各タスクをIssue番号で追跡）
      └─ /close-task #X               # タスク完了ごとに実行
      └─ make test && make lint
      └─ PRを作成（Closes #X, #Y...）
      └─ CIグリーン確認
      └─ mainにマージ
      └─ /close-task <epic番号>       # フェーズ完了
```

## ブランチ命名規則（フェーズ対応）

| フェーズ | ブランチ名 |
|---------|----------|
| Phase 2 | `feat/phase2-embedding` |
| Phase 3 | `feat/phase3-vectordb` |
| Phase 4 | `feat/phase4-ingestion` |
| Phase 5 | `feat/phase5-rag` |
| Phase 6 | `feat/phase6-recommender` |
| Phase 7 | `feat/phase7-tui` |
| Phase 8 | `feat/phase8-cli` |
| Phase 9 | `feat/phase9-quality` |

## PRとIssueの紐付け

PRの本文に `Closes #<issue番号>` を記載すると、マージ時に自動でcloseされる。

```markdown
## タスク
- Closes #10 (T-07)
- Closes #11 (T-08)
- Closes #12 (T-09)
- Closes #9 (Epic: Phase 2)
```
