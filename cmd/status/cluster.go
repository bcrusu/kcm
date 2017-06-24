package status

import (
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
)

type ClusterStatus struct {
	Active  bool
	Network NetworkStatus
	Nodes   []NodeStatus
	Cluster repository.Cluster
}

func Cluster(connection *libvirt.Connection, cluster repository.Cluster) (*ClusterStatus, error) {
	netStatus, err := Network(connection, cluster.Network)
	if err != nil {
		return nil, err
	}

	var nodes []NodeStatus
	for _, node := range cluster.Nodes {
		nodeStatus, err := Node(connection, node)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, *nodeStatus)
	}

	return &ClusterStatus{
		Active:  netStatus.Active,
		Network: *netStatus,
		Nodes:   nodes,
		Cluster: cluster,
	}, nil
}

func IsClusterActive(connection *libvirt.Connection, cluster repository.Cluster) (bool, error) {
	// simple check atm. - assume cluster is running if the network is active
	netStatus, err := Network(connection, cluster.Network)
	if err != nil {
		return false, err
	}

	return netStatus.Active, nil
}
