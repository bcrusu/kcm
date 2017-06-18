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
	p, err := lookupStoragePoolStrict(c.connect, pool)
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
	return listAllDomains(c.connect)
}

func (c *Connection) CreateStorageVolume(params CreateStorageVolumeParams) (*libvirtxml.StorageVolume, error) {
	if params.Content != nil && params.BackingVolumeName != "" {
		return nil, errors.New("libvirt: cannot create storage volume: invalid params")
	}

	pool, err := lookupStoragePoolStrict(c.connect, params.Pool)
	if err != nil {
		return nil, err
	}
	defer pool.Free()

	if params.Content != nil {
		return createStorageVolumeFromContent(c.connect, pool, params)
	}

	return createStorageVolumeFromBackingVolume(pool, params)
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

func (c *Connection) UndefineDomain(name string) error {
	domain, err := lookupDomainStrict(c.connect, name)
	if err != nil {
		return err
	}
	defer domain.Free()

	flags := libvirt.DOMAIN_UNDEFINE_MANAGED_SAVE |
		libvirt.DOMAIN_UNDEFINE_NVRAM |
		libvirt.DOMAIN_UNDEFINE_SNAPSHOTS_METADATA

	if err := domain.UndefineFlags(flags); err != nil {
		return errors.Wrapf(err, "libvirt: failed to undefine domain '%s'", name)
	}

	return nil
}

func (c *Connection) DestroyDomain(name string) error {
	domain, err := lookupDomainStrict(c.connect, name)
	if err != nil {
		return err
	}
	defer domain.Free()

	if err := domain.DestroyFlags(libvirt.DomainDestroyFlags(0)); err != nil {
		return errors.Wrapf(err, "libvirt: failed to destroy domain '%s'", name)
	}

	return nil
}

func (c *Connection) CreateDomain(name string) error {
	domain, err := lookupDomainStrict(c.connect, name)
	if err != nil {
		return err
	}
	defer domain.Free()

	if err := domain.Create(); err != nil {
		return errors.Wrapf(err, "libvirt: failed to create domain '%s'", name)
	}

	return nil
}

func (c *Connection) ShutdownDomain(name string) error {
	domain, err := lookupDomainStrict(c.connect, name)
	if err != nil {
		return err
	}
	defer domain.Free()

	if err := domain.Shutdown(); err != nil {
		return errors.Wrapf(err, "libvirt: failed to shutdown domain '%s'", name)
	}

	return nil
}

func (c *Connection) DomainIsActive(name string) (bool, error) {
	domain, err := lookupDomainStrict(c.connect, name)
	if err != nil {
		return false, err
	}
	defer domain.Free()

	active, err := domain.IsActive()

	if err != nil {
		return false, errors.Wrapf(err, "libvirt: failed to determine if domain is active '%s'", name)
	}

	return active, nil
}

func (c *Connection) DeleteStorageVolume(pool, name string) error {
	p, err := lookupStoragePoolStrict(c.connect, pool)
	if err != nil {
		return err
	}
	defer p.Free()

	volume, err := lookupStorageVolume(p, name)
	if err != nil {
		return err
	}

	if volume == nil {
		return nil
	}
	defer volume.Free()

	if err := volume.Delete(libvirt.STORAGE_VOL_DELETE_NORMAL); err != nil {
		return errors.Wrapf(err, "libvirt: failed to delete storage volume '%s'", name)
	}

	return nil
}

func (c *Connection) NetworkIsActive(name string) (bool, error) {
	n, err := lookupNetworkStrict(c.connect, name)
	if err != nil {
		return false, err
	}
	defer n.Free()

	active, err := n.IsActive()

	if err != nil {
		return false, errors.Wrapf(err, "libvirt: failed to determine if network active '%s'", name)
	}

	return active, nil
}

func (c *Connection) DestroyNetwork(name string) error {
	n, err := lookupNetworkStrict(c.connect, name)
	if err != nil {
		return err
	}
	defer n.Free()

	if err := n.Destroy(); err != nil {
		return errors.Wrapf(err, "libvirt: failed to destroy network '%s'", name)
	}

	return nil
}

func (c *Connection) CreateNetwork(name string) error {
	n, err := lookupNetworkStrict(c.connect, name)
	if err != nil {
		return err
	}
	defer n.Free()

	if err := n.Create(); err != nil {
		return errors.Wrapf(err, "libvirt: failed to create network '%s'", name)
	}

	return nil
}

func (c *Connection) UndefineNetwork(name string) error {
	n, err := lookupNetworkStrict(c.connect, name)
	if err != nil {
		return err
	}
	defer n.Free()

	if err := n.Undefine(); err != nil {
		return errors.Wrapf(err, "libvirt: failed to undefine network '%s'", name)
	}

	return nil
}

func (c *Connection) ListDomainInterfaceAddresses(domainName string) ([]string, error) {
	return listDomainInterfaceAddresses(c.connect, domainName)
}
