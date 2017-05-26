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
