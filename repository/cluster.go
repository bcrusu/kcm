package repository

type Cluster struct {
	spec ClusterSpec
}

func newCluster(path string) (*Cluster, error) {
	return nil, nil
}

func (c *Cluster) Save(path string) error {
	//TODO
	return nil
}

func (c *Cluster) Name() string {
	return c.spec.Name
}
