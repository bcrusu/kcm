package util

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

func CreateFile(fileName string, perm os.FileMode) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create file '%s'", fileName)
	}

	return file, nil
}

func DirectoryExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	if !stat.IsDir() {
		return false, errors.Errorf("path is not a directory '%s'", path)
	}

	return true, nil
}

func FileExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	if stat.IsDir() {
		return false, errors.Errorf("path is not a file '%s'", path)
	}

	return true, nil
}

func CreateDirectoryPath(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return errors.Wrapf(err, "failed to create directory '%s'", path)
	}

	return nil
}

func RemoveDirectory(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return errors.Wrapf(err, "failed to remove directory '%s'", path)
	}

	return nil
}

func WriteFile(path string, data []byte) error {
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return errors.Wrapf(err, "failed to write to file '%s'", path)
	}

	return nil
}

func WriteExecutableFile(path string, data []byte) error {
	if err := ioutil.WriteFile(path, data, 0655); err != nil {
		return errors.Wrapf(err, "failed to write to file '%s'", path)
	}

	return nil
}
