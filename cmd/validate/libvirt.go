package validate

import (
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
	"github.com/pkg/errors"
)

func LibvirtClusterObjects(connection *libvirt.Connection, cluster repository.Cluster) error {
	storagePool, err := connection.GetStoragePool(cluster.StoragePool)
	if err != nil {
		return err
	}
	if storagePool == nil {
		return errors.Errorf("validation: libvirt storage pool '%s' does not exist", cluster.StoragePool)
	}

	network, err := connection.GetNetwork(cluster.Network.Name)
	if err != nil {
		return err
	}
	if network != nil {
		return errors.Errorf("validation: libvirt network '%s' already exists", cluster.Network.Name)
	}

	for _, node := range cluster.Nodes {
		if err := LibvirtNodeObjects(connection, node); err != nil {
			return err
		}
	}

	return nil
}

func LibvirtNodeObjects(connection *libvirt.Connection, node repository.Node) error {
	domain, err := connection.GetDomain(node.Domain)
	if err != nil {
		return err
	}
	if domain != nil {
		return errors.Errorf("validation: libvirt domain '%s' already exists", node.Domain)
	}

	storageVolume, err := connection.GetStorageVolume(node.StoragePool, node.StorageVolume)
	if err != nil {
		return err
	}
	if storageVolume != nil {
		return errors.Errorf("validation: libvirt storage volume '%s' already exists", node.StorageVolume)
	}

	return nil
}
