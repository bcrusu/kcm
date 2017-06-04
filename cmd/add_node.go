package cmd

import (
	"strconv"

	"github.com/bcrusu/kcm/cmd/create"
	"github.com/bcrusu/kcm/cmd/start"
	"github.com/bcrusu/kcm/cmd/validate"
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
	"github.com/bcrusu/kcm/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type addNodeCmdState struct {
	ClusterName      string
	IsMaster         bool
	SSHPublicKeyPath string
	Start            bool
	CPUs             uint
	Memory           uint
}

func newAddNodeCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:          "add [NODE_NAME]",
		Short:        "Add a new node to an existing cluster",
		SilenceUsage: true,
	}

	state := &addNodeCmdState{}
	cmd.PersistentFlags().StringVarP(&state.ClusterName, "cluster", "c", "", "Cluster to add the node to. If not specified, the current cluster will be used")
	cmd.PersistentFlags().StringVar(&state.SSHPublicKeyPath, "ssh-public-key", util.GetUserDefaultSSHPublicKeyPath(), "SSH public key to use")
	cmd.PersistentFlags().BoolVarP(&state.Start, "start", "s", false, "Start the node immediately if the cluster is running")
	cmd.PersistentFlags().UintVar(&state.CPUs, "cpu", 1, "Node allocated CPUs")
	cmd.PersistentFlags().UintVar(&state.Memory, "memory", 512, "Node memory (in MiB)")
	cmd.PersistentFlags().BoolVarP(&state.IsMaster, "master", "m", false, "Adds a master node")

	cmd.RunE = state.runE
	return cmd
}

func (s *addNodeCmdState) runE(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return errors.New("invalid command arguments")
	}

	clusterRepository, err := newClusterRepository()
	if err != nil {
		return err
	}

	cluster, err := getWorkingCluster(clusterRepository, s.ClusterName)
	if err != nil {
		return err
	}

	var nodeName string
	if len(args) == 1 {
		nodeName = args[0]
		if _, ok := cluster.Nodes[nodeName]; ok {
			return errors.Errorf("cluster '%s' contains node '%s'", cluster.Name, nodeName)
		}
	} else {
		nodeName = s.nextNodeName(*cluster)
	}

	node := s.createNodeDefinition(nodeName, *cluster)

	// lightweight validation
	if err := node.Validate(); err != nil {
		return err
	}

	sshPublicKey, err := readSSHPublicKey(s.SSHPublicKeyPath)
	if err != nil {
		return err
	}

	connection, err := connectLibvirt()
	if err != nil {
		return err
	}
	defer connection.Close()

	// check for libvirt name conflicts
	if err := validate.LibvirtNodeObjects(connection, node); err != nil {
		return err
	}

	//persist cluster definition before creating any artefacts (libvirt objects/files on disk/etc.)
	cluster.Nodes[node.Name] = node
	if err := clusterRepository.Save(*cluster); err != nil {
		return err
	}

	clusterConfig, err := getClusterConfig(*cluster)
	if err != nil {
		return err
	}

	if err := create.Node(connection, clusterConfig, node, cluster.Network.Name, sshPublicKey); err != nil {
		return err
	}

	return s.startNode(connection, *cluster, node)
}

func (s *addNodeCmdState) createNodeDefinition(name string, cluster repository.Cluster) repository.Node {
	domainName := libvirtDomainName(cluster.Name, name)

	return repository.Node{
		Name:                 name,
		IsMaster:             s.IsMaster,
		Domain:               domainName,
		CPUs:                 s.CPUs,
		MemoryMiB:            s.Memory,
		StoragePool:          cluster.StoragePool,
		BackingStorageVolume: cluster.BackingStorageVolume,
		StorageVolume:        libvirtStorageVolumeName(domainName),
	}
}

func (s *addNodeCmdState) nextNodeName(cluster repository.Cluster) string {
	prefix := NodeNamePrefix
	if s.IsMaster {
		prefix = MasterNodeNamePrefix
	}

	// find the first available name
	for i := 1; ; i++ {
		name := prefix + strconv.FormatInt(int64(i), 10)
		if _, ok := cluster.Nodes[name]; !ok {
			return name
		}
	}
}

func (s *addNodeCmdState) startNode(connection *libvirt.Connection, cluster repository.Cluster, node repository.Node) error {
	if !s.Start {
		return nil
	}

	// start the node only if the cluster is running
	running, err := start.IsClusterRunning(connection, cluster)
	if err != nil {
		return err
	}

	if !running {
		return nil
	}

	return start.Node(connection, node)
}
