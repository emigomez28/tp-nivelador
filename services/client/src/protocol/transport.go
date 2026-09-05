package protocol

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/safe_socket"
)

func NewSendBuffer() []byte {
	initialLen := 0
	initialCapacity := _HEADER_SIZE + _MAX_PAYLOAD_SIZE
	return make([]byte, initialLen, initialCapacity)
}

func SendMessage(writer io.Writer, msg Message, buffer []byte) error {
	if len(msg.Payload) > _MAX_PAYLOAD_SIZE {
		return fmt.Errorf("payload exceeds %d bytes", _MAX_PAYLOAD_SIZE)
	}

	dataToSend := append(buffer[:0], byte(msg.Type))
	dataToSend = binary.BigEndian.AppendUint16(dataToSend, uint16(len(msg.Payload)))
	dataToSend = append(dataToSend, msg.Payload...)

	return safe_socket.SendAll(writer, dataToSend)
}

func RecvMessage(reader io.Reader) (Message, error) {
	header, err := safe_socket.RecvAll(reader, _HEADER_SIZE)
	if err != nil {
		return Message{}, err
	}

	msgType := MessageType(header[0])
	payloadSize := int(binary.BigEndian.Uint16(header[_TYPE_FIELD_SIZE:]))
	payload, err := safe_socket.RecvAll(reader, payloadSize)
	if err != nil {
		return Message{}, err
	}

	msg := NewMessage(msgType, payload)
	return msg, nil
}
