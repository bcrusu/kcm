package coreos

import (
	"bytes"
	"path"
	"text/template"

	"github.com/bcrusu/kcm/util"
)

func WriteCoreOSConfig(outDir string, params CloudConfigParams) error {
	userDataDir := path.Join(outDir, "openstack", "latest")
	if err := util.CreateDirectoryPath(userDataDir); err != nil {
		return err
	}

	data := generateCoreOSConfig(params)

	userDataFilename := path.Join(userDataDir, "user_data")
	if err := util.WriteFile(userDataFilename, data); err != nil {
		return err
	}

	return nil
}

func generateCoreOSConfig(params CloudConfigParams) []byte {
	t := template.New("coreos")

	t.Funcs(template.FuncMap{
		"Role": func() string {
			if params.IsMaster {
				return "master"
			}

			return "node"
		},
	})

	if _, err := t.Parse(cloudConfigTemplate); err != nil {
		panic(err)
	}

	buffer := &bytes.Buffer{}
	if err := t.ExecuteTemplate(buffer, "coreos", params); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}
