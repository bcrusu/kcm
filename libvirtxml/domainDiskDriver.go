package libvirtxml

type DomainDiskDriver struct {
	node *Node
}

func newDomainDiskDriver(node *Node) DomainDiskDriver {
	return DomainDiskDriver{
		node: node,
	}
}

func (s DomainDiskDriver) Name() string {
	return s.node.getAttribute(nameForLocal("name"))
}

func (s DomainDiskDriver) SetName(value string) {
	s.node.setAttribute(nameForLocal("name"), value)
}

func (s DomainDiskDriver) Type() string {
	return s.node.getAttribute(nameForLocal("type"))
}

func (s DomainDiskDriver) SetType(value string) {
	s.node.setAttribute(nameForLocal("type"), value)
}
