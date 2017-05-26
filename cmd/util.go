package cmd

import (
	"bufio"
	"compress/bzip2"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"time"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

func getHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func downloadHTTP(url string) (io.ReadCloser, error) {
	client := &http.Client{
		Timeout: time.Hour,
	}

	response, err := client.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to download '%s'", url)
	}

	if response.StatusCode != 200 {
		return nil, errors.Errorf("failed to download '%s'. Error code: %d", url, response.StatusCode)
	}

	return response.Body, nil
}

func createFile(fileName string, perm os.FileMode) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create file '%s'", fileName)
	}

	return file, nil
}

func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}

	return true
}

func closeNoError(c io.Closer) {
	if err := c.Close(); err != nil {
		panic(err)
	}
}

func downloadCoreOSImage(version, channel string, outDir string) error {
	fileName := fmt.Sprintf("coreos_production_qemu_%s.qcow2", version)
	filePath := path.Join(outDir, fileName)

	if fileExists(filePath) {
		glog.Infof("skipped CoreOS image download. File '%s' exists.", filePath)
		return nil
	}

	url := fmt.Sprintf("https://%s.release.core-os.net/amd64-usr/%s/coreos_production_qemu_image.img.bz2", channel, version)

	downloadReader, err := downloadHTTP(url)
	if err != nil {
		return err
	}
	defer closeNoError(downloadReader)

	decompressReader := bzip2.NewReader(downloadReader)

	file, err := createFile(filePath, 0644)
	if err != nil {
		return err
	}
	defer closeNoError(file)

	err = writeAll(decompressReader, file)
	if err != nil {
		return err
	}

	return nil
}

func writeAll(in io.Reader, out io.Writer) error {
	buffered := bufio.NewWriter(out)

	buf := make([]byte, 4096)
	for {
		n, err := in.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		if _, err := buffered.Write(buf[:n]); err != nil {
			return err
		}
	}

	if err := buffered.Flush(); err != nil {
		return err
	}

	return nil
}
