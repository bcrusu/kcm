package cmd

import (
	"strconv"

	"github.com/bcrusu/kcm/cmd/create"
	"github.com/bcrusu/kcm/repository"
	"github.com/bcrusu/kcm/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const DefaultKubernetesVersion = "1.6.4"
const DefaultCoreOSVersion = "1353.7.0"
const DefaultCoreOsChannel = "stable"

type createCmdState struct {
	KubernetesVersion  string
	CoreOSVersion      string
	CoreOSChannel      string
	LibvirtStoragePool string
	ClusterName        string
	MasterCount        uint
	NodesCount         uint
	KubernetesNetwork  string
	SSHPublicKey       string
	Start              bool
	IPv4CIDR           string
	IPv6CIDR           string
	MasterCPUs         uint
	MasterMemory       uint
	NondeCPUs          uint
	NodeMemory         uint
}

func newCreateCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:          "create CLUSTER_NAME",
		Short:        "Create a new cluster",
		SilenceUsage: true,
	}

	state := &createCmdState{}

	cmd.PersistentFlags().StringVar(&state.KubernetesVersion, "kubernetes-version", DefaultKubernetesVersion, "Kubernetes version to use")
	cmd.PersistentFlags().StringVar(&state.CoreOSVersion, "coreos-version", DefaultCoreOSVersion, "CoreOS version to use")
	cmd.PersistentFlags().StringVar(&state.CoreOSChannel, "coreos-channel", DefaultCoreOsChannel, "CoreOS release channel: stable, beta, alpha")
	cmd.PersistentFlags().StringVar(&state.LibvirtStoragePool, "libvirt-pool", "default", "Libvirt storage pool")
	cmd.PersistentFlags().UintVar(&state.MasterCount, "master-count", 1, "Initial number of masters in the cluster")
	cmd.PersistentFlags().UintVar(&state.NodesCount, "node-count", 3, "Initial number of nodes in the cluster")
	cmd.PersistentFlags().StringVar(&state.KubernetesNetwork, "kubernetes-network", "flannel", "Networking mode to use. Only flannel is suppoted at the moment")
	cmd.PersistentFlags().StringVar(&state.SSHPublicKey, "ssh-public-key", util.GetUserDefaultSSHPublicKeyPath(), "SSH public key to use")
	cmd.PersistentFlags().BoolVar(&state.Start, "start", false, "Start the cluster immediately")
	cmd.PersistentFlags().StringVar(&state.IPv4CIDR, "ipv4-cidr", "10.11.0.1/24", "Libvirt network IPv4 CIDR")
	cmd.PersistentFlags().StringVar(&state.IPv6CIDR, "ipv6-cidr", "", "Libvirt network IPv6 CIDR")
	cmd.PersistentFlags().UintVar(&state.MasterCPUs, "master-cpu", 1, "Master node allocated CPUs")
	cmd.PersistentFlags().UintVar(&state.MasterMemory, "master-memory", 512, "Master node memory (in MiB)")
	cmd.PersistentFlags().UintVar(&state.NondeCPUs, "node-cpu", 1, "Node allocated CPUs")
	cmd.PersistentFlags().UintVar(&state.NodeMemory, "node-memory", 512, "Node memory (in MiB)")

	cmd.RunE = state.runE
	return cmd
}

func (s *createCmdState) runE(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("invalid command arguments")
	}

	s.ClusterName = args[0]

	cluster := s.createClusterDefinition()
	// lightweight validation - empty strings, no nodes defined, etc.
	if err := cluster.Validate(); err != nil {
		return err
	}

	clusterRepository, err := newClusterRepository()
	if err != nil {
		return err
	}

	if clusterRepository.Exists(cluster.Name) {
		return errors.Errorf("failed to create cluster. Cluster '%s' already exists", cluster.Name)
	}

	connection, err := connectLibvirt()
	if err != nil {
		return err
	}
	defer connection.Close()

	// check for name conflicts/missing libvirt objects
	if err := create.ValidateLibvirtObjects(connection, cluster); err != nil {
		return err
	}

	if err := create.DownloadPrerequisites(connection, cluster, kubernetesCacheDir()); err != nil {
		return err
	}

	// persist cluster definition before creating any cluster-spepcific artefacts (libvirt/files on disk)
	if err := clusterRepository.Save(cluster); err != nil {
		return err
	}

	//TODO: create filesystems to be mounted

	if err := create.CreateLibvirtObjects(connection, cluster); err != nil {
		return err
	}

	// TODO: make the cluster active

	return nil
}

func (s *createCmdState) createClusterDefinition() repository.Cluster {
	cluster := repository.Cluster{
		Name:                 s.ClusterName,
		KubernetesVersion:    s.KubernetesVersion,
		CoreOSChannel:        s.CoreOSChannel,
		CoreOSVersion:        s.CoreOSVersion,
		StoragePool:          s.LibvirtStoragePool,
		BackingStorageVolume: coreOSStorageVolumeName(s.CoreOSVersion),
		Network: repository.Network{
			Name:     libvirtNetworkName(s.ClusterName),
			IPv4CIDR: s.IPv4CIDR,
			IPv6CIDR: s.IPv6CIDR,
		},
	}

	for i := uint(1); i <= s.MasterCount; i++ {
		domainName := libvirtDomainName(s.ClusterName, true, strconv.FormatUint(uint64(i), 10))

		cluster.Masters = append(cluster.Masters, repository.Node{
			Domain:        domainName,
			StoragePool:   s.LibvirtStoragePool,
			StorageVolume: libvirtStorageVolumeName(domainName),
			CPUs:          s.MasterCPUs,
			MemoryMiB:     s.MasterMemory,
		})
	}

	for i := uint(1); i <= s.NodesCount; i++ {
		domainName := libvirtDomainName(s.ClusterName, false, strconv.FormatUint(uint64(i), 10))

		cluster.Nodes = append(cluster.Nodes, repository.Node{
			Domain:        domainName,
			StoragePool:   s.LibvirtStoragePool,
			StorageVolume: libvirtStorageVolumeName(domainName),
			CPUs:          s.NondeCPUs,
			MemoryMiB:     s.NodeMemory,
		})
	}

	return cluster
}
