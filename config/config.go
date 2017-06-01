package config

import (
	"path"

	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
	"github.com/bcrusu/kcm/util"
	"github.com/pkg/errors"
)

type ClusterConfig struct {
	configDir    string
	cluster      repository.Cluster
	sshPublicKey string
}

type StageNodeResult struct {
	FilesystemMounts []libvirt.FilesystemMount
}

func New(configDir string, cluster repository.Cluster, sshPublicKey string) (*ClusterConfig, error) {
	err := util.CreateDirectoryPath(configDir)
	if err != nil {
		return nil, err
	}

	return &ClusterConfig{
		configDir:    configDir,
		cluster:      cluster,
		sshPublicKey: sshPublicKey,
	}, nil
}

func (c ClusterConfig) StageNode(name string) (*StageNodeResult, error) {
	node, ok := c.cluster.Nodes[name]
	if !ok {
		return nil, errors.Errorf("cluster '%s' does not contain node '%s'", c.cluster.Name, name)
	}

	nodeConfigDir := path.Join(c.configDir, node.Name)
	if err := emptyNodeConfigDir(nodeConfigDir); err != nil {
		return nil, err
	}

	// _ALL_
	//   kubernetes
	//   	bin

	// NODE_NAME
	//   config-2/openstack/latest/user_data
	//   kubernetes
	//	   /certs

	// masters
	//   manifests (static pods) - master nodes only
	//   addons - master nodes only
	// nodes

	{
		params := coreOSTemplateParams{
			Name:          node.Name,
			IsMaster:      node.IsMaster,
			MasterIP:      c.cluster.MasterIP,
			CoreOSChannel: c.cluster.CoreOSChannel,
			SSHPublicKey:  c.sshPublicKey,
		}

		if err := writeCoreOSConfig(nodeConfigDir, params); err != nil {
			return nil, err
		}
	}

	return &StageNodeResult{
		FilesystemMounts: getFilesystemMounts(nodeConfigDir),
	}, nil
}

func emptyNodeConfigDir(path string) error {
	if err := util.RemoveDirectory(path); err != nil {
		return errors.Wrapf(err, "failed to remove node config directory '%s'", path)
	}

	if err := util.CreateDirectoryPath(path); err != nil {
		return errors.Wrapf(err, "failed to create node config directory '%s'", path)
	}

	return nil
}

func getFilesystemMounts(nodeConfigDir string) []libvirt.FilesystemMount {
	return []libvirt.FilesystemMount{
		libvirt.FilesystemMount{
			HostPath:  path.Join(nodeConfigDir, "config-2"),
			GuestPath: "config-2",
		},
		libvirt.FilesystemMount{
			HostPath:  path.Join(nodeConfigDir, "kubernetes"),
			GuestPath: "kubernetes",
		},
	}
}
