package create

import (
	"github.com/bcrusu/kcm/config"
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
)

func Cluster(connection *libvirt.Connection, clusterConfig *config.ClusterConfig,
	cluster repository.Cluster) error {
	if err := Network(connection, cluster.Network); err != nil {
		return err
	}

	macAddresses, err := connection.GenerateUniqueMACAddresses(len(cluster.Nodes))
	if err != nil {
		return err
	}

	index := 0
	for _, node := range cluster.Nodes {
		if err := Node(connection, clusterConfig, node, cluster.Network.Name, macAddresses[index]); err != nil {
			return err
		}

		index++
	}

	return nil
}
