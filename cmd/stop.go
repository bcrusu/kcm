package cmd

import (
	"github.com/bcrusu/kcm/cmd/stop"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type stopCmdState struct {
	Force bool
}

func newStopCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:          "stop [CLUSTER_NAME]",
		Short:        "Stops the specified/current cluster",
		SilenceUsage: true,
	}

	state := &stopCmdState{}
	cmd.PersistentFlags().BoolVarP(&state.Force, "force", "f", false, "Does not use graceful shutdown (may produce inconsistent storage volume state)")

	cmd.RunE = state.runE
	return cmd
}

func (s stopCmdState) runE(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return errors.New("invalid command arguments")
	}

	clusterRepository, err := newClusterRepository()
	if err != nil {
		return err
	}

	clusterName := ""
	if len(args) == 1 {
		clusterName = args[0]
	}

	cluster, err := getWorkingCluster(clusterRepository, clusterName)
	if err != nil {
		return err
	}

	connection, err := connectLibvirt()
	if err != nil {
		return err
	}
	defer connection.Close()

	if err := stop.Cluster(connection, *cluster, s.Force); err != nil {
		return errors.Wrapf(err, "failed to stop cluster '%s'", cluster.Name)
	}

	return nil
}
