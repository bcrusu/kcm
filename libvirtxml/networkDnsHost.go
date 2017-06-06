package libvirtxml

type NetworkDNSHost struct {
	node *Node
}

func newNetworkDNSHost(node *Node) NetworkDNSHost {
	return NetworkDNSHost{
		node: node,
	}
}

func (s NetworkDNSHost) IP() string {
	return s.node.getAttribute(nameForLocal("ip"))
}

func (s NetworkDNSHost) SetEnable(value string) {
	s.node.setAttribute(nameForLocal("ip"), value)
}

func (s Network) Hostnames() []string {
	var result []string

	nodes := s.root.findNodes(nameForLocal("hostname"))
	for _, node := range nodes {
		result = append(result, node.CharData)
	}

	return result
}

func (s Network) SetHostnames(values []string) {
	name := nameForLocal("hostname")
	s.root.removeNodes(name)

	for _, value := range values {
		node := NewNode(name)
		node.CharData = value
		s.root.addNode(node)
	}
}
