package remove

import (
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
)

func Cluster(connection *libvirt.Connection, cluster repository.Cluster) error {
	for _, node := range cluster.Nodes {
		if err := Node(connection, node); err != nil {
			return err
		}
	}

	return Network(connection, cluster.Network)
}
