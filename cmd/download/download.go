package download

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"path"

	"bytes"

	"github.com/bcrusu/kcm/libvirt"
	"github.com/bcrusu/kcm/repository"
	"github.com/bcrusu/kcm/util"
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
	volume, err := connection.GetStorageVolume(cluster.StoragePool, cluster.BackingStorageVolume)
	if err != nil {
		return err
	}

	if volume != nil {
		// volume exists - will not download
		return nil
	}

	bytes, err := downloadCoreOS(cluster.CoreOSChannel, cluster.CoreOSVersion)
	if err != nil {
		return err
	}

	_, err = connection.CreateStorageVolume(libvirt.CreateStorageVolumeParams{
		Pool:        cluster.StoragePool,
		Name:        cluster.BackingStorageVolume,
		CapacityGiB: 10,
		Content:     bytes,
	})

	return err
}

func downloadCoreOS(channel, version string) ([]byte, error) {
	url := fmt.Sprintf("https://%s.release.core-os.net/amd64-usr/%s/coreos_production_qemu_image.img.bz2", channel, version)

	downloadReader, err := util.DownloadHTTP(url)
	if err != nil {
		return nil, err
	}
	defer util.CloseNoError(downloadReader)

	decompressReader := bzip2.NewReader(downloadReader)
	buffer := &bytes.Buffer{}

	_, err = io.Copy(buffer, decompressReader)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
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
