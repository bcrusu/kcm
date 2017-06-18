package libvirt

import (
	"github.com/bcrusu/kcm/libvirtxml"
	"github.com/libvirt/libvirt-go"
	"github.com/pkg/errors"
)

type CreateStorageVolumeParams struct {
	Pool        string
	Name        string
	CapacityGiB uint

	// only one should be set
	BackingVolumeName string
	Content           []byte
}

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

func createStorageVolumeFromBackingVolume(pool *libvirt.StoragePool, params CreateStorageVolumeParams) (*libvirtxml.StorageVolume, error) {
	backingVolume, err := lookupStorageVolume(pool, params.BackingVolumeName)
	{
		if err != nil {
			return nil, err
		}

		if backingVolume == nil {
			return nil, errors.Errorf("could not find storage volume '%s'", params.BackingVolumeName)
		}
		defer backingVolume.Free()
	}

	volumeXML, err := getStorageVolumeXML(backingVolume)
	{
		// create the new volume XML definition starting from the backing volume definition
		if err != nil {
			return nil, err
		}

		volumeType := volumeXML.Type()
		if volumeType != "file" {
			errors.Errorf("cannot clone storage volume '%s' - unsupported volume type '%s'", params.BackingVolumeName, volumeType)
		}

		volumeXML.SetName(params.Name)
		volumeXML.SetKey("")

		volumeXML.Capacity().SetUnit("GiB")
		volumeXML.Capacity().SetValue(uint64(params.CapacityGiB))

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
		return nil, err
	}

	storageVol, err := pool.StorageVolCreateXML(xmlString, libvirt.StorageVolCreateFlags(0))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to clone storage volume '%s' to '%s'", params.BackingVolumeName, params.Name)
	}
	defer storageVol.Free()

	return getStorageVolumeXML(storageVol)
}

func createStorageVolumeFromContent(connect *libvirt.Connect, pool *libvirt.StoragePool, params CreateStorageVolumeParams) (*libvirtxml.StorageVolume, error) {
	storageVol, err := createEmptyStorageVolume(pool, params.Name, params.CapacityGiB)
	if err != nil {
		return nil, err
	}
	defer storageVol.Free()

	stream, err := connect.NewStream(libvirt.STREAM_NONBLOCK)
	if err != nil {
		return nil, err
	}
	defer stream.Free()

	if err := storageVol.Upload(stream, 0, 0, libvirt.StorageVolUploadFlags(0)); err != nil {
		return nil, err
	}

	if err := streamSendAll(stream, params.Content); err != nil {
		return nil, errors.Wrap(err, "libvirt: failed to upload storage volume content")
	}

	return getStorageVolumeXML(storageVol)
}

func createEmptyStorageVolume(pool *libvirt.StoragePool, name string, capacityGiB uint) (*libvirt.StorageVol, error) {
	volumeXML := libvirtxml.NewStorageVolume()
	volumeXML.SetType("file")
	volumeXML.SetName(name)
	volumeXML.Target().Format().SetType("qcow2")

	volumeXML.Capacity().SetUnit("GiB")
	volumeXML.Capacity().SetValue(uint64(capacityGiB))
	volumeXML.Allocation().SetUnit("bytes")
	volumeXML.Allocation().SetValue(0)

	xmlString, err := volumeXML.MarshalToXML()
	if err != nil {
		return nil, err
	}

	storageVol, err := pool.StorageVolCreateXML(xmlString, libvirt.StorageVolCreateFlags(0))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create empty storage volume '%s'", name)
	}

	return storageVol, nil
}
