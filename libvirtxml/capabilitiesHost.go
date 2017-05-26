package libvirtxml

type CapabilitiesHost struct {
	node *Node
}

func newCapabilitiesHost(node *Node) CapabilitiesHost {
	return CapabilitiesHost{
		node: node,
	}
}

func (s CapabilitiesHost) UUID() string {
	return s.node.ensureNode(nameForLocal("uuid")).CharData
}

func (s CapabilitiesHost) CPU() CapabilitiesHostCPU {
	node := s.node.ensureNode(nameForLocal("cpu"))
	return newCapabilitiesHostCPU(node)
}
