package cmd

import (
	"fmt"
	"path"

	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
	"github.com/pkg/errors"
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

func libvirtDomainName(clusterName string, isMaster bool, shortName string) string {
	nodeType := "node"
	if isMaster {
		nodeType = "master"
	}

	return fmt.Sprintf("kcm.%s.%s.%s", clusterName, nodeType, shortName)
}

func coreOSStorageVolumeName(version string) string {
	return fmt.Sprintf("coreos_production_qemu_image_%s.qcow2", version)
}

func getWorkingCluster(clusterRepository repository.ClusterRepository, clusterName string) (*repository.Cluster, error) {
	var cluster *repository.Cluster
	var err error

	if clusterName != "" {
		cluster, err = clusterRepository.Load(clusterName)
	} else {
		cluster, err = clusterRepository.Current()
	}

	if err != nil {
		return nil, err
	}

	if cluster == nil {
		return nil, errors.Errorf("could not determine working cluster")
	}

	return cluster, nil
}
