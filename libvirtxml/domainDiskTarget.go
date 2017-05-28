package libvirtxml

type DomainDiskTarget struct {
	node *Node
}

func newDomainDiskTarget(node *Node) DomainDiskTarget {
	return DomainDiskTarget{
		node: node,
	}
}

func (s DomainDiskTarget) Dev() string {
	return s.node.getAttribute(nameForLocal("dev"))
}

func (s DomainDiskTarget) SetDev(value string) {
	s.node.setAttribute(nameForLocal("dev"), value)
}

func (s DomainDiskTarget) Bus() string {
	return s.node.getAttribute(nameForLocal("bus"))
}

func (s DomainDiskTarget) SetBus(value string) {
	s.node.setAttribute(nameForLocal("bus"), value)
}
