引数: フェーズ番号（例: `/start-phase 2`）

指定されたフェーズの実装を開始するための準備をしてください。

## 手順

1. `.claude/issue_management.md` を読んでEpic Issue番号を確認する
2. 対象フェーズのEpic Issueを `gh issue view <epic番号>` で確認し、タスク一覧を把握する
3. `git status` と `git branch` で現在のブランチを確認する
4. mainブランチにいる場合、`.claude/git_strategy.md` を参照して適切なブランチ名を決定し、`git checkout -b <branch名>` で作業ブランチを作成する
5. 対象フェーズの全タスクIssueを `gh issue list --label "phase$ARGUMENTS" --state open` で一覧表示する

## 出力形式

```
## Phase $ARGUMENTS 開始準備

### 作業ブランチ
`feat/phase$ARGUMENTS-<name>` を作成しました

### Epicチケット
#X: Epic: Phase $ARGUMENTS — <フェーズ名>

### 実装対象タスク（Issue番号順）
- [ ] #N: T-XX: タスク名
- [ ] #N: T-XX: タスク名
...

### 実装順序の推奨
<依存関係がある場合は順序を明示>

### 完了時のチェックリスト
- [ ] 全タスクIssueをcloseした
- [ ] make test がグリーン
- [ ] make lint がグリーン
- [ ] PRを作成してCI確認
- [ ] Epic Issueをcloseした
- [ ] mainにマージした
```

注意: ブランチを作成する前にファイルを編集しないでください。
