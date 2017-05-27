package cmd

import (
	"path"

	"github.com/bcrusu/kcm/cmd/create"
	"github.com/bcrusu/kcm/util"
	"github.com/spf13/cobra"
)

const LibvirtDefaultURI = "qemu:///system"

var (
	configDirAddress = RootCmd.PersistentFlags().String("config-dir", getDefaultConfigDir(), "Directory for the cluster files")
	libvirtURI       = RootCmd.PersistentFlags().String("libvirt-uri", LibvirtDefaultURI, "Libvirt URI")
)

func init() {
	RootCmd.AddCommand(create.Cmd)
	RootCmd.AddCommand(removeCmd)
	RootCmd.AddCommand(switchCmd)
}

var RootCmd = &cobra.Command{
	Use:          "kcm",
	SilenceUsage: true,
}

func getDefaultConfigDir() string {
	home := util.GetUserHomeDir()
	return path.Join(home, ".kcm")
}
