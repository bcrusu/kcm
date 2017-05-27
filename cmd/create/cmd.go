package create

import (
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const DefaultKubernetesVersion = "v1.6.4"
const DefaultCoreOSVersion = "1353.7.0"

var (
	kubernetesVersion    = Cmd.PersistentFlags().String("kubernetes-version", DefaultKubernetesVersion, "Kubernetes version to use")
	coreOsVersion        = Cmd.PersistentFlags().String("coreos-version", DefaultCoreOSVersion, "CoreOS version to use")
	coreOsReleaseChannel = Cmd.PersistentFlags().String("coreos-channel", "stable", "CoreOS release channel: stable, beta, alpha")
	libvirtStoragePool   = Cmd.PersistentFlags().String("libvirt-pool", "default", "Libvirt storage pool")
	clusterName          = Cmd.PersistentFlags().StringP("name", "n", "kube", "Cluster name")
	masterCount          = Cmd.PersistentFlags().IntP("master-count", "mc", 1, "Initial number of masters in the cluster")
	nodesCount           = Cmd.PersistentFlags().IntP("node-count", "nc", 3, "Initial number of nodes in the cluster")
	networking           = Cmd.PersistentFlags().String("networking", "flannel", "Networking mode to use. Only flannel is suppoted at the moment")
	sshPublicKey         = Cmd.PersistentFlags().String("ssh-public-key", util.GetUserDefaultSSHPublicKeyPath(), "SSH public key to use")
	start                = Cmd.PersistentFlags().Bool("start", false, "Start the cluster immediately")
)

func init() {
	Cmd.RunE = runE
	Cmd.MarkPersistentFlagRequired("name")
}

var Cmd = &cobra.Command{
	Use:          "create",
	Short:        "Create a new cluster",
	SilenceUsage: true,
}

func runE(cmd *cobra.Command, args []string) error {
	libvirtURI := cmd.Flag("libvirt-uri").Value.String()

	connection, err := libvirt.NewConnection(libvirtURI)
	if err != nil {
		return err
	}
	defer connection.Close()

	if err := ensureBackingStorageVolume(connection); err != nil {
		return err
	}

	return nil
}

func ensureBackingStorageVolume(connection *libvirt.LibvirtConnection) error {
	poolPath, err := getStoragePoolTargetPath(connection, *libvirtStoragePool)
	if err != nil {
		return err
	}

	volumeName := coreOSVolumeName(*coreOsVersion)

	volume, err := connection.GetStorageVolume(*libvirtStoragePool, volumeName)
	if err != nil {
		return err
	}

	if volume == nil {
		if err := downloadCoreOSImage(*coreOsVersion, *coreOsReleaseChannel, poolPath); err != nil {
			return err
		}
	}

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
