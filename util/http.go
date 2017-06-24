package util

import (
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

func DownloadHTTP(url string) (io.ReadCloser, error) {
	client := &http.Client{
		Timeout: time.Hour,
	}

	response, err := client.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "http: failed to download '%s'", url)
	}

	if response.StatusCode != 200 {
		return nil, errors.Errorf("http: failed to download '%s'. Error code: %d", url, response.StatusCode)
	}

	return response.Body, nil
}
