package libvirtxml

import (
	"strconv"

	"github.com/golang/glog"
)

type DomainVCPU struct {
	node *Node
}

func newDomainVCPU(node *Node) DomainVCPU {
	return DomainVCPU{
		node: node,
	}
}

func (s DomainVCPU) Value() uint {
	str := s.node.CharData
	if str == "" {
		return 0
	}

	result, err := strconv.Atoi(str)
	if err != nil || result < 0 {
		result = 0
		glog.Warningf("libvirtxml: ignoring invalid vcpu value '%s'", str)
	}

	return uint(result)
}

func (s DomainVCPU) SetValue(value uint) {
	s.node.CharData = strconv.FormatUint(uint64(value), 10)
}
