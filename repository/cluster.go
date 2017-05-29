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
	Name                 string  `json:"name"`
	KubernetesVersion    string  `json:"kubernetesVersion"`
	CoreOSVersion        string  `json:"coreOSVersion"`
	CoreOSChannel        string  `json:"coreOSChannel"`
	Masters              []Node  `json:"masters"`
	Nodes                []Node  `json:"nodes"`
	Network              Network `json:"network"`
	StoragePool          string  `json:"storagePool"`
	BackingStorageVolume string  `json:"backingStorageVolume"`
}

type Node struct {
	Domain        string `json:"domain"`
	MemoryMiB     uint   `json:"memory"`
	CPUs          uint   `json:"cpus"`
	StoragePool   string `json:"storagePool"`
	StorageVolume string `json:"storageVolume"`
}

type Network struct {
	Name     string `json:"name"`
	IPv4CIDR string `json:"ipv4cidr"`
	IPv6CIDR string `json:"ipv6cidr"`
}

func loadCluster(clusterDir string) (*Cluster, error) {
	clusterFile := path.Join(clusterDir, clusterSpecFileName)
	bytes, err := ioutil.ReadFile(clusterFile)
	if err != nil {
		return nil, errors.Wrapf(err, "repository: failed to read cluster '%s'", clusterFile)
	}

	cluster := &Cluster{}
	if err := json.Unmarshal(bytes, cluster); err != nil {
		return nil, errors.Wrapf(err, "repository: failed to unmarshall cluster '%s'", clusterFile)
	}

	return cluster, nil
}

func (c *Cluster) save(clusterDir string) error {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return errors.Wrapf(err, "repository: failed to marshall cluster '%s'", c.Name)
	}

	clusterFile := path.Join(clusterDir, clusterSpecFileName)
	if err := ioutil.WriteFile(clusterFile, bytes, 0644); err != nil {
		return errors.Wrapf(err, "repository: failed to write cluster '%s'", clusterFile)
	}

	return nil
}

func (c *Cluster) Validate() error {
	if c.Name == "" {
		return errors.New("repository: missing cluster name")
	}

	if len(strings.TrimSpace(c.Name)) != len(c.Name) {
		return errors.New("repository: invalid cluster name - cannot start/end with whitespaces")
	}

	if c.KubernetesVersion == "" {
		return errors.New("repository: missing Kubernetes version")
	}

	if c.CoreOSChannel == "" || c.CoreOSVersion == "" {
		return errors.New("repository: invalid CoreOS version/channel")
	}

	if len(c.Masters) < 1 {
		return errors.New("repository: no master node configured")
	}

	for _, node := range c.Masters {
		if err := node.validate(); err != nil {
			return err
		}
	}

	for _, node := range c.Nodes {
		if err := node.validate(); err != nil {
			return err
		}
	}

	if c.StoragePool == "" {
		return errors.New("repository: missing storage pool name")
	}

	if c.BackingStorageVolume == "" {
		return errors.New("repository: missing backing storage volume")
	}

	if err := c.Network.validate(); err != nil {
		return err
	}

	return nil
}

func (c *Cluster) Node(domain string) (Node, bool) {
	for _, node := range c.Nodes {
		if node.Domain == domain {
			return node, true
		}
	}

	return Node{}, false
}

func (c *Cluster) RemoveNode(domain string) {
	var filtered []Node
	for _, node := range c.Nodes {
		if node.Domain != domain {
			filtered = append(filtered, node)
		}
	}

	c.Nodes = filtered
}

func (c *Cluster) Master(domain string) (Node, bool) {
	for _, node := range c.Masters {
		if node.Domain == domain {
			return node, true
		}
	}

	return Node{}, false
}

func (c *Cluster) RemoveMaster(domain string) {
	var filtered []Node
	for _, node := range c.Masters {
		if node.Domain != domain {
			filtered = append(filtered, node)
		}
	}

	c.Masters = filtered
}

func (n *Node) validate() error {
	if n == nil {
		return errors.Errorf("repository: nil node")
	}

	if n.Domain == "" {
		return errors.Errorf("repository: missing node domain")
	}

	if n.StorageVolume == "" {
		return errors.Errorf("repository: missing node storage volume")
	}

	if n.StoragePool == "" {
		return errors.Errorf("repository: missing node storage pool")
	}

	if n.CPUs < 1 {
		return errors.Errorf("repository: invalid CPUs value")
	}

	if n.MemoryMiB < 128 {
		return errors.Errorf("repository: invalid memory value")
	}

	return nil
}

func (n *Network) validate() error {
	if n == nil {
		return errors.Errorf("repository: nil network")
	}

	if n.Name == "" {
		return errors.Errorf("repository: missing network name")
	}

	if len(n.IPv4CIDR) == 0 && len(n.IPv6CIDR) == 0 {
		return errors.Errorf("repository: missing network CIDR")
	}

	return nil
}
