package config

import (
	"path"

	"github.com/bcrusu/kcm/config/coreos"
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
	"github.com/bcrusu/kcm/util"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

type ClusterConfig struct {
	clusterDir       string
	cluster          repository.Cluster
	kubernetesBinDir string

	// private/k8s network
	podsNetworkCIDR     string
	servicesNetworkCIDR string
	nonMasqueradeCIDR   string

	// public/libvirt network
	Network util.NetworkInfo
}

type StageNodeResult struct {
	FilesystemMounts []libvirt.FilesystemMount
}

func New(clusterDir string, cluster repository.Cluster, kubernetesCacheDir string) (*ClusterConfig, error) {
	err := util.CreateDirectoryPath(clusterDir)
	if err != nil {
		return nil, err
	}

	network, err := util.ParseNetworkCIDR(cluster.Network.IPv4CIDR)
	if err != nil {
		return nil, err
	}

	return &ClusterConfig{
		clusterDir:          clusterDir,
		cluster:             cluster,
		kubernetesBinDir:    path.Join(kubernetesCacheDir, cluster.KubernetesVersion, "kubernetes", "server", "bin"),
		podsNetworkCIDR:     "10.2.0.0/17",
		servicesNetworkCIDR: "10.2.128.0/17",
		nonMasqueradeCIDR:   "10.2.0.0/16",
		Network:             *network,
	}, nil
}

func (c ClusterConfig) StageNode(name string, sshPublicKey string) (*StageNodeResult, error) {
	node, ok := c.cluster.Nodes[name]
	if !ok {
		return nil, errors.Errorf("cluster '%s' does not contain node '%s'", c.cluster.Name, name)
	}

	nodeDir := c.nodeConfigDir(node.Name)
	if err := prepareDirectory(nodeDir); err != nil {
		return nil, err
	}

	if err := c.stageKubernetesForNode(path.Join(nodeDir, "kubernetes"), node); err != nil {
		return nil, err
	}

	if err := c.stageCoreOS(path.Join(nodeDir, "coreos"), node, sshPublicKey); err != nil {
		return nil, err
	}

	return &StageNodeResult{
		FilesystemMounts: c.getFilesystemMounts(nodeDir),
	}, nil
}

func (c ClusterConfig) UnstageNode(name string) error {
	node, ok := c.cluster.Nodes[name]
	if !ok {
		return errors.Errorf("cluster '%s' does not contain node '%s'", c.cluster.Name, name)
	}

	nodeDir := c.nodeConfigDir(node.Name)
	exists, err := util.DirectoryExists(nodeDir)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	if err := util.RemoveDirectory(nodeDir); err != nil {
		return errors.Wrapf(err, "failed to remove node config directory '%s'", nodeDir)
	}

	return nil
}

func (c ClusterConfig) StageCluster() error {
	if err := prepareDirectory(c.clusterDir); err != nil {
		return err
	}

	if err := c.stageKubernetesForCluster(c.clusterDir); err != nil {
		return err
	}

	return nil
}

func (c ClusterConfig) UnstageCluster() error {
	exists, err := util.DirectoryExists(c.clusterDir)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	if err := util.RemoveDirectory(c.clusterDir); err != nil {
		return errors.Wrapf(err, "failed to remove cluster config directory '%s'", c.clusterDir)
	}

	return nil
}

func (c ClusterConfig) getFilesystemMounts(nodeDir string) []libvirt.FilesystemMount {
	return []libvirt.FilesystemMount{
		libvirt.FilesystemMount{
			HostPath:  path.Join(nodeDir, "coreos"),
			GuestPath: "config-2",
		},
		libvirt.FilesystemMount{
			HostPath:  path.Join(nodeDir, "kubernetes"),
			GuestPath: "k8sConfig",
		},
		libvirt.FilesystemMount{
			HostPath:  c.kubernetesBinDir,
			GuestPath: "k8sBin",
		},
		libvirt.FilesystemMount{
			HostPath:  path.Join(c.clusterDir, "manifests"),
			GuestPath: "k8sConfigManifests",
		},
		libvirt.FilesystemMount{
			HostPath:  path.Join(c.clusterDir, "kubeconfig"),
			GuestPath: "k8sConfigKubeconfig",
		},
	}
}

func (c ClusterConfig) stageCoreOS(outDir string, node repository.Node, sshPublicKey string) error {
	params := coreos.CloudConfigParams{
		Hostname:          node.Name,
		IsMaster:          node.IsMaster,
		SSHPublicKey:      sshPublicKey,
		NonMasqueradeCIDR: c.nonMasqueradeCIDR,
		Network:           c.Network,
		ClusterDomain:     c.cluster.DNSDomain,
	}

	if err := coreos.WriteCoreOSConfig(outDir, params); err != nil {
		return err
	}

	return nil
}

func (c ClusterConfig) nodeConfigDir(nodeName string) string {
	return path.Join(c.clusterDir, "nodes", nodeName)
}

func prepareDirectory(dir string) error {
	exists, err := util.DirectoryExists(dir)
	if err != nil {
		return err
	}

	if exists {
		glog.Warningf("config directory '%s' exists - will be deleted", dir)

		if err := util.RemoveDirectory(dir); err != nil {
			return errors.Wrapf(err, "failed to remove config directory '%s'", dir)
		}
	}

	if err := util.CreateDirectoryPath(dir); err != nil {
		return errors.Wrapf(err, "failed to create config directory '%s'", dir)
	}

	return nil
}
