package safe_socket

import "io"

func SendAll(socket io.Writer, bytes []byte) error {
	totalSent := 0

	for totalSent < len(bytes) {
		sent, err := socket.Write(bytes[totalSent:])
		totalSent += sent

		if err != nil {
			return err
		}
	}

	return nil
}

func RecvAll(socket io.Reader, size int) ([]byte, error) {
	buff := make([]byte, size)
	totalRead := 0

	for totalRead < size {
		read, err := socket.Read(buff[totalRead:])
		totalRead += read

		if totalRead == size {
			return buff, nil
		}

		if err != nil {
			return nil, err
		}
	}

	return buff, nil
}
