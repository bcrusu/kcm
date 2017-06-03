package config

import (
	"crypto/x509"
	"fmt"
	"path"

	"github.com/bcrusu/kcm/config/coreos"
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
	"github.com/bcrusu/kcm/util"
	"github.com/pkg/errors"
)

type ClusterConfig struct {
	clusterDir       string
	cluster          repository.Cluster
	kubernetesBinDir string
	caCertificate    *x509.Certificate
}

type StageNodeResult struct {
	FilesystemMounts []libvirt.FilesystemMount
}

func New(clusterDir string, cluster repository.Cluster, kubernetesCacheDir string) (*ClusterConfig, error) {
	err := util.CreateDirectoryPath(clusterDir)
	if err != nil {
		return nil, err
	}

	caCertificate, err := util.ParseCertificate(cluster.CACertificate)
	if err != nil {
		return nil, err
	}

	return &ClusterConfig{
		clusterDir:       clusterDir,
		cluster:          cluster,
		kubernetesBinDir: path.Join(kubernetesCacheDir, cluster.KubernetesVersion, "kubernetes", "server", "bin"),
		caCertificate:    caCertificate,
	}, nil
}

func (c ClusterConfig) StageNode(name string, sshPublicKey string) (*StageNodeResult, error) {
	node, ok := c.cluster.Nodes[name]
	if !ok {
		return nil, errors.Errorf("cluster '%s' does not contain node '%s'", c.cluster.Name, name)
	}

	nodeDir := path.Join(c.clusterDir, node.Name)
	if err := util.RemoveDirectory(nodeDir); err != nil {
		return nil, errors.Wrapf(err, "failed to remove node config directory '%s'", nodeDir)
	}

	if err := util.CreateDirectoryPath(nodeDir); err != nil {
		return nil, errors.Wrapf(err, "failed to create node config directory '%s'", nodeDir)
	}

	if err := c.stageKubernetes(path.Join(nodeDir, "kubernetes"), node); err != nil {
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

	nodeDir := path.Join(c.clusterDir, node.Name)
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

func (c ClusterConfig) Unstage() error {
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
			GuestPath: "kubernetesConfig",
		},
		libvirt.FilesystemMount{
			HostPath:  c.kubernetesBinDir,
			GuestPath: "kubernetesBin",
		},
	}
}

func (c ClusterConfig) stageCoreOS(outDir string, node repository.Node, sshPublicKey string) error {
	params := coreos.CloudConfigParams{
		Name:          node.Name,
		IsMaster:      node.IsMaster,
		MasterIP:      c.cluster.MasterIP,
		CoreOSChannel: c.cluster.CoreOSChannel,
		SSHPublicKey:  sshPublicKey,
	}

	if err := coreos.WriteCoreOSConfig(outDir, params); err != nil {
		return err
	}

	return nil
}

func (c ClusterConfig) nodeDNSName(nodeName string) string {
	return fmt.Sprintf("%s.%s", nodeName, c.cluster.DNSDomain)
}
