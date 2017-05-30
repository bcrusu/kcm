package create

import (
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
	"github.com/pkg/errors"
)

func ValidateLibvirtObjects(connection *libvirt.Connection, cluster repository.Cluster) error {
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

	checkNode := func(node repository.Node) error {
		domain, err := connection.GetDomain(node.Domain)
		if err != nil {
			return err
		}
		if domain != nil {
			return errors.Errorf("validation: libvirt domain '%s' already exists", node.Domain)
		}

		storageVolume, err := connection.GetStorageVolume(cluster.StoragePool, node.StorageVolume)
		if err != nil {
			return err
		}
		if storageVolume != nil {
			return errors.Errorf("validation: libvirt storage volume '%s' already exists", node.StorageVolume)
		}

		return nil
	}

	for _, node := range cluster.Masters {
		if err := checkNode(node); err != nil {
			return err
		}
	}

	for _, node := range cluster.Nodes {
		if err := checkNode(node); err != nil {
			return err
		}
	}

	return nil
}

func CreateLibvirtObjects(connection *libvirt.Connection, cluster repository.Cluster) error {
	if err := defineLibvirtNetwork(connection, cluster.Network); err != nil {
		return err
	}

	createDomain := func(node repository.Node, networkInterfaceMAC string) error {
		if err := connection.CreateStorageVolume(node.StoragePool, node.StorageVolume, cluster.BackingStorageVolume); err != nil {
			return err
		}

		if err := defineLibvirtDomain(connection, node, cluster.Network.Name, networkInterfaceMAC); err != nil {
			return err
		}

		return nil
	}

	macAddresses, err := connection.GenerateUniqueMACAddresses(len(cluster.Masters) + len(cluster.Nodes))
	if err != nil {
		return err
	}

	for i, node := range cluster.Masters {
		if err := createDomain(node, macAddresses[i]); err != nil {
			return err
		}
	}

	for i, node := range cluster.Nodes {
		macAddress := macAddresses[i+len(cluster.Masters)]
		if err := createDomain(node, macAddress); err != nil {
			return err
		}
	}
	return nil
}

func defineLibvirtNetwork(connection *libvirt.Connection, network repository.Network) error {
	params := libvirt.DefineNetworkParams{
		Name:     network.Name,
		IPv4CIDR: network.IPv4CIDR,
		IPv6CIDR: network.IPv6CIDR,
	}

	return connection.DefineNATNetwork(params)
}

func defineLibvirtDomain(connection *libvirt.Connection, node repository.Node, networkName, networkInterfaceMAC string) error {
	storageVolume, err := connection.GetStorageVolume(node.StoragePool, node.StorageVolume)
	if err != nil {
		return err
	}

	params := libvirt.DefineDomainParams{
		Name:                node.Domain,
		Network:             networkName,
		NetworkInterfaceMAC: networkInterfaceMAC,
		StorageVolumePath:   storageVolume.Target().Path(),
		MemoryMiB:           node.MemoryMiB,
		CPUs:                node.CPUs,
		//TODO: FilesystemMounts    :
	}

	return connection.DefineDomain(params)
}
