以下の手順でプロジェクトの現在状態を確認し、簡潔にレポートしてください。

## 手順

1. `gh issue list --label epic --state open` を実行して、現在進行中のフェーズを確認する
2. `gh issue list --label epic --state closed` を実行して、完了済みフェーズを確認する
3. 直近のopen epicに対して `gh issue list --label task --state open` で残タスクを確認する
4. `git branch` で現在のブランチを確認する
5. `git log --oneline -5` で最近のコミットを確認する

## レポート形式

```
## 現在のフェーズ
- 完了済み: Phase N (Epic #X)
- 進行中: Phase N (Epic #X) — 残タスク M件
- 未着手: Phase N+1〜

## 次にやること
- Issue #X: タスク名
- Issue #Y: タスク名

## ブランチ状態
- 現在: <branch名>
- 作業ブランチが必要なら: feat/phaseN-<name>
```

作業ブランチがまだmainなら、実装前にブランチを切るよう明示的に指摘してください。
