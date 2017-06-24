package util

import (
	"net"
	"regexp"

	"github.com/pkg/errors"
)

type NetworkInfo struct {
	Family         string // ipv4 or ipv6
	BridgeIP       net.IP
	DHCPRangeStart net.IP
	DHCPRangeEnd   net.IP
	Net            *net.IPNet
}

func ParseNetworkCIDR(cidr string) (*NetworkInfo, error) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, errors.Wrapf(err, "net: failed to parse CIDR '%s'", cidr)
	}

	var family string
	switch len(ipnet.IP) {
	case net.IPv4len:
		family = "ipv4"
	case net.IPv6len:
		family = "ipv6"
	default:
		return nil, errors.Wrapf(err, "net: failed to parse CIDR '%s'", cidr)
	}

	dhcpStart, dhcpEnd := getDHCPRange(ipnet)

	return &NetworkInfo{
		Family:         family,
		BridgeIP:       getBridgeIP(ipnet),
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

func getDHCPRange(ipnet *net.IPNet) (net.IP, net.IP) {
	ipLen := len(ipnet.IP)

	start := make([]byte, ipLen)
	{
		copy(start, ipnet.IP)
		start[ipLen-1] += 2 // first IP is assigned to the bridge
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

const dns1123LabelFmt string = "[a-z0-9]([-a-z0-9]*[a-z0-9])?"
const dns1123LabelErrMsg string = "a DNS-1123 label must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character"
const DNS1123LabelMaxLength int = 63

var dns1123LabelRegexp = regexp.MustCompile("^" + dns1123LabelFmt + "$")

// stolen from here: https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/util/validation/validation.go
// IsDNS1123Label tests for a string that conforms to the definition of a label in DNS (RFC 1123).
func IsDNS1123Label(value string) error {
	if len(value) > DNS1123LabelMaxLength {
		return errors.Errorf("net: invalid DNS label '%s' - max length of 63 chars exceeded", value)
	}

	if !dns1123LabelRegexp.MatchString(value) {
		return errors.Errorf("net: invalid DNS label '%s' - "+dns1123LabelErrMsg, value)
	}

	return nil
}
