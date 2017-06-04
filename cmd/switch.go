package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type switchCmdState struct {
	Clear bool
}

func newSwitchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "switch CLUSTER_NAME",
		Short:        "Switches the current cluster",
		SilenceUsage: true,
	}

	state := &switchCmdState{}
	cmd.PersistentFlags().BoolVarP(&state.Clear, "clear", "c", false, "Clears the current cluster")

	cmd.RunE = state.runE
	return cmd
}

func (s *switchCmdState) runE(cmd *cobra.Command, args []string) error {
	clusterRepository, err := newClusterRepository()
	if err != nil {
		return err
	}

	if s.Clear {
		return clusterRepository.SetCurrent("")
	}

	if len(args) != 1 {
		return errors.New("invalid command arguments")
	}

	clusterName := args[0]

	{
		current, err := clusterRepository.Current()
		if err != nil {
			return err
		}

		if current != nil && current.Name == clusterName {
			// is already the current cluster
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
