package create

import (
	"github.com/bcrusu/kcm/config"
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
)

func Node(connection *libvirt.Connection, clusterConfig *config.ClusterConfig,
	node repository.Node, networkName, networkInterfaceMAC string) error {

	storageVolume, err := connection.CreateStorageVolume(node.StoragePool, node.StorageVolume, node.BackingStorageVolume)
	if err != nil {
		return err
	}

	stageResult, err := clusterConfig.StageNode(node.Name)
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
