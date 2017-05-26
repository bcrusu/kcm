package cmd

import "github.com/spf13/cobra"

var createCmd = &cobra.Command{
	Use:          "create",
	Short:        "Create a new cluster",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
