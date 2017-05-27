package util

import (
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

func FileExists(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}

	return true
}
