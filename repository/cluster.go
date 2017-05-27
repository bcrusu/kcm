package repository

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"

	"github.com/pkg/errors"
)

const clusterSpecFileName = "cluster.json"

type Cluster struct {
	Name              string `json:"name"`
	KubernetesVersion string `json:"kubernetesVersion"`
	CoreOSVersion     string `json:"coreOSVersion"`
	CoreOSChannel     string `json:"coreOSChannel"`
	Nodes             []Node `json:"nodes"`
	StoragePool       string `json:"storagePool"`
}

type Node struct {
	DomainName string `json:"domainName"`
	IsMaster   bool   `json:"isMaster"`
}

func loadCluster(clusterDir string) (*Cluster, error) {
	clusterFile := path.Join(clusterDir, clusterSpecFileName)
	bytes, err := ioutil.ReadFile(clusterFile)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read cluster '%s'", clusterFile)
	}

	cluster := &Cluster{}
	if err := json.Unmarshal(bytes, cluster); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshall cluster '%s'", clusterFile)
	}

	return cluster, nil
}

func (c *Cluster) save(clusterDir string) error {
	bytes, err := json.Marshal(c)
	if err != nil {
		return errors.Wrapf(err, "failed to marshall cluster '%s'", c.Name)
	}

	clusterFile := path.Join(clusterDir, clusterSpecFileName)
	if err := ioutil.WriteFile(clusterFile, bytes, 0644); err != nil {
		return errors.Wrapf(err, "failed to write cluster '%s'", clusterFile)
	}

	return nil
}

func (c *Cluster) validate() error {
	if c.Name == "" {
		return errors.New("missing cluster name")
	}

	if len(strings.TrimSpace(c.Name)) != len(c.Name) {
		return errors.New("invalid cluster name - cannot start/end with whitespaces")
	}

	if c.KubernetesVersion == "" {
		return errors.New("missing Kubernetes version")
	}

	if c.CoreOSChannel == "" || c.CoreOSVersion == "" {
		return errors.New("invalid CoreOS version/channel")
	}

	masterCount := 0
	nodeCount := 0
	for i, node := range c.Nodes {
		if node.DomainName == "" {
			return errors.Errorf("missing domain name for node %d", i)
		}

		if node.IsMaster {
			masterCount++
		} else {
			nodeCount++
		}
	}

	if masterCount < 1 {
		return errors.New("no master node configured")
	}

	if nodeCount < 1 {
		return errors.New("no worker node configured")
	}

	if c.StoragePool == "" {
		return errors.New("missing storage pool name")
	}

	return nil
}
