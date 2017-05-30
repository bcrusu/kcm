package start

import (
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
)

func Cluster(connection *libvirt.Connection, cluster repository.Cluster) error {
	if err := Network(connection, cluster.Network); err != nil {
		return err
	}

	for _, node := range cluster.Masters {
		if err := Node(connection, node); err != nil {
			return err
		}
	}

	for _, node := range cluster.Nodes {
		if err := Node(connection, node); err != nil {
			return err
		}
	}

	return nil
}
