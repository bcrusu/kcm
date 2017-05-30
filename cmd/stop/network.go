package stop

import (
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
	"github.com/golang/glog"
)

func Network(connection *libvirt.Connection, network repository.Network) error {
	name := network.Name

	net, err := connection.GetNetwork(name)
	if err != nil {
		return err
	}

	if net == nil {
		glog.Warningf("cannot find network '%s'", name)
		return nil
	}

	active, err := connection.NetworkIsActive(name)
	if err != nil {
		return err
	}

	if !active {
		return nil
	}

	return connection.DestroyNetwork(name)
}
