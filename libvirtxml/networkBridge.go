package libvirtxml

type NetworkBridge struct {
	node *Node
}

func newNetworkBridge(node *Node) NetworkBridge {
	return NetworkBridge{
		node: node,
	}
}

func (s NetworkBridge) Name() string {
	return s.node.getAttribute(nameForLocal("name"))
}

func (s NetworkBridge) SetName(value string) {
	s.node.setAttribute(nameForLocal("name"), value)
}

func (s NetworkBridge) STP() bool {
	stp := s.node.getAttribute(nameForLocal("stp"))
	return stp == "on"
}

func (s NetworkBridge) SetSTP(value bool) {
	if value {
		s.node.setAttribute(nameForLocal("stp"), "on")
	} else {
		s.node.setAttribute(nameForLocal("stp"), "off")
	}
}
