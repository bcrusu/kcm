package repository

type ClusterSpec struct {
	Name              string     `json:"name"`
	KubernetesVersion string     `json:"kubernetesVersion"`
	CoreOSVersion     string     `json:"coreOSVersion"`
	CoreOSChannel     string     `json:"coreOSChannel"`
	MasterCount       int        `json:"masterCount"`
	NodeCount         int        `json:"nodeCount"`
	Nodes             []NodeSpec `json:"nodes"`
	StoragePool       string     `json:"storagePool"`
}

type NodeSpec struct {
	DomainName string `json:"domainName"`
	IsMaster   bool   `json:"isMaster"`
}
