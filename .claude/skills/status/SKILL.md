---
name: status
description: magic_wandプロジェクトの現在状態（フェーズ・残タスク・ブランチ）を確認する
disable-model-invocation: false
allowed-tools: Bash
---

以下の手順でプロジェクトの現在状態を確認し、レポートしてください。

1. `gh issue list --label epic --state open --json number,title` で進行中フェーズを確認
2. `gh issue list --label epic --state closed --json number,title` で完了済みフェーズを確認
3. open な epic の phaseラベルを見て、現フェーズの task issueを `gh issue list --label task --state open --json number,title,labels` で確認
4. `git branch --show-current` で現在のブランチを確認

以下の形式でレポートしてください：

```
## プロジェクト状態

### フェーズ進捗
- ✅ 完了: Phase N — <名前> (#Issue番号)
- 🔄 進行中: Phase N — <名前> (#Issue番号)  ← openのepicがあれば
- ⬜ 未着手: Phase N+1〜

### 次のタスク
- #N: タスク名
- #N: タスク名

### ブランチ
現在: <branch名>
```

現在のブランチが `main` の場合は「⚠️ 実装前に `/start-phase <N>` でブランチを作成してください」と警告を出すこと。
