package create

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"path"

	"github.com/bcrusu/kcm/util"
	"github.com/golang/glog"
)

func downloadKubernetes(version string, outDir string) error {
	url := fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/v%s/kubernetes-server-linux-amd64.tar.gz", version)

	glog.Infof("downloading Kubernetes from '%s' to path '%s'", url, "")

	downloadReader, err := util.DownloadHTTP(url)
	if err != nil {
		return err
	}
	defer util.CloseNoError(downloadReader)

	gzipReader, err := gzip.NewReader(downloadReader)
	if err != nil {
		return err
	}
	defer util.CloseNoError(gzipReader)

	buffered := bufio.NewReader(gzipReader)
	return util.ExtractTar(buffered, outDir)
}

func DownloadKubernetes(version string, kubeCacheDir string) error {
	kubePath := path.Join(kubeCacheDir, version)

	if _, err := os.Stat(kubePath); err != nil {
		if os.IsNotExist(err) {
			return downloadKubernetes(version, kubePath)
		}
		return err
	}

	// kubernetes version already on disk - no need to download
	return nil
}
