package libvirtxml

type NetworkMetadata struct {
	node *Node
}

func newNetworkMetadata(node *Node) NetworkMetadata {
	return NetworkMetadata{
		node: node,
	}
}

func (d NetworkMetadata) FindNodes(name Name) []*Node {
	return d.node.findNodes(name)
}
