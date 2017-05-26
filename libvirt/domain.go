package libvirt

import (
	"github.com/bcrusu/kcm/libvirtxml"
	"github.com/golang/glog"
	"github.com/libvirt/libvirt-go"
	"github.com/pkg/errors"
)

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
