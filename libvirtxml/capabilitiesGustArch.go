package libvirtxml

type CapabilitiesGustArch struct {
	node *Node
}

func newCapabilitiesGustArch(node *Node) CapabilitiesGustArch {
	return CapabilitiesGustArch{
		node: node,
	}
}

func (s CapabilitiesGustArch) Name() string {
	return s.node.getAttribute(nameForLocal("name"))
}

func (s CapabilitiesGustArch) Emulator() string {
	return s.node.ensureNode(nameForLocal("emulator")).CharData
}
