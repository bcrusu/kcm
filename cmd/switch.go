package cmd

import "github.com/spf13/cobra"

var switchCmd = &cobra.Command{
	Use:          "switch",
	Short:        "Switch of current (active) cluster",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

//TODO
