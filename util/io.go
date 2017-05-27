package util

import (
	"io"
)

func CloseNoError(c io.Closer) {
	if err := c.Close(); err != nil {
		panic(err)
	}
}
