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

## 実装開始時のフロー

1. 対象フェーズのEpic Issueを確認（子タスクの番号を把握）
2. `main` から作業ブランチを切る（`feat/phase2-embedding` 等）
3. 各タスク完了時に対応するTask Issueをcloseする
4. フェーズ完了・PRマージ後にEpic Issueをcloseする

## PRとIssueの紐付け

PRの本文に `Closes #<issue番号>` を記載すると、マージ時に自動でcloseされる。

```markdown
## タスク
- Closes #10 (T-07)
- Closes #11 (T-08)
- Closes #12 (T-09)
- Closes #9 (Epic: Phase 2)
```
