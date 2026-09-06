package bet_csv

import (
	"fmt"
	"strconv"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol"
)

var errFieldsAmount = fmt.Errorf("bet line must have %d fields", _FIELDS_AMOUNT)

const (
	_FIELD_SEPARATOR = ','
	_LINE_END        = '\n'

	_FIELDS_AMOUNT     = 5
	_SEPARATORS_AMOUNT = _FIELDS_AMOUNT - 1
	_LINE_END_SIZE     = 1

	_UINT_32_BITS       = 32
	_MAX_UINT_32_DIGITS = 10

	_MAX_FIRST_NAME_SIZE = protocol.MaxTextFieldLength
	_MAX_LAST_NAME_SIZE  = protocol.MaxTextFieldLength
	_MAX_DOCUMENT_SIZE   = _MAX_UINT_32_DIGITS
	_MAX_BIRTHDATE_SIZE  = protocol.MaxTextFieldLength
	_MAX_NUMBER_SIZE     = _MAX_UINT_32_DIGITS

	MaxLineSize = _MAX_FIRST_NAME_SIZE +
		_MAX_LAST_NAME_SIZE +
		_MAX_DOCUMENT_SIZE +
		_MAX_BIRTHDATE_SIZE +
		_MAX_NUMBER_SIZE +
		_SEPARATORS_AMOUNT +
		_LINE_END_SIZE
)

func ParseLine(line []byte, bet *protocol.Bet) error {
	fields, err := splitLine(line)
	if err != nil {
		return err
	}

	firstName, err := parseTextField(fields[0], "first name")
	if err != nil {
		return err
	}

	lastName, err := parseTextField(fields[1], "last name")
	if err != nil {
		return err
	}

	document, err := parseUint32(fields[2], "document")
	if err != nil {
		return err
	}

	birthdate, err := parseTextField(fields[3], "birthdate")
	if err != nil {
		return err
	}

	number, err := parseUint32(fields[4], "number")
	if err != nil {
		return err
	}

	bet.FirstName = firstName
	bet.LastName = lastName
	bet.Document = document
	bet.Birthdate = birthdate
	bet.Number = number

	return nil
}

func AppendLine(dst []byte, bet protocol.Bet) []byte {
	dst = append(dst, bet.FirstName...)
	dst = append(dst, _FIELD_SEPARATOR)
	dst = append(dst, bet.LastName...)
	dst = append(dst, _FIELD_SEPARATOR)
	dst = strconv.AppendUint(dst, uint64(bet.Document), 10)
	dst = append(dst, _FIELD_SEPARATOR)
	dst = append(dst, bet.Birthdate...)
	dst = append(dst, _FIELD_SEPARATOR)
	dst = strconv.AppendUint(dst, uint64(bet.Number), 10)

	return append(dst, _LINE_END)
}

func splitLine(line []byte) ([_FIELDS_AMOUNT][]byte, error) {
	var fields [_FIELDS_AMOUNT][]byte

	fieldIndex := 0
	fieldStart := 0

	for i, char := range line {
		if char != _FIELD_SEPARATOR {
			continue
		}

		if fieldIndex == _SEPARATORS_AMOUNT {
			return fields, errFieldsAmount
		}

		fields[fieldIndex] = line[fieldStart:i]
		fieldIndex++
		fieldStart = i + 1
	}

	if fieldIndex != _SEPARATORS_AMOUNT {
		return fields, errFieldsAmount
	}

	fields[fieldIndex] = line[fieldStart:]

	return fields, nil
}

func parseTextField(field []byte, fieldName string) ([]byte, error) {
	if len(field) == 0 {
		return nil, fmt.Errorf("empty %s field", fieldName)
	}

	if len(field) > protocol.MaxTextFieldLength {
		return nil, fmt.Errorf("%s field exceeds %d bytes", fieldName, protocol.MaxTextFieldLength)
	}

	return field, nil
}

func parseUint32(field []byte, fieldName string) (uint32, error) {
	value, err := strconv.ParseUint(string(field), 10, _UINT_32_BITS)
	if err != nil {
		return 0, fmt.Errorf("%s field is not a valid number: %q", fieldName, field)
	}

	return uint32(value), nil
}
