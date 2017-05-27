package libvirt

import (
	"github.com/bcrusu/kcm/libvirtxml"
	"github.com/libvirt/libvirt-go"
	"github.com/pkg/errors"
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

	if p == nil {
		// not found
		return nil, nil
	}
	defer p.Free()

	return getStoragePoolXML(p)
}

func (c *LibvirtConnection) GetStorageVolume(pool, volume string) (*libvirtxml.StorageVolume, error) {
	p, err := c.findStoragePool(pool)
	if err != nil {
		return nil, err
	}
	defer p.Free()

	v, err := lookupStorageVolume(p, volume)
	if err != nil {
		return nil, err
	}

	if v == nil {
		// not found
		return nil, nil
	}
	defer v.Free()

	return getStorageVolumeXML(v)
}

func (c *LibvirtConnection) GetDomain(domain string) (*libvirtxml.Domain, error) {
	d, err := lookupDomain(c.connect, domain)
	if err != nil {
		return nil, err
	}

	if d == nil {
		// not found
		return nil, nil
	}
	defer d.Free()

	return getDomainXML(d)
}

func (c *LibvirtConnection) CreateStorageVolume(pool string, name string, backingVolumeName string) error {
	p, err := c.findStoragePool(pool)
	if err != nil {
		return err
	}

	return createStorageVolume(p, name, backingVolumeName)
}

func (c *LibvirtConnection) findStoragePool(name string) (*libvirt.StoragePool, error) {
	pool, err := lookupStoragePool(c.connect, name)
	if err != nil {
		return nil, err
	}

	if pool == nil {
		return nil, errors.Errorf("could not find storage pool '%s'", pool)
	}

	return pool, nil
}
