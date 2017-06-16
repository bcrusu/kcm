package download

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"path"

	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
	"github.com/bcrusu/kcm/util"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

func DownloadPrerequisites(connection *libvirt.Connection, cluster repository.Cluster, cacheDir string) error {
	if err := downloadKubernetes(cluster.KubernetesVersion, path.Join(cacheDir, "kubernetes")); err != nil {
		return err
	}

	if err := downloadCNI(cluster.CNIVersion, path.Join(cacheDir, "cni")); err != nil {
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
		if err := downloadCoreOS(cluster.CoreOSChannel, cluster.CoreOSVersion, cluster.BackingStorageVolume, poolPath); err != nil {
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

func downloadTarGz(url string, outDir string) error {
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

func downloadKubernetes(version string, cacheDir string) error {
	kubePath := path.Join(cacheDir, version)

	exists, err := util.DirectoryExists(kubePath)
	if err != nil {
		return err
	}

	if exists {
		// kubernetes version already on disk - no need to download
		return nil
	}

	url := fmt.Sprintf("https://dl.k8s.io/release/v%s/kubernetes-server-linux-amd64.tar.gz", version)
	return downloadTarGz(url, kubePath)
}

func downloadCNI(version string, cacheDir string) error {
	cniPath := path.Join(cacheDir, version)

	exists, err := util.DirectoryExists(cniPath)
	if err != nil {
		return err
	}

	if exists {
		// CNI already on disk - no need to download
		return nil
	}

	url := fmt.Sprintf("https://dl.k8s.io/network-plugins/cni-amd64-%s.tar.gz", version)
	return downloadTarGz(url, cniPath)
}
