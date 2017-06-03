package remove

import (
	"github.com/bcrusu/kcm/config"
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
)

func Cluster(connection *libvirt.Connection, clusterConfig *config.ClusterConfig, cluster repository.Cluster) error {
	for _, node := range cluster.Nodes {
		if err := Node(connection, clusterConfig, node); err != nil {
			return err
		}
	}

	if err := Network(connection, cluster.Network); err != nil {
		return err
	}

	return clusterConfig.UnstageCluster()
}
