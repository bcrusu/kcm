package manifests

import (
	"bytes"
	"html/template"
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
	FlannelImageTag           string
}

func WriteManifests(outDir string, params Params) error {
	if err := util.CreateDirectoryPath(outDir); err != nil {
		return err
	}

	if err := util.WriteFile(path.Join(outDir, "kube-apiserver.yaml"),
		generateTemplate(apiServerTemplate, apiServerTemplateParams{
			ImageTag:            params.APIServerImageTag,
			ServicesNetworkCIDR: params.ServicesNetworkCIDR,
		})); err != nil {
		return err
	}

	if err := util.WriteFile(path.Join(outDir, "kube-controller-manager.yaml"),
		generateTemplate(controllerManagerTemplate, controllerManagerTemplateParams{
			ImageTag:            params.ControllerManagerImageTag,
			ClusterName:         params.ClusterName,
			PodsNetworkCIDR:     params.PodsNetworkCIDR,
			ServicesNetworkCIDR: params.ServicesNetworkCIDR,
		})); err != nil {
		return err
	}

	if err := util.WriteFile(path.Join(outDir, "kube-scheduler.yaml"),
		generateTemplate(schedulerTemplate, schedulerTemplateParams{
			ImageTag: params.SchedulerImageTag,
		})); err != nil {
		return err
	}

	if err := util.WriteFile(path.Join(outDir, "kube-proxy.yaml"),
		generateTemplate(proxyTemplate, proxyTemplateParams{
			ImageTag:        params.ProxyImageTag,
			PodsNetworkCIDR: params.PodsNetworkCIDR,
		})); err != nil {
		return err
	}

	if err := util.WriteFile(path.Join(outDir, "flannel.yaml"),
		generateTemplate(flannelTemplate, flannelTemplateParams{
			ImageTag:        params.FlannelImageTag,
			PodsNetworkCIDR: params.PodsNetworkCIDR,
		})); err != nil {
		return err
	}

	return nil
}

func generateTemplate(templateStr string, params interface{}) []byte {
	t := template.New("t")

	if _, err := t.Parse(templateStr); err != nil {
		panic(err)
	}

	buffer := &bytes.Buffer{}
	if err := t.ExecuteTemplate(buffer, "t", params); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}
