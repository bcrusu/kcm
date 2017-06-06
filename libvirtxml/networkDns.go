package libvirtxml

type NetworkDNS struct {
	node *Node
}

func newNetworkDNS(node *Node) NetworkDNS {
	return NetworkDNS{
		node: node,
	}
}

func (s NetworkDNS) Enable() bool {
	stp := s.node.getAttribute(nameForLocal("enable"))
	return stp == "yes"
}

func (s NetworkDNS) SetEnable(value bool) {
	if value {
		s.node.setAttribute(nameForLocal("enable"), "yes")
	} else {
		s.node.setAttribute(nameForLocal("enable"), "no")
	}
}

func (s NetworkDNS) ForwardPlainNames() bool {
	stp := s.node.getAttribute(nameForLocal("forwardPlainNames"))
	return stp == "yes"
}

func (s NetworkDNS) SetForwardPlainNames(value bool) {
	if value {
		s.node.setAttribute(nameForLocal("forwardPlainNames"), "yes")
	} else {
		s.node.setAttribute(nameForLocal("forwardPlainNames"), "no")
	}
}

func (s Network) Hosts() []NetworkDNSHost {
	var result []NetworkDNSHost

	nodes := s.root.findNodes(nameForLocal("host"))
	for _, node := range nodes {
		result = append(result, newNetworkDNSHost(node))
	}

	return result
}
