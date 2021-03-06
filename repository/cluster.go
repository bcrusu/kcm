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
	CNIVersion           string          `json:"cniVersion"`
	CoreOSVersion        string          `json:"coreOSVersion"`
	CoreOSChannel        string          `json:"coreOSChannel"`
	Nodes                map[string]Node `json:"nodes"` //map[NODE_NAME]NODE
	Network              Network         `json:"network"`
	StoragePool          string          `json:"storagePool"`
	BackingStorageVolume string          `json:"backingStorageVolume"`
	CACertificate        []byte          `json:"caCertificate"`
	CAPrivateKey         []byte          `json:"caPrivateKey"`
	DNSDomain            string          `json:"dnsDomain"`
	ServerURL            string          `json:"ServerUrl"`
}

type Node struct {
	Name                 string `json:"name"`
	IsMaster             bool   `json:"isMaster"`
	Domain               string `json:"domain"`
	MemoryMiB            uint   `json:"memory"`
	CPUs                 uint   `json:"cpus"`
	VolumeCapacityGiB    uint   `json:"volumeCapacity"`
	StoragePool          string `json:"storagePool"`
	BackingStorageVolume string `json:"backingStorageVolume"`
	StorageVolume        string `json:"storageVolume"`
	DNSName              string `json:"dnsName"`
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
	if c.Name == "" {
		return errors.New("missing cluster name")
	}

	if err := util.IsDNS1123Label(c.Name); err != nil {
		return err
	}

	if len(strings.TrimSpace(c.Name)) != len(c.Name) {
		return errors.New("invalid cluster name - cannot start/end with whitespaces")
	}

	if c.KubernetesVersion == "" {
		return errors.New("missing Kubernetes version")
	}

	if c.CNIVersion == "" {
		return errors.New("missing CNI version")
	}

	if c.CoreOSChannel == "" || c.CoreOSVersion == "" {
		return errors.New("invalid CoreOS version/channel")
	}

	if len(c.Nodes) < 1 {
		return errors.New("no node configured")
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
		return errors.New("no master node configured")
	}

	if mastersCount != 1 {
		return errors.New("multiple master clusters are not supported atm")
	}

	if c.StoragePool == "" {
		return errors.New("missing storage pool name")
	}

	if c.BackingStorageVolume == "" {
		return errors.New("missing backing storage volume")
	}

	if c.ServerURL == "" {
		return errors.New("missing server URL")
	}

	if err := c.Network.validate(); err != nil {
		return err
	}

	if c.CACertificate == nil {
		return errors.New("missing CA certificate")
	}

	if c.CAPrivateKey == nil {
		return errors.New("missing CA private key")
	}

	if c.DNSDomain == "" {
		return errors.New("missing cluster DNS domain name")
	}

	return nil
}

func (n *Node) Validate() error {
	if n == nil {
		return errors.Errorf("nil node")
	}

	if n.Name == "" {
		return errors.Errorf("missing node name")
	}

	if err := util.IsDNS1123Label(n.Name); err != nil {
		return err
	}

	if n.Domain == "" {
		return errors.Errorf("missing node domain")
	}

	if n.StorageVolume == "" {
		return errors.Errorf("missing node storage volume")
	}

	if n.StoragePool == "" {
		return errors.Errorf("missing node storage pool")
	}

	if n.BackingStorageVolume == "" {
		return errors.Errorf("missing backing storage volume")
	}

	if n.CPUs < 1 {
		return errors.Errorf("invalid CPUs value")
	}

	if n.MemoryMiB < 128 {
		return errors.Errorf("invalid memory value")
	}

	if n.VolumeCapacityGiB < 2 {
		return errors.Errorf("invalid volume capacity value")
	}

	if n.DNSName == "" {
		return errors.New("missing node DNS name")
	}

	return nil
}

func (n *Network) validate() error {
	if n == nil {
		return errors.Errorf("nil network")
	}

	if n.Name == "" {
		return errors.Errorf("missing network name")
	}

	if n.IPv4CIDR == "" {
		return errors.Errorf("missing network CIDR")
	}

	return nil
}
