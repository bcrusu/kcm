package cmd

import (
	"github.com/bcrusu/kcm/cmd/list"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List clusters",
		SilenceUsage: true,
	}

	cmd.RunE = listCmdRunE
	return cmd
}

func listCmdRunE(cmd *cobra.Command, args []string) error {
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

	allClusters, err := clusterRepository.LoadAll()
	if err != nil {
		return err
	}

	list.Print(allClusters, current)

	return nil
}
