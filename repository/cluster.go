package repository

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/bcrusu/kcm/util"
	"github.com/pkg/errors"
)

type Cluster struct {
	Name                 string          `json:"name"`
	KubernetesVersion    string          `json:"kubernetesVersion"`
	CoreOSVersion        string          `json:"coreOSVersion"`
	CoreOSChannel        string          `json:"coreOSChannel"`
	Nodes                map[string]Node `json:"nodes"` //map[NODE_NAME]NODE
	Network              Network         `json:"network"`
	StoragePool          string          `json:"storagePool"`
	BackingStorageVolume string          `json:"backingStorageVolume"`
	CACertificate        []byte          `json:"caCertificate"`
	CAPrivateKey         []byte          `json:"caPrivateKey"`
	DNSDomain            string          `json:"dnsDomain"`
	ServerURL            string          `json:"serverUrl"`
}

type Node struct {
	Name                 string `json:"name"`
	IsMaster             bool   `json:"isMaster"`
	Domain               string `json:"domain"`
	MemoryMiB            uint   `json:"memory"`
	CPUs                 uint   `json:"cpus"`
	StoragePool          string `json:"storagePool"`
	BackingStorageVolume string `json:"backingStorageVolume"`
	StorageVolume        string `json:"storageVolume"`
}

type Network struct {
	Name     string `json:"name"`
	IPv4CIDR string `json:"ipv4cidr"`
}

func loadCluster(clusterFile string) (*Cluster, error) {
	bytes, err := ioutil.ReadFile(clusterFile)
	if err != nil {
		return nil, errors.Wrapf(err, "repository: failed to read cluster '%s'", clusterFile)
	}

	cluster := &Cluster{}
	if err := json.Unmarshal(bytes, cluster); err != nil {
		return nil, errors.Wrapf(err, "repository: failed to unmarshall cluster '%s'", clusterFile)
	}

	if err := cluster.Validate(); err != nil {
		return nil, errors.Wrapf(err, "repository: validation failed for cluster '%s'", clusterFile)
	}

	return cluster, nil
}

func (c *Cluster) save(clusterFile string) error {
	if err := c.Validate(); err != nil {
		return errors.Wrapf(err, "repository: cluster '%s' validation failed", c.Name)
	}

	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return errors.Wrapf(err, "repository: failed to marshall cluster '%s'", c.Name)
	}

	if err := util.WriteFile(clusterFile, bytes); err != nil {
		return errors.Wrapf(err, "repository: failed to write cluster '%s'", clusterFile)
	}

	return nil
}

func (c *Cluster) Validate() error {
	//TODO: name must be a valid dns host name
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

	if len(c.Nodes) < 1 {
		return errors.New("repository: no node configured")
	}

	mastersCount := 0
	for _, node := range c.Nodes {
		if err := node.Validate(); err != nil {
			return err
		}

		if node.IsMaster {
			mastersCount++
		}
	}

	if mastersCount == 0 {
		return errors.New("repository: no master node configured")
	}

	if mastersCount != 1 {
		return errors.New("repository: multiple master clusters are not supported atm")
	}

	if c.StoragePool == "" {
		return errors.New("repository: missing storage pool name")
	}

	if c.BackingStorageVolume == "" {
		return errors.New("repository: missing backing storage volume")
	}

	if c.ServerURL == "" {
		return errors.New("repository: missing server URL")
	}

	if err := c.Network.validate(); err != nil {
		return err
	}

	if c.CACertificate == nil {
		return errors.New("repository: missing CA certificate")
	}

	if c.CAPrivateKey == nil {
		return errors.New("repository: missing CA private key")
	}

	if c.DNSDomain == "" {
		return errors.New("repository: missing cluster DNS domain name")
	}

	return nil
}

func (n *Node) Validate() error {
	//TODO: name must be a valid dns host name
	if n == nil {
		return errors.Errorf("repository: nil node")
	}

	if n.Name == "" {
		return errors.Errorf("repository: missing node name")
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

	if n.BackingStorageVolume == "" {
		return errors.Errorf("repository: missing backing storage volume")
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

	if n.IPv4CIDR == "" {
		return errors.Errorf("repository: missing network CIDR")
	}

	return nil
}
