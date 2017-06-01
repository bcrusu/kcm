package cmd

import (
	"github.com/bcrusu/kcm/cmd/remove"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newRemoveClusterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "cluster CLUSTER_NAME",
		Aliases:      []string{"c"},
		Short:        "Remove the specified clusters",
		SilenceUsage: true,
	}

	cmd.RunE = removeClusterCmdRunE
	return cmd
}

func removeClusterCmdRunE(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("invalid command arguments")
	}

	clusterName := args[0]

	clusterRepository, err := newClusterRepository()
	if err != nil {
		return err
	}

	cluster, err := clusterRepository.Load(clusterName)
	if err != nil {
		return err
	}

	if cluster == nil {
		return errors.Errorf("could not find cluster '%s'", clusterName)
	}

	connection, err := connectLibvirt()
	if err != nil {
		return err
	}
	defer connection.Close()

	clusterConfig, err := getClusterConfig(*cluster)
	if err != nil {
		return err
	}

	if err := remove.Cluster(connection, clusterConfig, *cluster); err != nil {
		return errors.Wrapf(err, "failed to remove cluster libvirt objects '%s'", clusterName)
	}

	if err := clusterRepository.Remove(cluster.Name); err != nil {
		return errors.Wrapf(err, "failed to remove cluster data '%s'", clusterName)
	}

	return nil
}
