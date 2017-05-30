package stop

import (
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
	"github.com/golang/glog"
)

func Node(connection *libvirt.Connection, node repository.Node, force bool) error {
	name := node.Domain

	domain, err := connection.GetDomain(name)
	if err != nil {
		return err
	}

	if domain == nil {
		glog.Warningf("cannot find domain '%s'", name)
		// ignore missing domains
		return nil
	}

	active, err := connection.DomainIsActive(name)
	if err != nil {
		return err
	}

	if !active {
		return nil
	}

	if force {
		return connection.DestroyDomain(name)
	}

	return connection.ShutdownDomain(name)
}
