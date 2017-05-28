package cmd

import (
	"github.com/bcrusu/kcm/cmd/create"
	"github.com/bcrusu/kcm/repository"
	"github.com/bcrusu/kcm/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const DefaultKubernetesVersion = "1.6.4"
const DefaultCoreOSVersion = "1353.7.0"
const DefaultCoreOsChannel = "stable"

var (
	kubernetesVersion  = createCmd.PersistentFlags().String("kubernetes-version", DefaultKubernetesVersion, "Kubernetes version to use")
	coreOSVersion      = createCmd.PersistentFlags().String("coreos-version", DefaultCoreOSVersion, "CoreOS version to use")
	coreOSChannel      = createCmd.PersistentFlags().String("coreos-channel", DefaultCoreOsChannel, "CoreOS release channel: stable, beta, alpha")
	libvirtStoragePool = createCmd.PersistentFlags().String("libvirt-pool", "default", "Libvirt storage pool")
	clusterName        = createCmd.PersistentFlags().String("name", "kube", "Cluster name")
	masterCount        = createCmd.PersistentFlags().Uint("master-count", 1, "Initial number of masters in the cluster")
	nodesCount         = createCmd.PersistentFlags().Uint("node-count", 3, "Initial number of nodes in the cluster")
	kubernetesNetwork  = createCmd.PersistentFlags().String("kubernetes-network", "flannel", "Networking mode to use. Only flannel is suppoted at the moment")
	sshPublicKey       = createCmd.PersistentFlags().String("ssh-public-key", util.GetUserDefaultSSHPublicKeyPath(), "SSH public key to use")
	start              = createCmd.PersistentFlags().Bool("start", false, "Start the cluster immediately")
	ipv4CIDR           = createCmd.PersistentFlags().String("ipv4-cidr", "10.11.0.1/24", "Libvirt network IPv4 CIDR")
	ipv6CIDR           = createCmd.PersistentFlags().String("ipv6-cidr", "", "Libvirt network IPv6 CIDR")
	masterCPUs         = createCmd.PersistentFlags().Uint("master-cpu", 2, "Master node allocated CPUs")
	masterMemory       = createCmd.PersistentFlags().Uint("master-memory", 1024, "Master node memory (in MiB)")
	nondeCPUs          = createCmd.PersistentFlags().Uint("node-cpu", 1, "Node allocated CPUs")
	nodeMemory         = createCmd.PersistentFlags().Uint("node-memory", 512, "Node memory (in MiB)")
)

func init() {
	createCmd.RunE = runE
	createCmd.MarkPersistentFlagRequired("name")
}

var createCmd = &cobra.Command{
	Use:          "create",
	Short:        "Create a new cluster",
	SilenceUsage: true,
}

func runE(cmd *cobra.Command, args []string) error {
	cluster := createClusterDefinition()
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

	// if err := create.CreateLibvirtObjects(connection, cluster); err != nil {
	// 	return err
	// }

	return nil
}

func createClusterDefinition() repository.Cluster {
	cluster := repository.Cluster{
		Name:                 *clusterName,
		KubernetesVersion:    *kubernetesVersion,
		CoreOSChannel:        *coreOSChannel,
		CoreOSVersion:        *coreOSVersion,
		StoragePool:          *libvirtStoragePool,
		BackingStorageVolume: coreOSStorageVolumeName(*coreOSVersion),
		Network: repository.Network{
			Name:     libvirtNetworkName(*clusterName),
			IPv4CIDR: *ipv4CIDR,
			IPv6CIDR: *ipv6CIDR,
		},
	}

	for i := uint(1); i <= *masterCount; i++ {
		domainName := libvirtDomainName(*clusterName, true, i)

		cluster.Masters = append(cluster.Masters, repository.Node{
			Domain:        domainName,
			StoragePool:   *libvirtStoragePool,
			StorageVolume: libvirtStorageVolumeName(domainName),
			CPUs:          *masterCPUs,
			MemoryMiB:     *masterMemory,
		})
	}

	for i := uint(1); i <= *nodesCount; i++ {
		domainName := libvirtDomainName(*clusterName, false, i)

		cluster.Nodes = append(cluster.Nodes, repository.Node{
			Domain:        domainName,
			StoragePool:   *libvirtStoragePool,
			StorageVolume: libvirtStorageVolumeName(domainName),
			CPUs:          *nondeCPUs,
			MemoryMiB:     *nodeMemory,
		})
	}

	return cluster
}
