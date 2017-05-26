package libvirtxml

type StorageVolumeTarget struct {
	node *Node
}

func newStorageVolumeTarget(node *Node) StorageVolumeTarget {
	return StorageVolumeTarget{
		node: node,
	}
}

func (s StorageVolumeTarget) Path() string {
	return s.node.ensureNode(nameForLocal("path")).CharData
}

func (s StorageVolumeTarget) SetPath(value string) {
	s.node.ensureNode(nameForLocal("path")).CharData = value
}

func (s StorageVolumeTarget) RemoveTimestamps() {
	s.node.removeNodes(nameForLocal("timestamps"))
}

func (s StorageVolumeTarget) Format() StorageVolumeTargetFormat {
	node := s.node.ensureNode(nameForLocal("format"))
	return newStorageVolumeTargetFormat(node)
}
