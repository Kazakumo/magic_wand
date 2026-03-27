---
name: close-task
description: タスクIssueを完了処理する（close + Epic進捗確認）
argument-hint: <Issue番号 例: 10>
disable-model-invocation: false
allowed-tools: Bash
---

Issue #$ARGUMENTS の完了処理をします。

1. `gh issue view $ARGUMENTS` で対象Issueの内容（タイトル・Parent番号）を確認する
2. `gh issue close $ARGUMENTS --comment "実装完了"` でIssueをcloseする
3. Issueの本文から `Parent: #<epic番号>` を読み取る
4. `gh issue list --label task --state open --json number,title` で同じフェーズのtaskが残っているか確認する
5. 残タスクがある場合は残数を表示する。全タスクが完了している場合は以下を通知する：

```
✅ Phase N の全タスクが完了しました。
PRをmainにマージした後、Epic Issue #X を /close-task <epic番号> でcloseしてください。
```

⚠️ Issueをcloseするのは実装・テストが完了した後のみ。Epicをcloseするのは PRマージ後。
