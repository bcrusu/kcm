package config

import (
	"io/ioutil"
	"path"

	"github.com/bcrusu/kcm/config/kubeconfig"
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
	if err := util.CreateDirectoryPath(path.Join(outDir, "manifests")); err != nil {
		return err
	}

	if err := c.writeCertificates(path.Join(outDir, "certs"), node); err != nil {
		return err
	}

	if err := kubeconfig.WriteKubeconfig(path.Join(outDir, "kubeconfig"), node, c.cluster); err != nil {
		return err
	}

	return nil
}

func (c ClusterConfig) stageKubernetesForCluster(outDir string) error {
	params := &manifests.Params{
		ClusterName:         c.cluster.Name,
		PodsNetworkCIDR:     c.podsNetworkCIDR,
		ServicesNetworkCIDR: c.servicesNetworkCIDR,
		FlannelImageTag:     "v0.7.1",
	}

	if err := c.readKubernetesImageTags(params); err != nil {
		return err
	}

	if err := manifests.WriteManifests(path.Join(outDir, "manifests"), *params); err != nil {
		return err
	}

	return nil
}

func (c ClusterConfig) readKubernetesImageTags(params *manifests.Params) error {
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

	if err := util.WriteFile(path.Join(outDir, "tls.pem"), node.Certificate); err != nil {
		return err
	}

	if err := util.WriteFile(path.Join(outDir, "tls-key.pem"), node.PrivateKey); err != nil {
		return err
	}

	return nil
}
