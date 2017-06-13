package status

import (
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
)

type NetworkStatus struct {
	Active  bool
	Missing bool
}

func Network(connection *libvirt.Connection, network repository.Network) (*NetworkStatus, error) {
	net, err := connection.GetNetwork(network.Name)
	if err != nil {
		return nil, err
	}

	if net == nil {
		return &NetworkStatus{
			Missing: true,
		}, nil
	}

	active, err := connection.NetworkIsActive(network.Name)
	if err != nil {
		return nil, err
	}

	return &NetworkStatus{
		Active: active,
	}, nil
}
