package libvirtxml

type DomainInterfaceSource struct {
	node *Node
}

func newDomainInterfaceSource(node *Node) DomainInterfaceSource {
	return DomainInterfaceSource{
		node: node,
	}
}

func (s DomainInterfaceSource) Network() string {
	return s.node.getAttribute(nameForLocal("network"))
}

func (s DomainInterfaceSource) SetNetwork(value string) {
	s.node.setAttribute(nameForLocal("network"), value)
}
