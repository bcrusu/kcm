package libvirtxml

import (
	"strconv"

	"github.com/golang/glog"
)

type DomainGraphic struct {
	node *Node
}

func newDomainGraphic(node *Node) DomainGraphic {
	return DomainGraphic{
		node: node,
	}
}

func (s DomainGraphic) Port() int {
	str := s.node.getAttribute(nameForLocal("port"))
	if str == "" {
		return -1
	}

	port, err := strconv.Atoi(str)
	if err != nil {
		port = 0
		glog.Warningf("libvirtxml: ignoring invalid domain graphics port '%s'", str)
	}
	return port
}

func (s DomainGraphic) SetPort(value int) {
	str := strconv.FormatInt(int64(value), 10)
	s.node.setAttribute(nameForLocal("port"), str)
}
