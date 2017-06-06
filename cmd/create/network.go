package create

import (
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
)

func Network(connection *libvirt.Connection, network repository.Network, clusterDomain string) error {
	params := libvirt.DefineNetworkParams{
		Name:     network.Name,
		IPv4CIDR: network.IPv4CIDR,
		Domain:   clusterDomain,
	}

	return connection.DefineNATNetwork(params)
}
