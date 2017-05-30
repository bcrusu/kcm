package cmd

import (
	"github.com/bcrusu/kcm/cmd/start"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newStartCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:          "start [CLUSTER_NAME]",
		Short:        "Starts the specified/current cluster",
		SilenceUsage: true,
	}

	cmd.RunE = startCmdRunE
	return cmd
}

func startCmdRunE(cmd *cobra.Command, args []string) error {
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

	if err := start.Cluster(connection, *cluster); err != nil {
		return errors.Wrapf(err, "failed to start cluster '%s'", cluster.Name)
	}

	return nil
}
