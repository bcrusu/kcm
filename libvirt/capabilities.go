package libvirt

import (
	"github.com/bcrusu/kcm/libvirtxml"
	"github.com/libvirt/libvirt-go"
	"github.com/pkg/errors"
)

func getCapabilitiesXML(connect *libvirt.Connect) (*libvirtxml.Capabilities, error) {
	xml, err := connect.GetCapabilities()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch libvirt capabilities")
	}

	return libvirtxml.NewCapabilitiesForXML(xml)
}
