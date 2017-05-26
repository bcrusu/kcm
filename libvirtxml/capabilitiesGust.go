package libvirtxml

type CapabilitiesGuest struct {
	node *Node
}

func newCapabilitiesGuest(node *Node) CapabilitiesGuest {
	return CapabilitiesGuest{
		node: node,
	}
}

func (s CapabilitiesGuest) OSType() string {
	return s.node.ensureNode(nameForLocal("os_type")).CharData
}

func (s CapabilitiesGuest) Arch() CapabilitiesGustArch {
	node := s.node.ensureNode(nameForLocal("arch"))
	return newCapabilitiesGustArch(node)
}
