package protocol

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/safe_socket"
)

func SendMessage(writer io.Writer, msg *Message) error {
	if len(msg.Payload) > _MAX_PAYLOAD_SIZE {
		return fmt.Errorf("payload exceeds %d bytes", _MAX_PAYLOAD_SIZE)
	}

	dataToSend := make([]byte, _HEADER_SIZE+len(msg.Payload))
	dataToSend[0] = byte(msg.Type)
	binary.BigEndian.PutUint16(dataToSend[_TYPE_FIELD_SIZE:_HEADER_SIZE], uint16(len(msg.Payload)))
	copy(dataToSend[_HEADER_SIZE:], msg.Payload)

	return safe_socket.SendAll(writer, dataToSend)
}

func RecvMessage(reader io.Reader) (*Message, error) {
	header, err := safe_socket.RecvAll(reader, _HEADER_SIZE)
	if err != nil {
		return nil, err
	}

	msgType := MessageType(header[0])
	payloadSize := int(binary.BigEndian.Uint16(header[_TYPE_FIELD_SIZE:]))
	payload, err := safe_socket.RecvAll(reader, payloadSize)
	if err != nil {
		return nil, err
	}

	msg := NewMessage(msgType, payload)
	return msg, nil
}
