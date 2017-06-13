package addons

import (
	"path"

	"github.com/bcrusu/kcm/util"
)

type Params struct {
	PodsNetworkCIDR string
}

func WriteManifests(outDir string, params Params) error {
	if err := util.CreateDirectoryPath(outDir); err != nil {
		return err
	}

	if err := util.WriteFile(path.Join(outDir, "flannel.yaml"),
		util.GenerateTextTemplate(flannelTemplate, flannelTemplateParams{
			ImageTag:        "v0.7.1",
			PodsNetworkCIDR: params.PodsNetworkCIDR,
		})); err != nil {
		return err
	}

	return nil
}
