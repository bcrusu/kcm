package start

import (
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
	"github.com/pkg/errors"
)

func Network(connection *libvirt.Connection, network repository.Network) error {
	name := network.Name

	active, err := isNetworkActive(connection, name)
	if err != nil {
		return err
	}

	if active {
		// network is already running
		return nil
	}

	return connection.CreateNetwork(name)
}

func isNetworkActive(connection *libvirt.Connection, name string) (bool, error) {
	net, err := connection.GetNetwork(name)
	if err != nil {
		return false, err
	}

	if net == nil {
		return false, errors.Errorf("cannot find network '%s'", name)
	}

	active, err := connection.NetworkIsActive(name)
	if err != nil {
		return false, err
	}

	return active, nil
}
