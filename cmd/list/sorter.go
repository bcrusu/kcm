package list

import (
	"sort"
	"strings"

	"github.com/bcrusu/kcm/repository"
)

type sorter struct {
	Clusters []*repository.Cluster
}

func Sort(nodes []*repository.Cluster) {
	sorter := &sorter{nodes}
	sort.Sort(sorter)
}

func (s *sorter) Len() int {
	return len(s.Clusters)
}

func (s *sorter) Less(i, j int) bool {
	return strings.Compare(s.Clusters[i].Name, s.Clusters[j].Name) < 0
}

func (s *sorter) Swap(i, j int) {
	s.Clusters[i], s.Clusters[j] = s.Clusters[j], s.Clusters[i]
}
