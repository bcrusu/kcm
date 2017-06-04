package util

import (
	"archive/tar"
	"io"
	"path"

	"os"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

func ExtractTar(in io.Reader, outDir string) error {
	tarReader := tar.NewReader(in)

	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		filePath := path.Join(outDir, hdr.Name)

		switch hdr.Typeflag {
		case tar.TypeReg, tar.TypeRegA:
			fileMode := hdr.FileInfo().Mode()
			fileMode |= os.FileMode(0044)
			if fileMode&0100 != 0 {
				fileMode |= os.FileMode(001)
			}

			file, err := CreateFile(filePath, fileMode)
			if err != nil {
				return errors.Wrapf(err, "extract tar: failed to create file '%s'", filePath)
			}

			if _, err := io.Copy(file, tarReader); err != nil {
				return errors.Wrapf(err, "extract tar: failed to write file contents '%s'", filePath)
			}

			if err := file.Close(); err != nil {
				return errors.Wrapf(err, "extract tar: failed to close file '%s'", filePath)
			}
		case tar.TypeDir:
			if err := CreateDirectoryPath(filePath); err != nil {
				return errors.Wrapf(err, "extract tar: failed to create directory '%s'", filePath)
			}
		default:
			glog.Warningf("extract tar: skipped special type header '%s'", hdr.Name)
		}
	}

	return nil
}
