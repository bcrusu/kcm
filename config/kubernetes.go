package config

import (
	"io/ioutil"
	"path"

	"github.com/bcrusu/kcm/repository"
	"github.com/bcrusu/kcm/util"
)

type imageTags struct {
	APIServer         string
	ControllerManager string
	Scheduler         string
	Proxy             string
}

func (c ClusterConfig) stageKubernetes(outDir string, node repository.Node) error {
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

	if err := c.writCertificates(path.Join(outDir, "certs"), node); err != nil {
		return err
	}

	//NonMasqueradeCIDR = "100.64.0.0/10"

	return nil
}

func (c ClusterConfig) loadImageTags() (*imageTags, error) {
	result := &imageTags{}
	var err error

	readTag := func(fileName string) (string, error) {
		bytes, err := ioutil.ReadFile(path.Join(c.kubernetesBinDir, fileName))
		if err != nil {
			return "", err
		}

		return string(bytes), nil
	}

	if result.APIServer, err = readTag("kube-apiserver.docker_tag"); err != nil {
		return nil, err
	}

	if result.ControllerManager, err = readTag("kube-controller-manager.docker_tag"); err != nil {
		return nil, err
	}

	if result.Proxy, err = readTag("kube-proxy.docker_tag"); err != nil {
		return nil, err
	}

	if result.Scheduler, err = readTag("kube-scheduler.docker_tag"); err != nil {
		return nil, err
	}

	return result, nil
}

func (c ClusterConfig) writCertificates(outDir string, node repository.Node) error {
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
