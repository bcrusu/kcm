package libvirt

import (
	"github.com/libvirt/libvirt-go"
)

const streamSendMaxSize = 16000000 // must be <= VIR_NET_MESSAGE_PAYLOAD_MAX

func streamSendAll(stream *libvirt.Stream, bytes []byte) error {
	for {
		chunkSize := streamSendMaxSize
		if chunkSize > len(bytes) {
			chunkSize = len(bytes)
		}

		if chunkSize == 0 {
			break
		}

		chunk := bytes[:chunkSize]
		bytes = bytes[chunkSize:]

		_, err := stream.Send(chunk)
		if err != nil {
			return err
		}
	}

	return nil
}
