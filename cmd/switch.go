package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newSwitchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "switch CLUSTER_NAME",
		Short:        "Switch of current (active) cluster",
		SilenceUsage: true,
	}

	cmd.RunE = switchCmdRunE
	return cmd
}

func switchCmdRunE(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("invalid command arguments")
	}

	clusterName := args[0]

	clusterRepository, err := newClusterRepository()
	if err != nil {
		return err
	}

	{
		current, err := clusterRepository.Current()
		if err != nil {
			return err
		}

		if current != nil && current.Name == clusterName {
			// is already the active cluster
			return nil
		}
	}

	cluster, err := clusterRepository.Load(clusterName)
	if err != nil {
		return err
	}

	if cluster == nil {
		return errors.Errorf("could not find cluster '%s'", clusterName)
	}

	return clusterRepository.SetCurrent(cluster.Name)
}
