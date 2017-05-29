package cmd

import "github.com/spf13/cobra"

func newRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "remove",
		Aliases:      []string{"rm"},
		Short:        "Removes the specified object (deletes all data)",
		SilenceUsage: true,
	}

	cmd.AddCommand(newRemoveClusterCmd())
	cmd.AddCommand(newRemoveMasterCmd())
	cmd.AddCommand(newRemoveNodeCmd())

	return cmd
}
