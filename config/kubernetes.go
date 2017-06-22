package config

import (
	"io/ioutil"
	"path"

	"github.com/bcrusu/kcm/config/addons"
	"github.com/bcrusu/kcm/config/kubeconfig"
	"github.com/bcrusu/kcm/config/manifests"
	"github.com/bcrusu/kcm/repository"
	"github.com/bcrusu/kcm/util"
)

type dockerImageTags struct {
	APIServer         string
	ControllerManager string
	Scheduler         string
	Proxy             string
	Flannel           string
	DNS               string
}

func (c ClusterConfig) stageKubernetesForNode(outDir string, node repository.Node) error {
	if err := util.CreateDirectoryPath(outDir); err != nil {
		return err
	}

	{
		// create mount points
		mountPoints := []string{"bin", "manifests", "addons", "kubeconfig"}
		for _, mountPoint := range mountPoints {
			if err := util.CreateDirectoryPath(path.Join(outDir, mountPoint)); err != nil {
				return err
			}
		}
	}

	if err := c.writeCertificates(path.Join(outDir, "certs"), node); err != nil {
		return err
	}

	return nil
}

func (c ClusterConfig) stageKubernetesForCluster(outDir string) error {
	imageTags, err := c.getDockerImageTags()
	if err != nil {
		return err
	}

	if err := manifests.Write(path.Join(outDir, "manifests"), manifests.Params{
		ClusterName:               c.cluster.Name,
		PodsNetworkCIDR:           c.podsNetworkCIDR,
		ServicesNetworkCIDR:       c.servicesNetworkCIDR,
		APIServerImageTag:         imageTags.APIServer,
		ControllerManagerImageTag: imageTags.ControllerManager,
		SchedulerImageTag:         imageTags.Scheduler,
	}); err != nil {
		return err
	}

	if err := addons.Write(path.Join(outDir, "addons"), addons.Params{
		ClusterDomain:   c.cluster.DNSDomain,
		DNSServiceIP:    c.dnsServiceIP,
		PodsNetworkCIDR: c.podsNetworkCIDR,
		ProxyImageTag:   imageTags.Proxy,
		FlannelImageTag: imageTags.Flannel,
		DNSImageTag:     imageTags.DNS,
	}); err != nil {
		return err
	}

	if err := kubeconfig.WriteKubeconfigFiles(path.Join(outDir, "kubeconfig"), c.cluster); err != nil {
		return err
	}

	return nil
}

func (c ClusterConfig) getDockerImageTags() (*dockerImageTags, error) {
	result := &dockerImageTags{
		Flannel: "v0.7.1",
		DNS:     "1.14.2",
	}

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

func (c ClusterConfig) writeCertificates(outDir string, node repository.Node) error {
	if err := util.CreateDirectoryPath(outDir); err != nil {
		return err
	}

	if err := util.WriteFile(path.Join(outDir, "ca.pem"), c.cluster.CACertificate); err != nil {
		return err
	}

	if err := util.WriteFile(path.Join(outDir, "ca-key.pem"), c.cluster.CAPrivateKey); err != nil {
		return err
	}

	caCert, err := util.ParseCertificate(c.cluster.CACertificate)
	if err != nil {
		return err
	}

	caKey, err := util.ParsePrivateKey(c.cluster.CAPrivateKey)
	if err != nil {
		return err
	}

	if node.IsMaster {
		clientCert, clientKey, err := util.CreateClientCertificate(node.DNSName, caCert, caKey)
		if err != nil {
			return err
		}

		if err := util.WriteFile(path.Join(outDir, "tls-client.pem"), clientCert); err != nil {
			return err
		}

		if err := util.WriteFile(path.Join(outDir, "tls-client-key.pem"), clientKey); err != nil {
			return err
		}
	}

	{
		hosts := []string{
			node.DNSName,
		}

		if node.IsMaster {
			hosts = append(hosts, []string{
				c.apiServerServiceIP,
				"kubernetes",
				"kubernetes.default",
				"kubernetes.default.svc",
				"kubernetes.default.svc." + c.cluster.DNSDomain,
			}...)
		}

		serverCert, serverKey, err := util.CreateServerCertificate(node.DNSName, caCert, caKey, hosts...)
		if err != nil {
			return err
		}

		if err := util.WriteFile(path.Join(outDir, "tls-server.pem"), serverCert); err != nil {
			return err
		}

		if err := util.WriteFile(path.Join(outDir, "tls-server-key.pem"), serverKey); err != nil {
			return err
		}
	}

	return nil
}
