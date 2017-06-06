package cmd

import (
	"fmt"
	"strconv"

	"github.com/bcrusu/kcm/cmd/create"
	"github.com/bcrusu/kcm/cmd/download"
	"github.com/bcrusu/kcm/cmd/start"
	"github.com/bcrusu/kcm/cmd/validate"
	"github.com/bcrusu/kcm/repository"
	"github.com/bcrusu/kcm/util"
	"github.com/golang/glog"
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
	NodesCount         uint
	KubernetesNetwork  string
	SSHPublicKeyPath   string
	Start              bool
	IPv4CIDR           string
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
	cmd.PersistentFlags().UintVar(&state.NodesCount, "node-count", 1, "Initial number of nodes in the cluster") //TODO: set value to 3
	cmd.PersistentFlags().StringVar(&state.KubernetesNetwork, "kubernetes-network", "flannel", "Networking mode to use. Only flannel is suppoted at the moment")
	cmd.PersistentFlags().StringVar(&state.SSHPublicKeyPath, "ssh-public-key", util.GetUserDefaultSSHPublicKeyPath(), "SSH public key to use")
	cmd.PersistentFlags().BoolVarP(&state.Start, "start", "s", false, "Start the cluster immediately")
	cmd.PersistentFlags().StringVar(&state.IPv4CIDR, "ipv4-cidr", "10.1.0.0/16", "Libvirt network IPv4 CIDR. Network 10.2.0.0/16 is reserved for pods/services network")
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

	cluster, err := s.createClusterDefinition(args[0])
	if err != nil {
		return err
	}

	// lightweight validation - valid names, empty strings, no nodes defined, etc.
	if err := cluster.Validate(); err != nil {
		return err
	}

	clusterRepository, err := newClusterRepository()
	if err != nil {
		return err
	}

	exists, err := clusterRepository.Exists(cluster.Name)
	if err != nil {
		return err
	}

	if exists {
		return errors.Errorf("failed to create cluster. Cluster '%s' already exists", cluster.Name)
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

	// check for name conflicts/missing libvirt objects
	if err := validate.LibvirtClusterObjects(connection, cluster); err != nil {
		return err
	}

	if err := download.DownloadPrerequisites(connection, cluster, kubernetesCacheDir()); err != nil {
		return err
	}

	//persist cluster definition before creating any artefacts (libvirt/files on disk/etc.)
	if err := clusterRepository.Save(cluster); err != nil {
		return err
	}

	clusterConfig, err := getClusterConfig(cluster)
	if err != nil {
		return err
	}

	if err := create.Cluster(connection, clusterConfig, cluster, sshPublicKey); err != nil {
		return err
	}

	s.setActiveCluster(clusterRepository, cluster.Name)

	if s.Start {
		return start.Cluster(connection, cluster)
	}

	return nil
}

func (s *createCmdState) createClusterDefinition(clusterName string) (repository.Cluster, error) {
	backingStorageVolume := coreOSStorageVolumeName(s.CoreOSVersion)

	caCertificateBytes, caKeyBytes, err := util.CreateCACertificate(clusterName + "-ca")
	if err != nil {
		return repository.Cluster{}, err
	}

	caCertificate, err := util.ParseCertificate(caCertificateBytes)
	if err != nil {
		return repository.Cluster{}, err
	}

	cluster := repository.Cluster{
		Name:                 clusterName,
		KubernetesVersion:    s.KubernetesVersion,
		CoreOSChannel:        s.CoreOSChannel,
		CoreOSVersion:        s.CoreOSVersion,
		StoragePool:          s.LibvirtStoragePool,
		BackingStorageVolume: backingStorageVolume,
		Network: repository.Network{
			Name:     libvirtNetworkName(clusterName),
			IPv4CIDR: s.IPv4CIDR,
		},
		Nodes:         make(map[string]repository.Node),
		CACertificate: caCertificateBytes,
		CAPrivateKey:  caKeyBytes,
		DNSDomain:     clusterName + ".local",
	}

	addNode := func(name string, isMaster bool) error {
		domainName := libvirtDomainName(clusterName, name)
		dnsName := nodeDNSName(name, cluster.DNSDomain)

		certificate, key, err := util.CreateCertificate(dnsName, caCertificate, dnsName)
		if err != nil {
			return err
		}

		node := repository.Node{
			Name:                 name,
			IsMaster:             isMaster,
			Domain:               domainName,
			StoragePool:          s.LibvirtStoragePool,
			BackingStorageVolume: backingStorageVolume,
			StorageVolume:        libvirtStorageVolumeName(domainName),
			Certificate:          certificate,
			PrivateKey:           key,
			DNSName:              dnsName,
		}

		if isMaster {
			node.CPUs = s.MasterCPUs
			node.MemoryMiB = s.MasterMemory
		} else {
			node.CPUs = s.NondeCPUs
			node.MemoryMiB = s.NodeMemory
		}

		cluster.Nodes[name] = node
		return nil
	}

	masterName := MasterNodeNamePrefix
	if err := addNode(masterName, true); err != nil {
		return repository.Cluster{}, err
	}

	for i := uint(1); i <= s.NodesCount; i++ {
		name := NodeNamePrefix + strconv.FormatUint(uint64(i), 10)
		if err := addNode(name, false); err != nil {
			return repository.Cluster{}, err
		}
	}

	cluster.ServerURL = fmt.Sprintf("https://%s:6443", cluster.Nodes[masterName].DNSName)

	return cluster, nil
}

func (s *createCmdState) setActiveCluster(clusterRepository repository.ClusterRepository, name string) {
	currentActiveCluster, err := clusterRepository.Current()
	if err != nil {
		glog.Error(err)
		return
	}

	if currentActiveCluster != nil {
		// only set the active cluster if none is currently set
		return
	}

	if err := clusterRepository.SetCurrent(name); err != nil {
		glog.Errorf("failed to set current cluster. Error: %v", err)
	}
}
