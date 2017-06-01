package create

import (
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
)

func Network(connection *libvirt.Connection, network repository.Network) error {
	params := libvirt.DefineNetworkParams{
		Name:     network.Name,
		IPv4CIDR: network.IPv4CIDR,
	}

	return connection.DefineNATNetwork(params)
}
