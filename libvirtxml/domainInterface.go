package libvirtxml

type DomainInterface struct {
	node *Node
}

func newDomainInterface(node *Node) DomainInterface {
	return DomainInterface{
		node: node,
	}
}

func (s DomainInterface) TargetDevice() string {
	node := s.node.ensureNode(nameForLocal("target"))
	return node.getAttribute(nameForLocal("dev"))
}

func (s DomainInterface) SetTargetDevice(value string) {
	node := s.node.ensureNode(nameForLocal("target"))
	node.setAttribute(nameForLocal("dev"), value)
}

func (s DomainInterface) MACAddress() string {
	node := s.node.ensureNode(nameForLocal("mac"))
	return node.getAttribute(nameForLocal("address"))
}

func (s DomainInterface) SetMACAddress(value string) {
	node := s.node.ensureNode(nameForLocal("mac"))
	node.setAttribute(nameForLocal("address"), value)
}
