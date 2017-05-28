package libvirtxml

import (
	"strconv"

	"github.com/golang/glog"
)

type StorageVolumeSize struct {
	node *Node
}

func newStorageVolumeSize(node *Node) StorageVolumeSize {
	return StorageVolumeSize{
		node: node,
	}
}

func (s StorageVolumeSize) Unit() string {
	return s.node.getAttribute(nameForLocal("unit"))
}

func (s StorageVolumeSize) SetUnit(value string) {
	s.node.setAttribute(nameForLocal("unit"), value)
}

func (s StorageVolumeSize) Value() uint64 {
	str := s.node.CharData
	if str == "" {
		return 0
	}

	result, err := strconv.Atoi(str)
	if err != nil || result < 0 {
		result = 0
		glog.Warningf("libvirtxml: ignoring invalid storage volume size '%s'", str)
	}

	return uint64(result)
}

func (s StorageVolumeSize) SetValue(value uint64) {
	s.node.CharData = strconv.FormatUint(value, 10)
}
