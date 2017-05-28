package libvirtxml

type DomainDisk struct {
	node *Node
}

func newDomainDisk(node *Node) DomainDisk {
	return DomainDisk{
		node: node,
	}
}

func (s DomainDisk) Type() string {
	return s.node.getAttribute(nameForLocal("type"))
}

func (s DomainDisk) SetType(value string) {
	s.node.setAttribute(nameForLocal("type"), value)
}

func (s DomainDisk) Device() string {
	return s.node.getAttribute(nameForLocal("device"))
}

func (s DomainDisk) SetDevice(value string) {
	s.node.setAttribute(nameForLocal("device"), value)
}

func (s DomainDisk) Readonly() bool {
	return s.node.hasNode(nameForLocal("readonly"))
}

func (s DomainDisk) SetReadonly(value bool) {
	if value {
		s.node.ensureNode(nameForLocal("readonly"))
	} else {
		s.node.removeNodes(nameForLocal("readonly"))
	}
}

func (s DomainDisk) Shareable() bool {
	return s.node.hasNode(nameForLocal("shareable"))
}

func (s DomainDisk) SetShareable(value bool) {
	if value {
		s.node.ensureNode(nameForLocal("shareable"))
	} else {
		s.node.removeNodes(nameForLocal("shareable"))
	}
}

func (s DomainDisk) Source() DomainDiskSource {
	node := s.node.ensureNode(nameForLocal("source"))
	return newDomainDiskSource(node)
}

func (s DomainDisk) Target() DomainDiskTarget {
	node := s.node.ensureNode(nameForLocal("target"))
	return newDomainDiskTarget(node)
}

func (s DomainDisk) Driver() DomainDiskDriver {
	node := s.node.ensureNode(nameForLocal("driver"))
	return newDomainDiskDriver(node)
}
