package start

import (
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
	"github.com/pkg/errors"
)

func Network(connection *libvirt.Connection, network repository.Network) error {
	name := network.Name

	net, err := connection.GetNetwork(name)
	if err != nil {
		return err
	}

	if net == nil {
		return errors.Errorf("cannot find network '%s'", name)
	}

	active, err := connection.NetworkIsActive(name)
	if err != nil {
		return err
	}

	if active {
		// network is already running
		return nil
	}

	return connection.CreateNetwork(name)
}
