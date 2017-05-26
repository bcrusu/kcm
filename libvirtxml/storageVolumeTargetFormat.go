package libvirtxml

type StorageVolumeTargetFormat struct {
	node *Node
}

func newStorageVolumeTargetFormat(node *Node) StorageVolumeTargetFormat {
	return StorageVolumeTargetFormat{
		node: node,
	}
}

func (s StorageVolumeTargetFormat) Type() string {
	return s.node.getAttribute(nameForLocal("type"))
}

func (s StorageVolumeTargetFormat) SetType(value string) {
	s.node.setAttribute(nameForLocal("type"), value)
}
