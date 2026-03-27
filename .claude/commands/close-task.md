引数: Issue番号（例: `/close-task 10`）

タスクの完了処理をしてください。

## 手順

1. `gh issue view $ARGUMENTS` で対象Issueの内容を確認する
2. `gh issue close $ARGUMENTS -c "実装完了"` でIssueをcloseする
3. Issueの本文から `Parent: #<epic番号>` を読み取る
4. `gh issue view <epic番号>` でEpic Issueを確認し、チェックリストの該当タスクにチェックを入れる（`gh issue edit` で更新）
5. 同じEpicの子タスクが全てcloseされているか `gh issue list --label "task" --state open` で確認する
6. 全タスク完了していたら「Epic #X の全タスクが完了しました。PRをマージ後に `/close-task <epic番号>` でEpicをcloseしてください」と通知する

## 注意

- Issueをcloseするのは実装・テストが完了した後のみ
- Epicをcloseするのは PRがマージされた後
