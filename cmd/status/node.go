package status

import (
	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
)

type NodeStatus struct {
	Active             bool
	Missing            bool
	InterfaceAddresses []string
}

func Node(connection *libvirt.Connection, node repository.Node) (*NodeStatus, error) {
	domain, err := connection.GetDomain(node.Domain)
	if err != nil {
		return nil, err
	}

	if domain == nil {
		return &NodeStatus{
			Missing: true,
		}, nil
	}

	active, err := connection.DomainIsActive(node.Domain)
	if err != nil {
		return nil, err
	}

	var addresses []string
	if active {
		addresses, err = connection.ListDomainInterfaceAddresses(node.Domain)
		if err != nil {
			return nil, err
		}
	}

	return &NodeStatus{
		Active:             active,
		Missing:            false,
		InterfaceAddresses: addresses,
	}, nil
}
