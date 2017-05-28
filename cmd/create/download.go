package create

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
	"github.com/bcrusu/kcm/util"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

func DownloadPrerequisites(connection *libvirt.Connection, cluster repository.Cluster, kubernetesCacheDir string) error {
	if err := downloadKubernetes(cluster.KubernetesVersion, kubernetesCacheDir); err != nil {
		return err
	}

	if err := downloadBackingStorageImage(connection, cluster); err != nil {
		return err
	}

	return nil
}

func downloadBackingStorageImage(connection *libvirt.Connection, cluster repository.Cluster) error {
	poolPath, err := getStoragePoolTargetPath(connection, cluster.StoragePool)
	if err != nil {
		return err
	}

	volume, err := connection.GetStorageVolume(cluster.StoragePool, cluster.BackingStorageVolume)
	if err != nil {
		return err
	}

	if volume == nil {
		if err := downloadCoreOS(cluster.CoreOSVersion, cluster.CoreOSChannel, cluster.BackingStorageVolume, poolPath); err != nil {
			return err
		}
	}

	// volume exists - will not download
	return nil
}

func getStoragePoolTargetPath(connection *libvirt.Connection, pool string) (string, error) {
	p, err := connection.GetStoragePool(pool)
	if err != nil {
		return "", err
	}

	if p == nil {
		return "", errors.Errorf("could not find storage pool '%s'", pool)
	}

	poolType := p.Type()
	if poolType != "dir" && poolType != "fs" {
		return "", errors.Errorf("unsupported storage pool type '%s'", poolType)
	}

	poolPath := p.Target().Path()
	if poolPath == "" {
		return "", errors.Errorf("invalid storage pool '%s'. Target path is empty", poolType)
	}

	return poolPath, nil
}

func downloadCoreOS(channel, version string, volumeName string, outDir string) error {
	filePath := path.Join(outDir, volumeName)

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

func runDownloadKubernetes(version string, outDir string) error {
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

func downloadKubernetes(version string, kubeCacheDir string) error {
	kubePath := path.Join(kubeCacheDir, version)

	if _, err := os.Stat(kubePath); err != nil {
		if os.IsNotExist(err) {
			return runDownloadKubernetes(version, kubePath)
		}
		return err
	}

	// kubernetes version already on disk - no need to download
	return nil
}
