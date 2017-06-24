package libvirt

import (
	"github.com/bcrusu/kcm/libvirtxml"
	"github.com/golang/glog"
	"github.com/libvirt/libvirt-go"
	"github.com/pkg/errors"
)

func lookupStoragePool(connect *libvirt.Connect, lookup string) (*libvirt.StoragePool, error) {
	if len(lookup) == uuidStringLength {
		pool, err := connect.LookupStoragePoolByUUIDString(lookup)
		if err != nil {
			if lverr, ok := err.(libvirt.Error); ok && lverr.Code != libvirt.ERR_NO_STORAGE_POOL {
				glog.Infof("libvirt: storage pool lookup by ID '%s' failed. Error: %v", lookup, lverr)
			}
		}

		if pool != nil {
			return pool, nil
		}
	}

	pool, err := connect.LookupStoragePoolByName(lookup)
	if err != nil {
		if lverr, ok := err.(libvirt.Error); ok && lverr.Code == libvirt.ERR_NO_STORAGE_POOL {
			return nil, nil
		}

		return nil, errors.Wrapf(err, "libvirt: storage pool lookup failed '%s'", lookup)
	}

	return pool, nil
}

func lookupStoragePoolStrict(connect *libvirt.Connect, lookup string) (*libvirt.StoragePool, error) {
	pool, err := lookupStoragePool(connect, lookup)
	if err != nil {
		return nil, err
	}

	if pool == nil {
		return nil, errors.Errorf("libvirt: could not find storage pool '%s'", lookup)
	}

	return pool, nil
}

func getStoragePoolXML(pool *libvirt.StoragePool) (*libvirtxml.StoragePool, error) {
	xml, err := pool.GetXMLDesc(libvirt.StorageXMLFlags(0))
	if err != nil {
		return nil, errors.Wrapf(err, "libvirt: failed to fetch storage pool XML description")
	}

	return libvirtxml.NewStoragePoolForXML(xml)
}
