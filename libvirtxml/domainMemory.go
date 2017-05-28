package libvirtxml

import "strconv"

type DomainMemory struct {
	node *Node
}

func newDomainMemory(node *Node) DomainMemory {
	return DomainMemory{
		node: node,
	}
}

func (s DomainMemory) Unit() string {
	return s.node.getAttribute(nameForLocal("unit"))
}

func (s DomainMemory) SetUnit(value string) {
	s.node.setAttribute(nameForLocal("unit"), value)
}

func (s DomainMemory) Value() string {
	return s.node.CharData
}

func (s DomainMemory) SetValue(value uint64) {
	s.node.CharData = strconv.FormatUint(value, 10)
}
