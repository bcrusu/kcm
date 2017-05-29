package cmd

import (
	"path"

	"github.com/bcrusu/kcm/util"
	"github.com/spf13/cobra"
)

const LibvirtDefaultURI = "qemu:///system"

var (
	dataDir    = RootCmd.PersistentFlags().String("data-dir", getDefaultDataDir(), "kcm data directory")
	libvirtURI = RootCmd.PersistentFlags().String("libvirt-uri", LibvirtDefaultURI, "Libvirt URI")
)

func init() {
	RootCmd.AddCommand(createCmd)
	RootCmd.AddCommand(newRemoveCmd())
	RootCmd.AddCommand(switchCmd)
}

var RootCmd = &cobra.Command{
	Use:          "kcm",
	SilenceUsage: true,
}

func getDefaultDataDir() string {
	home := util.GetUserHomeDir()
	return path.Join(home, ".kcm")
}
