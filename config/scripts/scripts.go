package scripts

import (
	"path"

	"github.com/bcrusu/kcm/util"
)

func Write(outDir string) error {
	if err := util.CreateDirectoryPath(outDir); err != nil {
		return err
	}

	if err := util.WriteExecutableFile(path.Join(outDir, "install_socat"), []byte(installSocat)); err != nil {
		return err
	}

	if err := util.WriteExecutableFile(path.Join(outDir, "load_addons"), []byte(loadAddons)); err != nil {
		return err
	}

	return nil
}
