package libvirt

import (
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/bcrusu/kcm/libvirtxml"
	"github.com/pkg/errors"
)

const uuidStringLength = 36

var random *rand.Rand

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	random = rand.New(source)
}

func randomMACAddress(uri string) (string, error) {
	url, err := url.Parse(uri)
	if err != nil {
		return "", errors.Wrapf(err, "libvirt: failed to parse libvirt connection uri")
	}

	var mac []byte

	if isQemuURL(url) {
		mac = []byte{0x52, 0x54, 0x00}
	} else if isXenURL(url) {
		mac = []byte{0x00, 0x16, 0x3E}
	}

	for len(mac) < 6 {
		b := random.Uint32()
		mac = append(mac, byte(b))
	}

	result := fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
	return strings.ToUpper(result), nil
}

func isQemuURL(url *url.URL) bool {
	return strings.HasPrefix(url.Scheme, "qemu")
}

func isXenURL(url *url.URL) bool {
	return strings.HasPrefix(url.Scheme, "xen") ||
		strings.HasPrefix(url.Scheme, "libxl")
}

func setMetadataValues(metadata libvirtxml.Metadata, kv map[string]string) {
	for name, value := range kv {
		nodeName := libvirtxml.NewName(MetadataXMLNamespace, name)
		node := metadata.NewNode(nodeName)
		node.CharData = value
	}
}
