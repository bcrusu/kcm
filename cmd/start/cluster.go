package start

import (
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
)

func Cluster(connection *libvirt.Connection, cluster repository.Cluster) error {
	if err := Network(connection, cluster.Network); err != nil {
		return err
	}

	// start masters first
	for _, node := range cluster.Nodes {
		if !node.IsMaster {
			continue
		}

		if err := Node(connection, node); err != nil {
			return err
		}
	}

	for _, node := range cluster.Nodes {
		if node.IsMaster {
			continue
		}

		if err := Node(connection, node); err != nil {
			return err
		}
	}

	return nil
}

func IsClusterRunning(connection *libvirt.Connection, cluster repository.Cluster) (bool, error) {
	// simple check atm. - assume cluster is running if the network is active
	return isNetworkActive(connection, cluster.Network.Name)
}
