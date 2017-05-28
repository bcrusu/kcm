package libvirt

import (
	"github.com/bcrusu/kcm/libvirtxml"
	"github.com/golang/glog"
	"github.com/libvirt/libvirt-go"
	"github.com/pkg/errors"
)

type DefineDomainParams struct {
	Name                string
	Network             string
	NetworkInterfaceMAC string
	StoragePool         string
	StorageVolume       string
	FilesystemMounts    map[string]string // map[HOST_PATH]GUEST_PATH
	MemoryMiB           uint              // max domain memory
	CPUs                uint              // number of CPU cores
	Metadata            map[string]string // map[NAME]VALUE
}

func lookupDomain(connect *libvirt.Connect, lookup string) (*libvirt.Domain, error) {
	if len(lookup) == uuidStringLength {
		domain, err := connect.LookupDomainByUUIDString(lookup)
		if err != nil {
			if lverr, ok := err.(libvirt.Error); ok && lverr.Code != libvirt.ERR_NO_DOMAIN {
				glog.Infof("domain lookup by ID '%s' failed. Error: %v", lookup, lverr)
			}
		}

		if domain != nil {
			return domain, nil
		}
	}

	domain, err := connect.LookupDomainByName(lookup)
	if err != nil {
		if lverr, ok := err.(libvirt.Error); ok && lverr.Code == libvirt.ERR_NO_DOMAIN {
			return nil, nil

		}
		return nil, errors.Wrapf(err, "domain lookup failed '%s'", lookup)
	}

	return domain, nil
}

func getDomainXML(domain *libvirt.Domain) (*libvirtxml.Domain, error) {
	xml, err := domain.GetXMLDesc(libvirt.DomainXMLFlags(0))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch domain XML description")
	}

	return libvirtxml.NewDomainForXML(xml)
}

func defineDomain(connect *libvirt.Connect, params DefineDomainParams, emulatorPath string) error {
	domainXML, err := libvirtxml.NewDomainForXML(defaultDomainTemplateXML)
	if err != nil {
		return err
	}

	domainXML.SetID("")
	domainXML.SetUUID("")
	domainXML.SetName(params.Name)

	domainXML.VCPU().SetValue(params.CPUs)
	domainXML.Memory().SetUnit("MiB")
	domainXML.Memory().SetValue(uint64(params.MemoryMiB))

	for hostPath, guestPath := range params.FilesystemMounts {
		fs := domainXML.Devices().NewFilesystem()
		fs.SetType("mount")
		fs.SetAccessmode("squash")
		fs.SetSourceDir(hostPath)
		fs.SetTargetDir(guestPath)
		fs.SetReadonly(true)
	}

	{
		iface := domainXML.Devices().NewInterface()
		iface.SetType("network")
		iface.SetMACAddress(params.NetworkInterfaceMAC)
		iface.Source().SetNetwork(params.Network)
		iface.SetModelType("virtio")
	}

	{
		disk := domainXML.Devices().NewDisk()
		disk.SetType("volume")
		disk.SetDevice("disk")

		disk.Driver().SetName("qemu")
		disk.Driver().SetType("qcow2")

		disk.Source().SetPool(params.StoragePool)
		disk.Source().SetVolume(params.StorageVolume)

		disk.Target().SetDev("vda")
		disk.Target().SetBus("virtio")
	}

	setMetadataValues(domainXML.Metadata(), params.Metadata)

	// Set the graphics device port to auto, in order to avoid conflicts
	graphics := domainXML.Devices().Graphics()
	for _, graphic := range graphics {
		graphic.SetPort(-1)
	}

	// reset path for guest agent channel
	channels := domainXML.Devices().Channels()
	for _, channel := range channels {
		if channel.Type() != "unix" {
			continue
		}

		// will be set by libvirt
		channel.SetSourcePath("")
	}

	domainXML.Devices().SetEmulator(emulatorPath)

	xml, err := domainXML.MarshalToXML()
	if err != nil {
		return err
	}

	domain, err := connect.DomainDefineXML(xml)
	if err != nil {
		return errors.Wrapf(err, "libvirt: failed to define domain '%s'", params.Name)
	}
	defer domain.Free()

	return nil
}
