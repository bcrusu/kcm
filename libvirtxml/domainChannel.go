package libvirtxml

type DomainChannel struct {
	node *Node
}

func newDomainChannel(node *Node) DomainChannel {
	return DomainChannel{
		node: node,
	}
}

func (s DomainChannel) Type() string {
	return s.node.getAttribute(nameForLocal("type"))
}

func (s DomainChannel) SetType(value string) {
	s.node.setAttribute(nameForLocal("type"), value)
}

func (s DomainChannel) SourcePath() string {
	node := s.node.ensureNode(nameForLocal("source"))
	return node.getAttribute(nameForLocal("path"))
}

func (s DomainChannel) SetSourcePath(value string) {
	node := s.node.ensureNode(nameForLocal("source"))
	node.setAttribute(nameForLocal("path"), value)
}
