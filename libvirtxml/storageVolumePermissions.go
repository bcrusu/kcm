package libvirtxml

import (
	"strconv"

	"github.com/golang/glog"
)

type StorageVolumePermissions struct {
	node *Node
}

func newStorageVolumePermissions(node *Node) StorageVolumePermissions {
	return StorageVolumePermissions{
		node: node,
	}
}

func (s StorageVolumePermissions) Owner() uint64 {
	node := s.node.ensureNode(nameForLocal("owner"))

	str := node.CharData
	if str == "" {
		return 0
	}

	result, err := strconv.Atoi(str)
	if err != nil || result < 0 {
		result = 0
		glog.Warningf("libvirtxml: ignoring invalid storage volume owner '%s'", str)
	}

	return uint64(result)
}

func (s StorageVolumePermissions) SetOwner(value uint64) {
	node := s.node.ensureNode(nameForLocal("owner"))
	node.CharData = strconv.FormatUint(value, 10)
}

func (s StorageVolumePermissions) Group() uint64 {
	node := s.node.ensureNode(nameForLocal("group"))

	str := node.CharData
	if str == "" {
		return 0
	}

	result, err := strconv.Atoi(str)
	if err != nil || result < 0 {
		result = 0
		glog.Warningf("libvirtxml: ignoring invalid storage volume group '%s'", str)
	}

	return uint64(result)
}

func (s StorageVolumePermissions) SetGroup(value uint64) {
	node := s.node.ensureNode(nameForLocal("group"))
	node.CharData = strconv.FormatUint(value, 10)
}

func (s StorageVolumePermissions) Mode() string {
	return s.node.ensureNode(nameForLocal("mode")).CharData
}

func (s StorageVolumePermissions) SetMode(value string) {
	s.node.ensureNode(nameForLocal("mode")).CharData = value
}
