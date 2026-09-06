package protocol

import (
	"encoding/binary"
	"errors"
)

const (
	_UINT_8_SIZE  = 1
	_UINT_32_SIZE = 4

	_FIRST_NAME_LENGTH_SIZE = _UINT_8_SIZE
	_LAST_NAME_LENGTH_SIZE  = _UINT_8_SIZE
	_DOCUMENT_SIZE          = _UINT_32_SIZE
	_BIRTHDATE_LENGTH_SIZE  = _UINT_8_SIZE
	_NUMBER_SIZE            = _UINT_32_SIZE

	_RECORD_FIXED_SIZE = _FIRST_NAME_LENGTH_SIZE +
		_LAST_NAME_LENGTH_SIZE +
		_DOCUMENT_SIZE +
		_BIRTHDATE_LENGTH_SIZE +
		_NUMBER_SIZE

	MaxTextFieldLength = 255
)

var errInvalidBet = errors.New("invalid bet line")

type Bet struct {
	FirstName []byte
	LastName  []byte
	Document  uint32
	Birthdate []byte
	Number    uint32
}

func (bet Bet) EncodedSize() int {
	return _RECORD_FIXED_SIZE + len(bet.FirstName) + len(bet.LastName) + len(bet.Birthdate)
}

func AppendBet(dst []byte, bet Bet) []byte {
	dst = appendTextField(dst, bet.FirstName)
	dst = appendTextField(dst, bet.LastName)
	dst = binary.BigEndian.AppendUint32(dst, bet.Document)
	dst = appendTextField(dst, bet.Birthdate)
	dst = binary.BigEndian.AppendUint32(dst, bet.Number)

	return dst
}

func DecodeBet(payload []byte, offset int, bet *Bet) (int, error) {
	firstName, offset, err := readTextField(payload, offset)
	if err != nil {
		return 0, err
	}

	lastName, offset, err := readTextField(payload, offset)
	if err != nil {
		return 0, err
	}

	document, offset, err := readUint32(payload, offset)
	if err != nil {
		return 0, err
	}

	birthdate, offset, err := readTextField(payload, offset)
	if err != nil {
		return 0, err
	}

	number, offset, err := readUint32(payload, offset)
	if err != nil {
		return 0, err
	}

	bet.FirstName = firstName
	bet.LastName = lastName
	bet.Document = document
	bet.Birthdate = birthdate
	bet.Number = number

	return offset, nil
}

func appendTextField(dst []byte, value []byte) []byte {
	dst = append(dst, uint8(len(value)))
	dst = append(dst, value...)
	return dst
}

func readTextField(payload []byte, offset int) ([]byte, int, error) {
	length, offset, err := readUint8(payload, offset)
	if err != nil {
		return nil, 0, err
	}

	end := offset + int(length)
	if end > len(payload) {
		return nil, 0, errInvalidBet
	}

	return payload[offset:end], end, nil
}

func readUint32(payload []byte, offset int) (uint32, int, error) {
	end := offset + _UINT_32_SIZE
	if end > len(payload) {
		return 0, 0, errInvalidBet
	}

	number := binary.BigEndian.Uint32(payload[offset:end])
	return number, end, nil
}

func readUint8(payload []byte, offset int) (uint8, int, error) {
	if offset >= len(payload) {
		return 0, 0, errInvalidBet
	}

	number := payload[offset]
	end := offset + _UINT_8_SIZE
	return number, end, nil
}
