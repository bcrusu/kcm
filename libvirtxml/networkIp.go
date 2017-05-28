package libvirtxml

import (
	"strconv"

	"github.com/golang/glog"
)

type NetworkIP struct {
	node *Node
}

func newNetworkIP(node *Node) NetworkIP {
	return NetworkIP{
		node: node,
	}
}

func (s NetworkIP) Address() string {
	return s.node.getAttribute(nameForLocal("address"))
}

func (s NetworkIP) SetAddress(value string) {
	s.node.setAttribute(nameForLocal("address"), value)
}

func (s NetworkIP) Family() string {
	return s.node.getAttribute(nameForLocal("family"))
}

func (s NetworkIP) SetFamily(value string) {
	s.node.setAttribute(nameForLocal("family"), value)
}

func (s NetworkIP) Netmask() string {
	return s.node.getAttribute(nameForLocal("netmask"))
}

func (s NetworkIP) SetNetmask(value string) {
	s.node.setAttribute(nameForLocal("netmask"), value)
}

func (s NetworkIP) Prefix() int {
	str := s.node.getAttribute(nameForLocal("prefix"))
	if str == "" {
		return 0
	}

	prefix, err := strconv.Atoi(str)
	if err != nil {
		prefix = 0
		glog.Warningf("libvirtxml: ignoring invalid network IP prefix '%s'", str)
	}

	return prefix
}

func (s NetworkIP) SetPrefix(value int) {
	str := strconv.FormatInt(int64(value), 10)
	s.node.setAttribute(nameForLocal("prefix"), str)
}

func (s NetworkIP) SetDHCPRange(start string, end string) {
	dhcp := s.node.ensureNode(nameForLocal("dhcp"))
	rng := dhcp.ensureNode(nameForLocal("range"))

	rng.setAttribute(nameForLocal("start"), start)
	rng.setAttribute(nameForLocal("end"), end)
}
