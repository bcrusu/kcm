package remove

import (
	"github.com/bcrusu/kcm/config"
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
)

func Node(connection *libvirt.Connection, clusterConfig *config.ClusterConfig, node repository.Node) error {
	{
		domainName := node.Domain
		domain, err := connection.GetDomain(domainName)
		if err != nil {
			return err
		}

		if domain != nil {
			active, err := connection.DomainIsActive(domainName)
			if err != nil {
				return err
			}

			if active {
				if err := connection.DestroyDomain(domainName); err != nil {
					return err
				}
			}

			if err := connection.UndefineDomain(domainName); err != nil {
				return err
			}
		}
	}

	{
		storageVolume, err := connection.GetStorageVolume(node.StoragePool, node.StorageVolume)
		if err != nil {
			return err
		}

		if storageVolume != nil {
			if err := connection.DeleteStorageVolume(node.StoragePool, node.StorageVolume); err != nil {
				return err
			}
		}
	}

	if err := clusterConfig.UnstageNode(node.Name); err != nil {
		return err
	}

	return nil
}
