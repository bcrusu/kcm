package cmd

import (
	"github.com/bcrusu/kcm/cmd/create"
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const DefaultKubernetesVersion = "1.6.4"
const DefaultCoreOSVersion = "1353.7.0"
const DefaultCoreOsReleaseChannel = "stable"

var (
	kubernetesVersion    = createCmd.PersistentFlags().String("kubernetes-version", DefaultKubernetesVersion, "Kubernetes version to use")
	coreOsVersion        = createCmd.PersistentFlags().String("coreos-version", DefaultCoreOSVersion, "CoreOS version to use")
	coreOsReleaseChannel = createCmd.PersistentFlags().String("coreos-channel", DefaultCoreOsReleaseChannel, "CoreOS release channel: stable, beta, alpha")
	libvirtStoragePool   = createCmd.PersistentFlags().String("libvirt-pool", "default", "Libvirt storage pool")
	clusterName          = createCmd.PersistentFlags().String("name", "kube", "Cluster name")
	masterCount          = createCmd.PersistentFlags().Int("master-count", 1, "Initial number of masters in the cluster")
	nodesCount           = createCmd.PersistentFlags().Int("node-count", 3, "Initial number of nodes in the cluster")
	networking           = createCmd.PersistentFlags().String("networking", "flannel", "Networking mode to use. Only flannel is suppoted at the moment")
	sshPublicKey         = createCmd.PersistentFlags().String("ssh-public-key", util.GetUserDefaultSSHPublicKeyPath(), "SSH public key to use")
	start                = createCmd.PersistentFlags().Bool("start", false, "Start the cluster immediately")
)

func init() {
	createCmd.RunE = runE
}

var createCmd = &cobra.Command{
	Use:          "create",
	Short:        "Create a new cluster",
	SilenceUsage: true,
}

func runE(cmd *cobra.Command, args []string) error {
	repository, err := newClusterRepository()
	if err != nil {
		return err
	}

	if repository.Exists(*clusterName) {
		return errors.Errorf("failed to create cluster. Cluster '%s' already exists", *clusterName)
	}

	if err := create.DownloadKubernetes(*kubernetesVersion, kubernetesCacheDir()); err != nil {
		return err
	}

	connection, err := connectLibvirt()
	if err != nil {
		return err
	}
	defer connection.Close()

	if err := downloadBackingStorageImage(connection); err != nil {
		return err
	}

	return nil
}

func downloadBackingStorageImage(connection *libvirt.LibvirtConnection) error {
	poolPath, err := getStoragePoolTargetPath(connection, *libvirtStoragePool)
	if err != nil {
		return err
	}

	volumeName := create.CoreOSVolumeName(*coreOsVersion)

	volume, err := connection.GetStorageVolume(*libvirtStoragePool, volumeName)
	if err != nil {
		return err
	}

	if volume == nil {
		if err := create.DownloadCoreOS(*coreOsVersion, *coreOsReleaseChannel, poolPath); err != nil {
			return err
		}
	}

	// volume exists - will not download
	return nil
}

func getStoragePoolTargetPath(connection *libvirt.LibvirtConnection, pool string) (string, error) {
	p, err := connection.GetStoragePool(pool)
	if err != nil {
		return "", err
	}

	if p == nil {
		return "", errors.Errorf("could not find storage pool '%s'", pool)
	}

	poolType := p.Type()
	if poolType != "dir" && poolType != "fs" {
		return "", errors.Errorf("unsupported storage pool type '%s'", poolType)
	}

	poolPath := p.Target().Path()
	if poolPath == "" {
		return "", errors.Errorf("invalid storage pool '%s'. Target path is empty", poolType)
	}

	return poolPath, nil
}
