package manifests

import (
	"path"

	"github.com/bcrusu/kcm/util"
)

type Params struct {
	ClusterName         string
	PodsNetworkCIDR     string
	ServicesNetworkCIDR string

	APIServerImageTag         string
	ControllerManagerImageTag string
	SchedulerImageTag         string
	ProxyImageTag             string
}

func WriteManifests(outDir string, params Params) error {
	if err := util.CreateDirectoryPath(outDir); err != nil {
		return err
	}

	if err := util.WriteFile(path.Join(outDir, "kube-apiserver.yaml"),
		util.GenerateTextTemplate(apiServerTemplate, apiServerTemplateParams{
			ImageTag:            params.APIServerImageTag,
			ServicesNetworkCIDR: params.ServicesNetworkCIDR,
		})); err != nil {
		return err
	}

	if err := util.WriteFile(path.Join(outDir, "kube-controller-manager.yaml"),
		util.GenerateTextTemplate(controllerManagerTemplate, controllerManagerTemplateParams{
			ImageTag:        params.ControllerManagerImageTag,
			ClusterName:     params.ClusterName,
			PodsNetworkCIDR: params.PodsNetworkCIDR,
		})); err != nil {
		return err
	}

	if err := util.WriteFile(path.Join(outDir, "kube-scheduler.yaml"),
		util.GenerateTextTemplate(schedulerTemplate, schedulerTemplateParams{
			ImageTag: params.SchedulerImageTag,
		})); err != nil {
		return err
	}

	if err := util.WriteFile(path.Join(outDir, "kube-proxy.yaml"),
		util.GenerateTextTemplate(proxyTemplate, proxyTemplateParams{
			ImageTag:        params.ProxyImageTag,
			PodsNetworkCIDR: params.PodsNetworkCIDR,
		})); err != nil {
		return err
	}

	return nil
}
