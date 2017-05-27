package cmd

import (
	"path"

	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
)

func newClusterRepository() (repository.ClusterRepository, error) {
	repoPath := path.Join(*dataDir, "clusters")
	return repository.New(repoPath)
}

func kubernetesCacheDir() string {
	return path.Join(*dataDir, "kubernetes")
}

func connectLibvirt() (*libvirt.LibvirtConnection, error) {
	return libvirt.NewConnection(*libvirtURI)
}
