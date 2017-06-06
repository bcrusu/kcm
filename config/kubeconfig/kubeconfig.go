package kubeconfig

import (
	"github.com/bcrusu/kcm/repository"
	"github.com/bcrusu/kcm/util"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
)

func WriteKubeconfig(filename string, node repository.Node, cluster repository.Cluster) error {
	config := &KubectlConfig{
		ApiVersion: "v1",
		Kind:       "Config",
		Users: []*KubectlUserWithName{
			{
				Name: "kubelet",
				User: KubectlUser{
					ClientCertificateData: node.Certificate,
					ClientKeyData:         node.PrivateKey,
				},
			},
		},
		Clusters: []*KubectlClusterWithName{
			{
				Name: "local",
				Cluster: KubectlCluster{
					CertificateAuthorityData: cluster.CACertificate,
					Server: cluster.ServerURL,
				},
			},
		},
		Contexts: []*KubectlContextWithName{
			{
				Name: "service-account-context",
				Context: KubectlContext{
					Cluster: "local",
					User:    "kubelet",
				},
			},
		},
		CurrentContext: "service-account-context",
	}

	bytes, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal kubeconfig to yaml")
	}

	return util.WriteFile(filename, bytes)
}
