package cmd

import "github.com/spf13/cobra"

var removeCmd = &cobra.Command{
	Use:          "remove",
	Short:        "Remove the current or specified cluster (delete all data)",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
