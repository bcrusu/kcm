package remove

import (
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
)

func RemoveNode(connection *libvirt.Connection, node repository.Node) error {
	domainName := node.Domain

	domain, err := connection.GetDomain(domainName)
	if err != nil {
		return err
	}

	if domain == nil {
		// domain does not exist
		return nil
	}

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

	if err := connection.DeleteStorageVolume(node.StoragePool, node.StorageVolume); err != nil {
		return err
	}

	return nil
}
