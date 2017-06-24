package status

import (
	"sort"
	"strings"
)

type nodeStatusSorter struct {
	Nodes []NodeStatus
}

func Sort(nodes []NodeStatus) {
	sorter := &nodeStatusSorter{nodes}
	sort.Sort(sorter)
}

func (s *nodeStatusSorter) Len() int {
	return len(s.Nodes)
}

func (s *nodeStatusSorter) Less(i, j int) bool {
	s1 := s.Nodes[i]
	s2 := s.Nodes[j]

	if s1.Node.IsMaster != s2.Node.IsMaster {
		// masters above minions (a hard m' fact of life)
		return s1.Node.IsMaster
	}

	return strings.Compare(s1.Node.Name, s2.Node.Name) < 0
}

func (s *nodeStatusSorter) Swap(i, j int) {
	s.Nodes[i], s.Nodes[j] = s.Nodes[j], s.Nodes[i]
}
