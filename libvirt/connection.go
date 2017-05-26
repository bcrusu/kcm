package libvirt

import (
	"github.com/bcrusu/kcm/libvirtxml"
	"github.com/libvirt/libvirt-go"
)

type LibvirtConnection struct {
	uri     string
	connect *libvirt.Connect
}

func NewConnection(uri string) (*LibvirtConnection, error) {
	//TODO: allow only local connections

	connect, err := libvirt.NewConnect(uri)
	if err != nil {
		return nil, err
	}

	return &LibvirtConnection{
		uri:     uri,
		connect: connect,
	}, nil
}

func (c *LibvirtConnection) Close() {
	c.connect.Close()
	c.connect = nil
}

func (c *LibvirtConnection) GetStoragePool(pool string) (*libvirtxml.StoragePool, error) {
	p, err := lookupStoragePool(c.connect, pool)
	if err != nil {
		return nil, err
	}
	defer p.Free()

	return getStoragePoolXML(p)
}
