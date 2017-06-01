package cmd

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/bcrusu/kcm/config"
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
	"github.com/pkg/errors"
)

func newClusterRepository() (repository.ClusterRepository, error) {
	repoPath := path.Join(*dataDir, "repository")
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

func libvirtDomainName(clusterName string, nodeName string) string {
	return fmt.Sprintf("kcm.%s.%s", clusterName, nodeName)
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

func getClusterConfig(cluster repository.Cluster, sshPublicKeyPath string) (*config.ClusterConfig, error) {
	sshPublicKey, err := readSSHPublicKey(sshPublicKeyPath)
	if err != nil {
		return nil, err
	}

	configDir := path.Join(*dataDir, "config", cluster.Name)
	return config.New(configDir, cluster, sshPublicKey)
}

func readSSHPublicKey(path string) (string, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "cannot load SSH public key from file '%s'", path)
	}

	return string(bytes), nil
}
