package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newIngestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ingest <path>",
		Short: "指定パスを手動インジェスト",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: T-48で実装
			fmt.Printf("ingest: %s (coming soon)\n", args[0])
			return nil
		},
	}
}

func newQueryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "query <text>",
		Short: "非インタラクティブRAG検索",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: T-49で実装
			fmt.Printf("query: %s (coming soon)\n", args[0])
			return nil
		},
	}
}

func newWatchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "watch",
		Short: "バックグラウンドウォッチャー起動",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: T-47で実装
			fmt.Println("watch: (coming soon)")
			return nil
		},
	}
}
