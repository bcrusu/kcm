package libvirtxml

import (
	"strconv"
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

func (s StorageVolumeSize) Value() string {
	return s.node.CharData
}

func (s StorageVolumeSize) SetValue(value uint64) {
	s.node.CharData = strconv.FormatUint(value, 10)
}
