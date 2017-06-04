package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newCurrentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "current",
		Short:        "Prints the current cluster",
		SilenceUsage: true,
	}

	cmd.RunE = currentCmdRunE
	return cmd
}

func currentCmdRunE(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return errors.New("invalid command arguments")
	}

	clusterRepository, err := newClusterRepository()
	if err != nil {
		return err
	}

	current, err := clusterRepository.Current()
	if err != nil {
		return err
	}

	if current == nil {
		return nil
	}

	fmt.Println(current.Name)

	return nil
}
