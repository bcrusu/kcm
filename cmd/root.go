package cmd

import (
	"path"

	"github.com/spf13/cobra"
)

const LibvirtDefaultURI = "qemu:///system"

var (
	configDirAddress = RootCmd.PersistentFlags().String("config-dir", getDefaultConfigDir(), "Directory for the cluster files")
	libvirtUri       = RootCmd.PersistentFlags().String("libvirt-uri", LibvirtDefaultURI, "Libvirt URI")
)

func init() {
	RootCmd.AddCommand(createCmd)
	RootCmd.AddCommand(removeCmd)
	RootCmd.AddCommand(switchCmd)
}

var RootCmd = &cobra.Command{
	Use:          "kcm",
	SilenceUsage: true,
}

func getDefaultConfigDir() string {
	home := getHomeDir()
	return path.Join(home, ".kcm")
}
