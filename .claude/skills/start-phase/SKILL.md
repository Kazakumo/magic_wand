---
name: start-phase
description: 指定フェーズの実装を開始する（ブランチ作成・Issue一覧表示）
argument-hint: <フェーズ番号 例: 2>
disable-model-invocation: false
allowed-tools: Bash, Read
---

フェーズ $ARGUMENTS の実装を開始します。以下の手順を実行してください。

1. `.claude/issue_management.md` を読んでフェーズ $ARGUMENTS の Epic Issue番号を確認する
2. `gh issue view <epic番号>` でタスク一覧を確認する
3. `git branch --show-current` で現在のブランチを確認する
4. ブランチが `main` または別フェーズのブランチの場合、`.claude/issue_management.md` のブランチ命名規則を参照して `git checkout -b <branch名>` を実行する
5. `gh issue list --label "phase$ARGUMENTS" --state open --json number,title` でタスク一覧を取得する

以下の形式で出力してください：

```
## Phase $ARGUMENTS 開始

### ブランチ
`<branch名>` を作成しました（または既存）

### Epic
#N: Epic: Phase $ARGUMENTS — <フェーズ名>

### 実装タスク（Issue番号順）
- [ ] #N: T-XX タスク名
- [ ] #N: T-XX タスク名

### 完了チェックリスト
- [ ] 全タスクIssueをclose（/close-task #N）
- [ ] make test グリーン
- [ ] make lint グリーン
- [ ] PR作成・CIグリーン確認
- [ ] mainマージ後に Epic Issueをclose
```

⚠️ ブランチを作成してから初めてファイルを編集すること。
