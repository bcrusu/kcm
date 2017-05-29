package remove

import (
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
)

func RemoveCluster(connection *libvirt.Connection, cluster repository.Cluster) error {
	for _, node := range cluster.Nodes {
		if err := RemoveNode(connection, node); err != nil {
			return err
		}
	}

	for _, node := range cluster.Masters {
		if err := RemoveNode(connection, node); err != nil {
			return err
		}
	}

	return RemoveNetwork(connection, cluster.Network)
}
