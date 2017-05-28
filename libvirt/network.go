package libvirt

import (
	"net"

	"github.com/bcrusu/kcm/libvirtxml"
	"github.com/golang/glog"
	"github.com/libvirt/libvirt-go"
	"github.com/pkg/errors"
)

func lookupNetwork(connect *libvirt.Connect, lookup string) (*libvirt.Network, error) {
	if len(lookup) == uuidStringLength {
		net, err := connect.LookupNetworkByUUIDString(lookup)
		if err != nil {
			if lverr, ok := err.(libvirt.Error); ok && lverr.Code != libvirt.ERR_NO_NETWORK {
				glog.Infof("libvirt: network lookup by ID '%s' failed. Error: %v", lookup, lverr)
			}
		}

		if net != nil {
			return net, nil
		}
	}

	net, err := connect.LookupNetworkByName(lookup)
	if err != nil {
		if lverr, ok := err.(libvirt.Error); ok && lverr.Code == libvirt.ERR_NO_NETWORK {
			return nil, nil
		}

		return nil, errors.Wrapf(err, "libvirt: network lookup failed '%s'", lookup)
	}

	return net, nil
}

func getNetworkXML(network *libvirt.Network) (*libvirtxml.Network, error) {
	xml, err := network.GetXMLDesc(libvirt.NetworkXMLFlags(0))
	if err != nil {
		return nil, errors.Wrapf(err, "libvirt: failed to fetch network XML description")
	}

	return libvirtxml.NewNetworkForXML(xml)
}

func defineNATNetwork(connect *libvirt.Connect, name string, ipv4CIDR, ipv6CIDR string) error {
	networkXML := libvirtxml.NewNetwork()
	networkXML.SetName(name)
	networkXML.Forward().SetMode("nat")
	networkXML.Forward().SetNATPortRange(1024, 65535)

	networkXML.Bridge().SetSTP(true)

	if ipv4CIDR != "" {
		addIP(networkXML, ipv4CIDR)
	}

	if ipv6CIDR != "" {
		addIP(networkXML, ipv6CIDR)
	}

	if len(networkXML.IPs()) == 0 {
		return errors.New("libvirt: failed to define network - missing CIDR")
	}

	xmlString, err := networkXML.MarshalToXML()
	if err != nil {
		return err
	}

	network, err := connect.NetworkDefineXML(xmlString)
	if err != nil {
		return errors.Wrapf(err, "libvirt: failed to define network '%s'", name)
	}
	defer network.Free()

	return err
}

func addIP(network libvirtxml.Network, cidr string) error {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return errors.Wrapf(err, "libvirt: failed to define network - invalid CIDR '%s'", cidr)
	}

	var family string
	switch len(ipnet.IP) {
	case net.IPv4len:
		family = "ipv4"
	case net.IPv6len:
		family = "ipv6"
	default:
		return errors.Wrapf(err, "libvirt: failed to define network - invalid CIDR IP '%s'", ip)
	}

	prefix, bits := ipnet.Mask.Size()
	if bits-prefix < 3 {
		return errors.Wrapf(err, "libvirt: failed to define network - network too small '%s'", cidr)
	}

	ipXML := network.NewIP()
	ipXML.SetFamily(family)
	ipXML.SetAddress(getBridgeIP(ipnet).String())
	ipXML.SetPrefix(prefix)

	dhcpStart, dhcpEnd := getDHCPRange(ipnet)
	ipXML.SetDHCPRange(dhcpStart.String(), dhcpEnd.String())

	return nil
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
