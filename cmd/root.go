package cmd

import (
	"path"

	"github.com/bcrusu/kcm/cmd/create"
	"github.com/bcrusu/kcm/util"
	"github.com/spf13/cobra"
)

const LibvirtDefaultURI = "qemu:///system"

var (
	repositoryAddress = RootCmd.PersistentFlags().String("repository", getDefaultRepositoryPath(), "Cluster repository path. Cluster definitions will be placed here.")
	libvirtURI        = RootCmd.PersistentFlags().String("libvirt-uri", LibvirtDefaultURI, "Libvirt URI")
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

func getDefaultRepositoryPath() string {
	home := util.GetUserHomeDir()
	return path.Join(home, ".kcm")
}
