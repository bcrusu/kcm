package libvirt

import (
	"strings"

	"github.com/bcrusu/kcm/libvirtxml"
	"github.com/libvirt/libvirt-go"
	"github.com/pkg/errors"
)

const MetadataXMLNamespace = "https://github.com/bcrusu/kcm"

type Connection struct {
	uri          string
	connect      *libvirt.Connect
	capabilities *libvirtxml.Capabilities
}

func NewConnection(uri string) (*Connection, error) {
	//TODO: allow only local connections

	connect, err := libvirt.NewConnect(uri)
	if err != nil {
		return nil, err
	}

	return &Connection{
		uri:     uri,
		connect: connect,
	}, nil
}

func (c *Connection) Close() {
	c.connect.Close()
	c.connect = nil
}

func (c *Connection) GetCapabilities() (*libvirtxml.Capabilities, error) {
	if c.capabilities == nil {
		capabilities, err := getCapabilitiesXML(c.connect)
		if err != nil {
			return nil, err
		}

		c.capabilities = capabilities
	}

	return c.capabilities, nil
}

func (c *Connection) GetStoragePool(pool string) (*libvirtxml.StoragePool, error) {
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

func (c *Connection) GetStorageVolume(pool, volume string) (*libvirtxml.StorageVolume, error) {
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

func (c *Connection) GetDomain(domain string) (*libvirtxml.Domain, error) {
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

func (c *Connection) ListAllDomains() ([]libvirtxml.Domain, error) {
	var result []libvirtxml.Domain

	flags := libvirt.CONNECT_LIST_DOMAINS_ACTIVE |
		libvirt.CONNECT_LIST_DOMAINS_INACTIVE |
		libvirt.CONNECT_LIST_DOMAINS_PERSISTENT |
		libvirt.CONNECT_LIST_DOMAINS_TRANSIENT |
		libvirt.CONNECT_LIST_DOMAINS_RUNNING |
		libvirt.CONNECT_LIST_DOMAINS_PAUSED |
		libvirt.CONNECT_LIST_DOMAINS_SHUTOFF |
		libvirt.CONNECT_LIST_DOMAINS_OTHER |
		libvirt.CONNECT_LIST_DOMAINS_MANAGEDSAVE |
		libvirt.CONNECT_LIST_DOMAINS_NO_MANAGEDSAVE |
		libvirt.CONNECT_LIST_DOMAINS_AUTOSTART |
		libvirt.CONNECT_LIST_DOMAINS_NO_AUTOSTART |
		libvirt.CONNECT_LIST_DOMAINS_HAS_SNAPSHOT |
		libvirt.CONNECT_LIST_DOMAINS_NO_SNAPSHOT

	domains, err := c.connect.ListAllDomains(flags)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list domains")
	}

	for _, domain := range domains {
		domainXML, err := getDomainXML(&domain)
		if err != nil {
			return nil, err
		}

		result = append(result, *domainXML)
		domain.Free()
	}

	//TODO(bcrusu): cache the result
	return result, nil
}

func (c *Connection) CreateStorageVolume(pool string, name string, backingVolumeName string) error {
	p, err := c.findStoragePool(pool)
	if err != nil {
		return err
	}

	return createStorageVolume(p, name, backingVolumeName)
}

func (c *Connection) findStoragePool(name string) (*libvirt.StoragePool, error) {
	pool, err := lookupStoragePool(c.connect, name)
	if err != nil {
		return nil, err
	}

	if pool == nil {
		return nil, errors.Errorf("could not find storage pool '%s'", pool)
	}

	return pool, nil
}

func (c *Connection) GenerateUniqueMACAddresses(count int) ([]string, error) {
	// TODO: take into account network bridge MAC address
	allDomains, err := c.ListAllDomains()
	if err != nil {
		return nil, err
	}

	allMACs := make(map[string]bool)
	for _, domain := range allDomains {
		interfaces := domain.Devices().Interfaces()
		for _, iface := range interfaces {
			ifaceMAC := iface.MACAddress()
			allMACs[strings.ToUpper(ifaceMAC)] = true
		}
	}

	var result []string

enclosingLoop:
	for j := 0; j < count; j++ {
		for i := 0; i < 256; i++ {
			mac, err := randomMACAddress(c.uri)
			if err != nil {
				return nil, err
			}

			if _, ok := allMACs[mac]; !ok {
				// if no colisions add to result; else retry, up to 256 times
				result = append(result, mac)

				// add to list to avoid conflicts between generated MACs
				allMACs[mac] = true

				continue enclosingLoop
			}
		}

		return nil, errors.Errorf("failed to generate non-conflicting MAC address. Too many colisions")
	}

	return result, nil
}

func (c *Connection) GetNetwork(network string) (*libvirtxml.Network, error) {
	d, err := lookupNetwork(c.connect, network)
	if err != nil {
		return nil, err
	}

	if d == nil {
		// not found
		return nil, nil
	}
	defer d.Free()

	return getNetworkXML(d)
}

func (c *Connection) DefineNATNetwork(params DefineNetworkParams) error {
	return defineNATNetwork(c.connect, params)
}

func (c *Connection) DefineDomain(params DefineDomainParams) error {
	capabilities, err := c.GetCapabilities()
	if err != nil {
		return err
	}

	qemuEmulatorPath, err := findQemuEmulatorPath(capabilities)
	if err != nil {
		return err
	}

	return defineDomain(c.connect, params, qemuEmulatorPath)
}
