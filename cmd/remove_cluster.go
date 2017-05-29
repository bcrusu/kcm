package cmd

import (
	"github.com/bcrusu/kcm/cmd/remove"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newRemoveClusterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "cluster",
		Aliases:      []string{"cluster", "c"},
		Short:        "Remove the specified clusters",
		SilenceUsage: true,
	}

	cmd.RunE = runRemoveClusterCmdE
	return cmd
}

func runRemoveClusterCmdE(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("specify the cluster names to remove")
	}

	connection, err := connectLibvirt()
	if err != nil {
		return err
	}
	defer connection.Close()

	clusterRepository, err := newClusterRepository()
	if err != nil {
		return err
	}

	for _, name := range args {
		cluster, err := clusterRepository.Load(name)
		if err != nil {
			return err
		}

		if cluster == nil {
			glog.Warningf("could not find cluster '%s'", name)
			return nil
		}

		if err := remove.RemoveCluster(connection, *cluster); err != nil {
			return errors.Wrapf(err, "failed to remove cluster libvirt objects '%s'", name)
		}

		if err := clusterRepository.Remove(cluster.Name); err != nil {
			return errors.Wrapf(err, "failed to remove cluster data '%s'", name)
		}
	}

	return nil
}
