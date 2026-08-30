package protocol

type MessageType uint8

const (
	MsgStartTransmission MessageType = 0x01
	MsgBet               MessageType = 0x02
	MsgEndTransmission   MessageType = 0x03
	MsgOk                MessageType = 0x10
	MsgWinners           MessageType = 0x11
	MsgError             MessageType = 0x12
)

const (
	_TYPE_FIELD_SIZE   = 1
	_LENGTH_FIELD_SIZE = 2
	_HEADER_SIZE       = _TYPE_FIELD_SIZE + _LENGTH_FIELD_SIZE
	_MAX_UINT_16       = 65535
	_MAX_PAYLOAD_SIZE  = _MAX_UINT_16
)

type Message struct {
	Type    MessageType
	Payload []byte
}

func NewMessage(messageType MessageType, payload []byte) *Message {
	return &Message{
		messageType,
		payload,
	}
}

func (messageType MessageType) String() string {
	switch messageType {
	case MsgStartTransmission:
		return "START_TRANSMISSION"
	case MsgBet:
		return "BET"
	case MsgEndTransmission:
		return "END_TRANSMISSION"
	case MsgOk:
		return "OK"
	case MsgWinners:
		return "WINNERS"
	case MsgError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
