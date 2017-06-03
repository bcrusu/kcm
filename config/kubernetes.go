package config

import (
	"io/ioutil"
	"path"

	"github.com/bcrusu/kcm/config/metadata"
	"github.com/bcrusu/kcm/repository"
	"github.com/bcrusu/kcm/util"
)

func (c ClusterConfig) stageKubernetesForNode(outDir string, node repository.Node) error {
	if err := util.CreateDirectoryPath(outDir); err != nil {
		return err
	}

	// create mount point for kubernetes bin
	if err := util.CreateDirectoryPath(path.Join(outDir, "bin")); err != nil {
		return err
	}

	// create mount point for static pods manifests
	if err := util.CreateDirectoryPath(path.Join(outDir, "metadata")); err != nil {
		return err
	}

	if err := c.writeCertificates(path.Join(outDir, "certs"), node); err != nil {
		return err
	}

	return nil
}

func (c ClusterConfig) stageKubernetesForCluster(outDir string) error {

	params := &metadata.MetadataParams{
		ClusterName:         c.cluster.Name,
		PodsNetworkCIDR:     c.podsNetworkCIDR,
		ServicesNetworkCIDR: c.servicesNetworkCIDR,
		FlannelImageTag:     "v0.7.1",
	}

	if err := c.readKubernetesImageTags(params); err != nil {
		return err
	}

	if err := metadata.WriteMetadataFiles(path.Join(outDir, "metadata"), *params); err != nil {
		return err
	}

	return nil
}

func (c ClusterConfig) readKubernetesImageTags(params *metadata.MetadataParams) error {
	var err error
	readTag := func(fileName string) (string, error) {
		bytes, err := ioutil.ReadFile(path.Join(c.kubernetesBinDir, fileName))
		if err != nil {
			return "", err
		}

		return string(bytes), nil
	}

	if params.APIServerImageTag, err = readTag("kube-apiserver.docker_tag"); err != nil {
		return err
	}

	if params.ControllerManagerImageTag, err = readTag("kube-controller-manager.docker_tag"); err != nil {
		return err
	}

	if params.ProxyImageTag, err = readTag("kube-proxy.docker_tag"); err != nil {
		return err
	}

	if params.SchedulerImageTag, err = readTag("kube-scheduler.docker_tag"); err != nil {
		return err
	}

	return nil
}

func (c ClusterConfig) writeCertificates(outDir string, node repository.Node) error {
	if err := util.CreateDirectoryPath(outDir); err != nil {
		return err
	}

	if err := util.WriteFile(path.Join(outDir, "ca.pem"), c.cluster.CACertificate); err != nil {
		return err
	}

	nodeDNSName := c.nodeDNSName(node.Name)

	hosts := []string{nodeDNSName}
	if node.IsMaster {
		hosts = append(hosts, "kubernetes", "kubernetes.default", "kubernetes.default.svc", "kubernetes.default.svc."+c.cluster.DNSDomain)
	}

	tlsCert, tlsKey, err := util.CreateCertificate(nodeDNSName, c.caCertificate, hosts...)
	if err != nil {
		return err
	}

	if err := util.WriteFile(path.Join(outDir, "tls.pem"), tlsCert); err != nil {
		return err
	}

	if err := util.WriteFile(path.Join(outDir, "tls-key.pem"), tlsKey); err != nil {
		return err
	}

	return nil
}
