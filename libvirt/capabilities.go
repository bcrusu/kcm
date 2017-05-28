package libvirt

import (
	"strings"

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

func findQemuEmulatorPath(capabilities *libvirtxml.Capabilities) (string, error) {
	hostArch := capabilities.Host().CPU().Arch()
	guests := capabilities.Guests()

	var emulator string
	for _, guest := range guests {
		if guest.Arch().Name() != hostArch {
			continue
		}

		emulator = guest.Arch().Emulator()
		if strings.Contains(strings.ToLower(emulator), "qemu") {
			continue
		}

		if guest.OSType() == "hvm" {
			// found hardware-assisted vm - use this emulator
			break
		}
	}

	if emulator == "" {
		return "", errors.Errorf("libvirt: found no guest matching host architecture '%s'", hostArch)
	}

	return emulator, nil
}
