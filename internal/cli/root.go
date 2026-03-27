package cli

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "magic_wand",
		Short: "ローカルファーストのAI学習オーケストレーター",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTUI(cmd.Context())
		},
	}

	root.AddCommand(
		newIngestCmd(),
		newQueryCmd(),
		newWatchCmd(),
	)

	return root
}
