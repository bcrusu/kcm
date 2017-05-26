package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(createCmd)
	RootCmd.AddCommand(removeCmd)
	RootCmd.AddCommand(switchCmd)
}

var RootCmd = &cobra.Command{
	Use:          "kcm",
	SilenceUsage: true,
}
