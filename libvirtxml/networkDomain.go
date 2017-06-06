package libvirtxml

type NetworDomain struct {
	node *Node
}

func newNetworkDomain(node *Node) NetworDomain {
	return NetworDomain{
		node: node,
	}
}

func (s NetworDomain) Name() string {
	return s.node.getAttribute(nameForLocal("name"))
}

func (s NetworDomain) SetName(value string) {
	s.node.setAttribute(nameForLocal("name"), value)
}

func (s NetworDomain) LocalOnly() bool {
	stp := s.node.getAttribute(nameForLocal("localOnly"))
	return stp == "yes"
}

func (s NetworDomain) SetLocalOnly(value bool) {
	if value {
		s.node.setAttribute(nameForLocal("localOnly"), "yes")
	} else {
		s.node.setAttribute(nameForLocal("localOnly"), "no")
	}
}
