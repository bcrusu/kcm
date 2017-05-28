package libvirtxml

type DomainMetadata struct {
	node *Node
}

func newDomainMetadata(node *Node) DomainMetadata {
	return DomainMetadata{
		node: node,
	}
}

func (d DomainMetadata) FindNodes(name Name) []*Node {
	return d.node.findNodes(name)
}

func (s DomainMetadata) NewNode(name Name) *Node {
	node := NewNode(name)
	s.node.addNode(node)
	return node
}
