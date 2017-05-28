package libvirtxml

import "strconv"

type NetworkForward struct {
	node *Node
}

func newNetworkForward(node *Node) NetworkForward {
	return NetworkForward{
		node: node,
	}
}

func (s NetworkForward) Mode() string {
	return s.node.getAttribute(nameForLocal("mode"))
}

func (s NetworkForward) SetMode(value string) {
	s.node.setAttribute(nameForLocal("mode"), value)
}

func (s NetworkForward) SetNATPortRange(start int, end int) {
	nat := s.node.ensureNode(nameForLocal("nat"))
	port := nat.ensureNode(nameForLocal("port"))

	startStr := strconv.FormatInt(int64(start), 10)
	endStr := strconv.FormatInt(int64(end), 10)

	port.setAttribute(nameForLocal("start"), startStr)
	port.setAttribute(nameForLocal("end"), endStr)
}
