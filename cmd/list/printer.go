package list

import (
	"os"

	"github.com/bcrusu/kcm/repository"
	"github.com/bcrusu/kcm/util"
)

func Print(clusters []*repository.Cluster, current *repository.Cluster) {
	Sort(clusters)

	writer := util.NewTabWriter(os.Stdout)

	writer.Print("CURRENT\tCLUSTER")
	writer.Nl()

	for _, cluster := range clusters {
		mark := ""
		if current != nil && cluster.Name == current.Name {
			mark = "*"
		}

		writer.Print(mark)
		writer.Tab()
		writer.Print(cluster.Name)

		writer.Nl()
	}

	writer.Flush()
}
