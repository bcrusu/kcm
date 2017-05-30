package libvirtxml

type StorageVolumeBackingStore struct {
	node *Node
}

func newStorageVolumeBackingStore(node *Node) StorageVolumeBackingStore {
	return StorageVolumeBackingStore{
		node: node,
	}
}

func (s StorageVolumeBackingStore) Path() string {
	return s.node.ensureNode(nameForLocal("path")).CharData
}

func (s StorageVolumeBackingStore) SetPath(value string) {
	s.node.ensureNode(nameForLocal("path")).CharData = value
}

func (s StorageVolumeBackingStore) RemoveTimestamps() {
	s.node.removeNodes(nameForLocal("timestamps"))
}

func (s StorageVolumeBackingStore) Format() StorageVolumeTargetFormat {
	node := s.node.ensureNode(nameForLocal("format"))
	return newStorageVolumeTargetFormat(node)
}

func (s StorageVolumeBackingStore) Permissions() StorageVolumePermissions {
	node := s.node.ensureNode(nameForLocal("permissions"))
	return newStorageVolumePermissions(node)
}
