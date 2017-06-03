package cmd

import (
	"fmt"

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

	for _, cluster := range allClusters {
		mark := " "
		if current != nil && cluster.Name == current.Name {
			mark = "*"
		}

		fmt.Printf("%s%s\n", mark, cluster.Name)
	}

	return nil
}
