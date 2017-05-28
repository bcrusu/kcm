package cmd

import (
	"fmt"
	"path"

	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
)

func newClusterRepository() (repository.ClusterRepository, error) {
	repoPath := path.Join(*dataDir, "clusters")
	return repository.New(repoPath)
}

func kubernetesCacheDir() string {
	return path.Join(*dataDir, "kubernetes")
}

func connectLibvirt() (*libvirt.Connection, error) {
	return libvirt.NewConnection(*libvirtURI)
}

func libvirtNetworkName(clusterName string) string {
	return fmt.Sprintf("kcm.%s", clusterName)
}

func libvirtStorageVolumeName(domainName string) string {
	return fmt.Sprintf("%s.qcow2", domainName)
}

func libvirtDomainName(clusterName string, isMaster bool, number uint) string {
	nodeType := "node"
	if isMaster {
		nodeType = "master"
	}

	return fmt.Sprintf("kcm.%s.%s.%d", clusterName, nodeType, number)
}

func coreOSStorageVolumeName(version string) string {
	return fmt.Sprintf("coreos_production_qemu_image_%s.qcow2", version)
}
