package kubeconfig

import (
	"crypto/rsa"
	"crypto/x509"
	"path"

	"github.com/bcrusu/kcm/repository"
	"github.com/bcrusu/kcm/util"
	"github.com/ghodss/yaml"
)

func WriteKubeconfigFiles(outDir string, cluster repository.Cluster) error {
	if err := util.CreateDirectoryPath(outDir); err != nil {
		return err
	}

	caCert, err := util.ParseCertificate(cluster.CACertificate)
	if err != nil {
		return err
	}

	caKey, err := util.ParsePrivateKey(cluster.CAPrivateKey)
	if err != nil {
		return err
	}

	if err := generateKubeconfigFile(path.Join(outDir, "kubelet"), "kubelet", cluster, caCert, caKey); err != nil {
		return err
	}

	if err := generateKubeconfigFile(path.Join(outDir, "kube-scheduler"), KubeScheduler, cluster, caCert, caKey); err != nil {
		return err
	}

	if err := generateKubeconfigFile(path.Join(outDir, "kube-controller-manager"), KubeControllerManager, cluster, caCert, caKey); err != nil {
		return err
	}

	if err := generateKubeconfigFile(path.Join(outDir, "kube-proxy"), KubeProxy, cluster, caCert, caKey); err != nil {
		return err
	}

	if err := generateKubeconfigFile(path.Join(outDir, "kubectl"), "kubectl", cluster, caCert, caKey); err != nil {
		return err
	}

	return nil
}

func generateKubeconfigFile(filename string, user string, cluster repository.Cluster, caCert *x509.Certificate, caKey *rsa.PrivateKey) error {
	clientCert, clientKey, err := util.CreateClientCertificate(user, caCert, caKey)
	if err != nil {
		return err
	}

	bytes, err := CreateKubeconfig(user, cluster.Name, cluster.ServerURL, cluster.CACertificate, clientCert, clientKey)
	if err != nil {
		return err
	}

	return util.WriteFile(filename, bytes)
}

func CreateKubeconfig(user, clusterName, serverURL string, caCert, clientCert, clientKey []byte) ([]byte, error) {
	config := &KubectlConfig{
		ApiVersion: "v1",
		Kind:       "Config",
		Users: []*KubectlUserWithName{
			{
				Name: user,
				User: KubectlUser{
					ClientCertificateData: clientCert,
					ClientKeyData:         clientKey,
				},
			},
		},
		Clusters: []*KubectlClusterWithName{
			{
				Name: clusterName,
				Cluster: KubectlCluster{
					CertificateAuthorityData: caCert,
					Server: serverURL,
				},
			},
		},
		Contexts: []*KubectlContextWithName{
			{
				Name: "ctx",
				Context: KubectlContext{
					Cluster: clusterName,
					User:    user,
				},
			},
		},
		CurrentContext: "ctx",
	}

	return yaml.Marshal(config)
}
