package addons

import (
	"path"

	"github.com/bcrusu/kcm/util"
)

type Params struct {
	PodsNetworkCIDR string
	ProxyImageTag   string
	FlannelImageTag string
}

func Write(outDir string, params Params) error {
	if err := util.CreateDirectoryPath(outDir); err != nil {
		return err
	}

	if err := util.WriteFile(path.Join(outDir, "kube-proxy.yaml"),
		util.GenerateTextTemplate(proxyTemplate, proxyTemplateParams{
			ImageTag:        params.ProxyImageTag,
			PodsNetworkCIDR: params.PodsNetworkCIDR,
		})); err != nil {
		return err
	}

	if err := util.WriteFile(path.Join(outDir, "flannel.yaml"),
		util.GenerateTextTemplate(flannelTemplate, flannelTemplateParams{
			ImageTag:        params.FlannelImageTag,
			PodsNetworkCIDR: params.PodsNetworkCIDR,
		})); err != nil {
		return err
	}

	return nil
}
