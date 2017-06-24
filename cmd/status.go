package cmd

import (
	"fmt"

	"github.com/bcrusu/kcm/cmd/status"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:          "status [CLUSTER_NAME]",
		Aliases:      []string{"stat"},
		Short:        "Prints the status for the specified/current cluster",
		SilenceUsage: true,
	}

	cmd.RunE = statusCmdRunE
	return cmd
}

func statusCmdRunE(cmd *cobra.Command, args []string) error {
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

	clusterStatus, err := status.Cluster(connection, *cluster)
	if err != nil {
		return errors.Wrapf(err, "failed to get status for cluster '%s'", cluster.Name)
	}

	status.PrintCluster(*clusterStatus)
	fmt.Println()

	status.PrintNetwork(clusterStatus.Network)
	fmt.Println()

	status.PrintNodes(clusterStatus.Nodes)

	return nil
}
