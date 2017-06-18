package create

import (
	"github.com/bcrusu/kcm/config"
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
)

func Node(connection *libvirt.Connection, clusterConfig *config.ClusterConfig,
	node repository.Node, networkName, sshPublicKey string) error {
	macAddresses, err := connection.GenerateUniqueMACAddresses(1)
	if err != nil {
		return err
	}

	return nodeInternal(connection, clusterConfig, node, networkName, macAddresses[0], sshPublicKey)
}

func nodeInternal(connection *libvirt.Connection, clusterConfig *config.ClusterConfig,
	node repository.Node, networkName, networkInterfaceMAC string, sshPublicKey string) error {

	storageVolume, err := connection.CreateStorageVolume(libvirt.CreateStorageVolumeParams{
		Pool:              node.StoragePool,
		Name:              node.StorageVolume,
		CapacityGiB:       node.VolumeCapacityGiB,
		BackingVolumeName: node.BackingStorageVolume,
	})
	if err != nil {
		return err
	}

	stageResult, err := clusterConfig.StageNode(node.Name, sshPublicKey)
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
		FilesystemMounts:    stageResult.FilesystemMounts,
	}

	return connection.DefineDomain(params)
}
