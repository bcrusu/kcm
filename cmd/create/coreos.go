package create

import (
	"compress/bzip2"
	"fmt"
	"path"

	"github.com/bcrusu/kcm/util"
	"github.com/golang/glog"
)

func coreOSVolumeName(version string) string {
	return fmt.Sprintf("coreos_production_qemu_image_%s.qcow2", version)
}

func downloadCoreOSImage(channel, version string, outDir string) error {
	fileName := coreOSVolumeName(version)
	filePath := path.Join(outDir, fileName)

	if util.FileExists(filePath) {
		glog.Infof("skipped CoreOS image download. File '%s' exists.", filePath)
		return nil
	}

	url := fmt.Sprintf("https://%s.release.core-os.net/amd64-usr/%s/coreos_production_qemu_image.img.bz2", channel, version)

	glog.Infof("downloading '%s' to path '%s'", url, filePath)

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

	err = util.TransferAllBytes(decompressReader, file)
	if err != nil {
		return err
	}

	return nil
}
