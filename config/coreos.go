package config

import (
	"bytes"
	"path"
	"text/template"

	"github.com/bcrusu/kcm/util"
)

type coreOSTemplateParams struct {
	Name          string
	IsMaster      bool
	SSHPublicKey  string
	MasterIP      string
	CoreOSChannel string
}

func writeCoreOSConfig(configDir string, params coreOSTemplateParams) error {
	userDataDir := path.Join(configDir, "config-2", "openstack", "latest")
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

func generateCoreOSConfig(params coreOSTemplateParams) []byte {
	t := template.New("coreos")
	if _, err := t.Parse(coreOSCloudConfigTemplate); err != nil {
		panic(err)
	}

	buffer := &bytes.Buffer{}
	if err := t.ExecuteTemplate(buffer, "coreos", params); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}
