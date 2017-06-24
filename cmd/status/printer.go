package status

import (
	"fmt"
	"os"
	"strings"

	"github.com/bcrusu/kcm/util"
)

func PrintNodes(stats []NodeStatus) {
	Sort(stats)

	writer := util.NewTabWriter(os.Stdout)

	writer.Print("NODE\tSTATUS\tDNS NAME\tIP")
	writer.Nl()

	for _, stat := range stats {
		writer.Print(stat.Node.Name)
		writer.Tab()

		if stat.Missing {
			writer.Print("Missing")
		} else if stat.Active {
			writer.Print("Active")
		} else {
			writer.Print("Inactive")
		}
		writer.Tab()

		writer.Print(stat.Node.DNSName)
		writer.Tab()

		if len(stat.Addresses) > 0 {
			writer.Print(strings.Join(stat.Addresses, ", "))
		}

		writer.Nl()
	}

	writer.Flush()
}

func PrintNetwork(stat NetworkStatus) {
	writer := util.NewTabWriter(os.Stdout)

	writer.Print("NETWORK\tSTATUS\tDNS SERVER")
	writer.Nl()

	writer.Print(stat.Network.Name)
	writer.Tab()

	if stat.Missing {
		writer.Print("Missing")
	} else if stat.Active {
		writer.Print("Active")
	} else {
		writer.Print("Inactive")
	}
	writer.Tab()

	{
		networkInfo, err := util.ParseNetworkCIDR(stat.Network.IPv4CIDR)
		if err != nil {
			panic("failed to parse network CIDR")
		}

		writer.Print(networkInfo.BridgeIP.String())
	}
	writer.Nl()

	writer.Flush()
}

func PrintCluster(stat ClusterStatus) {
	writer := util.NewTabWriter(os.Stdout)

	writer.Print("CLUSTER\tSTATUS\tDNS DOMAIN\tKUBE VERSION\tCOREOS VERSION")
	writer.Nl()

	writer.Print(stat.Cluster.Name)
	writer.Tab()

	if stat.Active {
		writer.Print("Active")
	} else {
		writer.Print("Inactive")
	}
	writer.Tab()

	writer.Print(stat.Cluster.DNSDomain)
	writer.Tab()

	writer.Print(stat.Cluster.KubernetesVersion)
	writer.Tab()

	writer.Print(fmt.Sprintf("%s/%s", stat.Cluster.CoreOSChannel, stat.Cluster.CoreOSVersion))
	writer.Nl()

	writer.Flush()
}
