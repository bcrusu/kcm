package status

import (
	"net"

	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
)

type NodeStatus struct {
	Active    bool
	Missing   bool
	Addresses []string
	Node      repository.Node
}

func Node(connection *libvirt.Connection, node repository.Node) (*NodeStatus, error) {
	domain, err := connection.GetDomain(node.Domain)
	if err != nil {
		return nil, err
	}

	if domain == nil {
		return &NodeStatus{
			Missing: true,
			Node:    node,
		}, nil
	}

	active, err := connection.DomainIsActive(node.Domain)
	if err != nil {
		return nil, err
	}

	var addresses []string
	if active {
		addresses, err = net.LookupHost(node.DNSName)
		if err != nil {
			if _, ok := err.(*net.DNSError); ok {
				addresses = []string{"DNS lookup failed"}
			} else {
				return nil, err
			}
		}
	}

	return &NodeStatus{
		Active:    active,
		Missing:   false,
		Addresses: addresses,
		Node:      node,
	}, nil
}
