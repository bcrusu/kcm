package util

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

type NetworkInfo struct {
	Family         string // ipv4 or ipv6
	BridgeIP       net.IP
	MasterIP       net.IP
	MasterAddress  string // e.g. 10.1.0.2/16
	DHCPRangeStart net.IP
	DHCPRangeEnd   net.IP
	Net            *net.IPNet
}

func ParseNetworkCIDR(cidr string) (*NetworkInfo, error) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse CIDR '%s'", cidr)
	}

	var family string
	switch len(ipnet.IP) {
	case net.IPv4len:
		family = "ipv4"
	case net.IPv6len:
		family = "ipv6"
	default:
		return nil, errors.Wrapf(err, "failed to parse CIDR '%s'", cidr)
	}

	dhcpStart, dhcpEnd := getDHCPRange(ipnet)

	return &NetworkInfo{
		Family:         family,
		BridgeIP:       getBridgeIP(ipnet),
		MasterIP:       getMasterIP(ipnet),
		MasterAddress:  getMasterAddress(ipnet),
		DHCPRangeStart: dhcpStart,
		DHCPRangeEnd:   dhcpEnd,
		Net:            ipnet,
	}, nil
}

func getBridgeIP(net *net.IPNet) net.IP {
	result := make([]byte, len(net.IP))
	copy(result, net.IP)

	result[len(result)-1]++
	return result
}

func getMasterIP(net *net.IPNet) net.IP {
	result := make([]byte, len(net.IP))
	copy(result, net.IP)

	result[len(result)-1] += 2
	return result
}

func getMasterAddress(net *net.IPNet) string {
	ip := getMasterIP(net)
	size, _ := net.Mask.Size()

	return fmt.Sprintf("%s/%d", ip, size)
}

func getDHCPRange(ipnet *net.IPNet) (net.IP, net.IP) {
	ipLen := len(ipnet.IP)

	start := make([]byte, ipLen)
	{
		copy(start, ipnet.IP)
		start[ipLen-1] += 3 // first IP is assigned to the bridge and 2nd to the master/load balancer
	}

	end := make([]byte, ipLen)
	{
		copy(end, ipnet.IP)

		for i, b := range ipnet.Mask {
			end[i] += ^b
		}

		if ipLen == net.IPv4len {
			end[ipLen-1]-- // exclude broadcast address
		}
	}

	return start, end
}
