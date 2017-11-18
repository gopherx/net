package nat

import (
	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/binary/write"
)

const (
	EvenPortAttributeType    AttributeType = 0x0018
	EvenPortAttributeRfcName string        = "EVEN-PORT"
	EvenPortAttributeSize    uint16        = 4
	EvenPortAttributeMask    byte          = 0x80
)

func init() {
	RegisterEvenPortAttribute(DefaultRegistry)
}

func RegisterEvenPortAttribute(r AttributeRegistry) {
	r.Register(
		EvenPortAttributeType,
		EvenPortAttributeRfcName,
		func(r *read.BigEndian, l uint16) (Attribute, error) {
			return ParseEvenPortAttribute(r, l)
		},
	)
}

func ParseEvenPortAttribute(r *read.BigEndian, l uint16) (EvenPortAttribute, error) {
	b0 := r.Byte()
	return EvenPortAttribute{b0&EvenPortAttributeMask == EvenPortAttributeMask}, nil
}

func (a EvenPortAttribute) Print(w *write.BigEndian) error {
	WriteTLVHeader(w, EvenPortAttributeType, EvenPortAttributeSize)
	v := byte(0)
	if a.Even {
		v = EvenPortAttributeMask
	}
	w.Byte(v)
	WriteTLVPadding(w, 1)
	return nil
}

type EvenPortAttribute struct {
	Even bool
}

func (r EvenPortAttribute) Type() AttributeType {
	return EvenPortAttributeType
}
