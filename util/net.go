package util

import "net"

func GetBridgeIP(net *net.IPNet) net.IP {
	result := make([]byte, len(net.IP))
	copy(result, net.IP)

	result[len(result)-1]++
	return result
}

func GetMasterIP(net *net.IPNet) net.IP {
	result := make([]byte, len(net.IP))
	copy(result, net.IP)

	result[len(result)-1] += 2
	return result
}

func GetDHCPRange(ipnet *net.IPNet) (net.IP, net.IP) {
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
