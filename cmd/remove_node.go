package cmd

import (
	"github.com/bcrusu/kcm/cmd/remove"
	"github.com/bcrusu/kcm/repository"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type removeNodeCmdState struct {
	ClusterName string
	IsMaster    bool
}

func newRemoveNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "node NODE_NAME",
		Aliases:      []string{"n"},
		Short:        "Remove the specified cluster node",
		SilenceUsage: true,
	}

	state := &removeNodeCmdState{}
	cmd.PersistentFlags().StringVarP(&state.ClusterName, "cluster", "c", "", "Cluster that owns the node")
	cmd.PersistentFlags().BoolVarP(&state.IsMaster, "master", "m", false, "Remove master node")

	cmd.RunE = state.runE
	return cmd
}

func (s *removeNodeCmdState) runE(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("invalid command arguments")
	}

	nodeName := args[0]

	clusterRepository, err := newClusterRepository()
	if err != nil {
		return err
	}

	cluster, err := getWorkingCluster(clusterRepository, s.ClusterName)
	if err != nil {
		return err
	}

	var toRemove *repository.Node
	for _, node := range cluster.Nodes {
		if nodeName == node.Name {
			toRemove = &node
			break
		}
	}

	if toRemove == nil {
		glog.Warningf("cluster '%s' does not contain node '%s'", cluster.Name, nodeName)
		return nil
	}

	//TODO: do not allow to remove last node?

	connection, err := connectLibvirt()
	if err != nil {
		return err
	}
	defer connection.Close()

	clusterConfig, err := getClusterConfig(*cluster)
	if err != nil {
		return err
	}

	if err := remove.Node(connection, clusterConfig, *toRemove); err != nil {
		return errors.Wrapf(err, "failed to remove node '%s' in cluster '%s'", nodeName, cluster.Name)
	}

	delete(cluster.Nodes, toRemove.Name)
	if err := clusterRepository.Save(*cluster); err != nil {
		return errors.Wrapf(err, "failed to persist state for cluster '%s' after removing node '%s'", cluster.Name, nodeName)
	}

	return nil
}
