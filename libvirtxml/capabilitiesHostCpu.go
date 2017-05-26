package libvirtxml

type CapabilitiesHostCPU struct {
	node *Node
}

func newCapabilitiesHostCPU(node *Node) CapabilitiesHostCPU {
	return CapabilitiesHostCPU{
		node: node,
	}
}

func (s CapabilitiesHostCPU) Arch() string {
	return s.node.ensureNode(nameForLocal("arch")).CharData
}
