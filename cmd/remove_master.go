package cmd

import (
	"github.com/bcrusu/kcm/cmd/remove"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type removeMasterCmdState struct {
	ClusterName string
}

func newRemoveMasterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "master MASTER_NAMES",
		Aliases:      []string{"m"},
		Short:        "Remove the specified cluster master nodes",
		SilenceUsage: true,
	}

	state := &removeMasterCmdState{}
	cmd.PersistentFlags().StringVarP(&state.ClusterName, "cluster", "c", "", "Cluster that owns the master node")

	cmd.RunE = state.runE
	return cmd
}

func (s *removeMasterCmdState) runE(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("invalid command arguments - no master node specified")
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
	for _, node := range cluster.Masters {
		existing[node.Domain] = true
	}

	toRemove := make(map[string]bool)
	for _, arg := range args {
		if _, ok := existing[arg]; ok {
			toRemove[arg] = true
			continue
		}

		// try the 'short name'
		tmp := libvirtDomainName(cluster.Name, true, arg)
		if _, ok := existing[tmp]; ok {
			toRemove[tmp] = true
			continue
		}

		glog.Warningf("cluster '%s' does not contain master node '%s'", cluster.Name, arg)
	}

	if len(toRemove) == 0 {
		return nil
	}

	//TODO: do not delete last master

	connection, err := connectLibvirt()
	if err != nil {
		return err
	}
	defer connection.Close()

	for domainName := range toRemove {
		node, _ := cluster.Master(domainName)
		if err := remove.RemoveNode(connection, node); err != nil {
			return errors.Wrapf(err, "failed to remove master node '%s' in cluster '%s'", domainName, cluster.Name)
		}

		// persist the new cluster state after each node removal to minimize the potential diff beteen state in repository vs. state in libvirt
		cluster.RemoveMaster(domainName)
		if err := clusterRepository.Save(*cluster); err != nil {
			return errors.Wrapf(err, "failed to persist state for cluster '%s' after removing master node '%s'", cluster.Name, domainName)
		}
	}

	return nil
}
