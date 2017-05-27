package create

import (
	"compress/bzip2"
	"fmt"
	"io"
	"path"

	"github.com/bcrusu/kcm/util"
	"github.com/golang/glog"
)

func CoreOSVolumeName(version string) string {
	return fmt.Sprintf("coreos_production_qemu_image_%s.qcow2", version)
}

func DownloadCoreOS(channel, version string, outDir string) error {
	fileName := CoreOSVolumeName(version)
	filePath := path.Join(outDir, fileName)

	url := fmt.Sprintf("https://%s.release.core-os.net/amd64-usr/%s/coreos_production_qemu_image.img.bz2", channel, version)

	glog.Infof("downloading CoreOS from '%s' to path '%s'", url, filePath)

	downloadReader, err := util.DownloadHTTP(url)
	if err != nil {
		return err
	}
	defer util.CloseNoError(downloadReader)

	decompressReader := bzip2.NewReader(downloadReader)

	file, err := util.CreateFile(filePath, 0644)
	if err != nil {
		return err
	}
	defer util.CloseNoError(file)

	_, err = io.Copy(file, decompressReader)
	if err != nil {
		return err
	}

	return nil
}
