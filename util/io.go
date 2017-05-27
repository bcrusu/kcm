package util

import (
	"bufio"
	"io"
)

func TransferAllBytes(in io.Reader, out io.Writer) error {
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

func CloseNoError(c io.Closer) {
	if err := c.Close(); err != nil {
		panic(err)
	}
}
