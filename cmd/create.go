package cmd

import "github.com/spf13/cobra"
import "fmt"

const DefaultKubernetesVersion = "v1.6.4"
const DefaultCoreOSVersion = "1353.7.0"

var (
	kubernetesVersion        = RootCmd.PersistentFlags().String("kube-version", DefaultKubernetesVersion, "Kubernetes version to use")
	kubernetesReleaseChannel = RootCmd.PersistentFlags().String("kube-channel", DefaultKubernetesVersion, "Kubernetes version to use")
	coreOsVersion            = RootCmd.PersistentFlags().String("coreos-version", DefaultCoreOSVersion, "CoreOS version to use")
	coreOsReleaseChannel     = RootCmd.PersistentFlags().String("coreos-channel", "stable", "CoreOS release channel: stable, beta, alpha")
)

var createCmd = &cobra.Command{
	Use:          "create",
	Short:        "Create a new cluster",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := downloadCoreOSImage(DefaultCoreOSVersion, "stable", "/home/bcrusu/Downloads")
		fmt.Println(err)
		return err
	},
}
