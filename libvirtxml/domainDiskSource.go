package libvirtxml

type DomainDiskSource struct {
	node *Node
}

func newDomainDiskSource(node *Node) DomainDiskSource {
	return DomainDiskSource{
		node: node,
	}
}

func (s DomainDiskSource) File() string {
	return s.node.getAttribute(nameForLocal("file"))
}

func (s DomainDiskSource) SetFile(value string) {
	s.node.setAttribute(nameForLocal("file"), value)
}
