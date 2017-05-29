package cmd

import (
	"github.com/bcrusu/kcm/cmd/remove"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type removeNodeCmdState struct {
	ClusterName string
}

func newRemoveNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "node NODE_NAMES",
		Aliases:      []string{"n"},
		Short:        "Remove the specified cluster nodes",
		SilenceUsage: true,
	}

	state := &removeNodeCmdState{}
	cmd.PersistentFlags().StringVarP(&state.ClusterName, "cluster", "c", "", "Cluster that owns the node")

	cmd.RunE = state.runE
	return cmd
}

func (s *removeNodeCmdState) runE(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("invalid command arguments - no node specified")
	}

	clusterRepository, err := newClusterRepository()
	if err != nil {
		return err
	}

	cluster, err := getWorkingCluster(clusterRepository, s.ClusterName)
	if err != nil {
		return err
	}

	existing := make(map[string]bool)
	for _, node := range cluster.Nodes {
		existing[node.Domain] = true
	}

	toRemove := make(map[string]bool)
	for _, arg := range args {
		if _, ok := existing[arg]; ok {
			toRemove[arg] = true
			continue
		}

		// try the 'short name'
		tmp := libvirtDomainName(cluster.Name, false, arg)
		if _, ok := existing[tmp]; ok {
			toRemove[tmp] = true
			continue
		}

		glog.Warningf("cluster '%s' does not contain node '%s'", cluster.Name, arg)
	}

	if len(toRemove) == 0 {
		return nil
	}

	connection, err := connectLibvirt()
	if err != nil {
		return err
	}
	defer connection.Close()

	for domainName := range toRemove {
		node, _ := cluster.Node(domainName)
		if err := remove.RemoveNode(connection, node); err != nil {
			return errors.Wrapf(err, "failed to remove node '%s' in cluster '%s'", domainName, cluster.Name)
		}

		// persist the new cluster state after each node removal to minimize the potential diff beteen state in repository vs. state in libvirt
		cluster.RemoveNode(domainName)
		if err := clusterRepository.Save(*cluster); err != nil {
			return errors.Wrapf(err, "failed to persist state for cluster '%s' after removing node '%s'", cluster.Name, domainName)
		}
	}

	return nil
}
