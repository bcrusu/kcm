package libvirt

import (
	"github.com/bcrusu/kcm/libvirtxml"
	"github.com/libvirt/libvirt-go"
	"github.com/pkg/errors"
)

func lookupStorageVolume(pool *libvirt.StoragePool, lookup string) (*libvirt.StorageVol, error) {
	volume, err := pool.LookupStorageVolByName(lookup)
	if err != nil {
		if lverr, ok := err.(libvirt.Error); ok && lverr.Code == libvirt.ERR_NO_STORAGE_VOL {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "storage volume lookup failed '%s'", lookup)
	}

	return volume, nil
}

func getStorageVolumeXML(volume *libvirt.StorageVol) (*libvirtxml.StorageVolume, error) {
	xml, err := volume.GetXMLDesc(0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch storage volume XML description")
	}

	return libvirtxml.NewStorageVolumeForXML(xml)
}

func createStorageVolume(pool *libvirt.StoragePool, name string, backingVolumeName string) error {
	backingVolume, err := lookupStorageVolume(pool, backingVolumeName)
	{
		if err != nil {
			return err
		}

		if backingVolume == nil {
			return errors.Errorf("could not find storage volume '%s'", backingVolumeName)
		}
		defer backingVolume.Free()
	}

	volumeXML, err := getStorageVolumeXML(backingVolume)
	{
		// create the new volume XML definition starting from the backing volume definition
		if err != nil {
			return err
		}

		volumeType := volumeXML.Type()
		if volumeType != "file" {
			errors.Errorf("cannot clone storage volume '%s' - unsupported volume type '%s'", backingVolumeName, volumeType)
		}

		volumeXML.SetName(name)
		volumeXML.SetKey("")

		targetXML := volumeXML.Target()
		targetXML.RemoveTimestamps()

		sourcePath := targetXML.Path()
		targetXML.SetPath("") // will be filled-in by libvirt

		{
			// set backing store as the souorce target
			backingStoreXML := volumeXML.BackingStore()
			backingStoreXML.SetPath(sourcePath)
			backingStoreXML.Format().SetType(targetXML.Format().Type())
			backingStoreXML.RemoveTimestamps()
		}

		// switch to a format that supports backing store
		switch targetXML.Format().Type() {
		case "raw":
			targetXML.Format().SetType("qcow2")
		}
	}

	xmlString, err := volumeXML.MarshalToXML()
	if err != nil {
		return err
	}

	storageVol, err := pool.StorageVolCreateXML(xmlString, libvirt.StorageVolCreateFlags(0))
	if err != nil {
		return errors.Wrapf(err, "failed to clone storage volume '%s' to '%s'", backingVolumeName, name)
	}
	defer storageVol.Free()

	return err
}
