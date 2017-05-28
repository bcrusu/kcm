package libvirtxml

type Metadata struct {
	node *Node
}

func newMetadata(node *Node) Metadata {
	return Metadata{
		node: node,
	}
}

func (d Metadata) FindNodes(name Name) []*Node {
	return d.node.findNodes(name)
}

func (s Metadata) NewNode(name Name) *Node {
	node := NewNode(name)
	s.node.addNode(node)
	return node
}
