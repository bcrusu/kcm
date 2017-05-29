package remove

import (
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
)

func RemoveNetwork(connection *libvirt.Connection, network repository.Network) error {
	name := network.Name

	net, err := connection.GetNetwork(name)
	if err != nil {
		return err
	}

	if net == nil {
		// network does not exist
		return nil
	}

	active, err := connection.NetworkIsActive(name)
	if err != nil {
		return err
	}

	if active {
		if err := connection.DestroyNetwork(name); err != nil {
			return err
		}
	}

	return connection.UndefineNetwork(name)
}
