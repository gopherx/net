package nat

import (
	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/binary/write"
)

const (
	ErrorCodeAttributeType    AttributeType = 0x0009
	ErrorCodeAttributeRfcName string        = "ERROR-CODE"
)

func init() {
	RegisterErrorCodeAttribute(DefaultParser)
}

func RegisterErrorCodeAttribute(p *MessageParser) {
	p.Register(
		ErrorCodeAttributeType,
		ErrorCodeAttributeRfcName,
		func(r *read.BigEndian, l uint16) (Attribute, error) {
			return ParseErrorCodeAttribute(r, l)
		},
		func(w *write.BigEndian, a Attribute) error {
			return PrintErrorCodeAttribute(w, a.(ErrorCodeAttribute))
		},
	)
}

func ParseErrorCodeAttribute(r *read.BigEndian, l uint16) (ErrorCodeAttribute, error) {
	code := ErrorCodeAttribute{}

	tmp := r.Uint32()
	number := byte(tmp & 0x000000FF)
	class := byte((tmp & 0x00000700) >> 8)

	l -= uint16(4)
	reason, err := Read127CharString(r, l)
	if err != nil {
		return code, err
	}

	code.Class = class
	code.Number = number
	code.Reason = reason
	return code, nil
}

func PrintErrorCodeAttribute(w *write.BigEndian, a ErrorCodeAttribute) error {
	bytes, err := Check127CharString(a.Reason)
	if err != nil {
		return err
	}

	written := uint16(4 + len(bytes))

	WriteTLVHeader(w, ErrorCodeAttributeType, written)
	class := uint32(a.Class) << 8
	v := class | uint32(a.Number)
	w.Uint32(v)
	w.Bytes(bytes)
	WriteTLVPadding(w, written)
	return nil
}

type ErrorCodeAttribute struct {
	Class  byte
	Number byte
	Reason string
}

func (f ErrorCodeAttribute) Type() AttributeType {
	return ErrorCodeAttributeType
}
