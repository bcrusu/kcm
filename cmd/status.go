package cmd

import (
	"fmt"

	"strings"

	"github.com/bcrusu/kcm/cmd/status"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:          "status [CLUSTER_NAME]",
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

	fmt.Printf("Cluster %s\n", cluster.Name)

	{
		fmt.Println("Nodes:")
		for name, node := range clusterStatus.Nodes {
			fmt.Printf("%s:\t", name)
			if node.Missing {
				fmt.Print("missing")
			} else if node.Active {
				fmt.Printf("active")
				if len(node.InterfaceAddresses) > 0 {
					fmt.Printf(" (%s)", strings.Join(node.InterfaceAddresses, ", "))
				}
			} else {
				fmt.Print("inactive")
			}

			fmt.Println()
		}
	}

	return nil
}
