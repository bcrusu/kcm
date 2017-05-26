package libvirtxml

type StoragePoolTarget struct {
	node *Node
}

func newStoragePoolTarget(node *Node) StoragePoolTarget {
	return StoragePoolTarget{
		node: node,
	}
}

func (s StoragePoolTarget) Path() string {
	return s.node.ensureNode(nameForLocal("path")).CharData
}

func (s StoragePoolTarget) SetPath(value string) {
	s.node.ensureNode(nameForLocal("path")).CharData = value
}
