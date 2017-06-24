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

const MasterNodeNamePrefix = "master"
const NodeNamePrefix = "node"

func newClusterRepository() (repository.ClusterRepository, error) {
	repoPath := path.Join(*dataDir, "clusters")
	return repository.New(repoPath)
}

func cacheDir() string {
	return path.Join(*dataDir, "cache")
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
		if clusterName != "" {
			return nil, errors.Errorf("could not find cluster '%s'", clusterName)
		}

		return nil, errors.Errorf("current cluster is not set. Use the 'switch' command to set the current cluster")
	}

	return cluster, nil
}

func getClusterConfig(cluster repository.Cluster) (*config.ClusterConfig, error) {
	clusterDir := path.Join(*dataDir, "config", cluster.Name)

	kubernetesBinDir := kubernetesBinDir(cluster.KubernetesVersion)
	cniBinDir := path.Join(cacheDir(), "cni", cluster.CNIVersion, "bin")

	return config.New(clusterDir, cluster, kubernetesBinDir, cniBinDir)
}

func readSSHPublicKey(path string) (string, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "cannot load SSH public key from file '%s'", path)
	}

	return string(bytes), nil
}

func nodeDNSName(nodeName string, clusterDomain string) string {
	return fmt.Sprintf("%s.%s", nodeName, clusterDomain)
}

func kubernetesBinDir(kubernetesVersion string) string {
	return path.Join(cacheDir(), "kubernetes", kubernetesVersion, "kubernetes", "server", "bin")
}
