package nat

import (
	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/binary/write"
)

const (
	PriorityAttributeType    AttributeType = 0x0024
	PriorityAttributeRfcName string        = "PRIORITY"
	PriorityAttributeSize    uint16        = 0x04
)

func init() {
	RegisterPriorityAttribute(DefaultParser)
}

func RegisterPriorityAttribute(p *MessageParser) {
	p.Register(
		PriorityAttributeType,
		PriorityAttributeRfcName,
		func(r *read.BigEndian, l uint16) (Attribute, error) {
			return ParsePriorityAttribute(r, l)
		},
		func(w *write.BigEndian, a Attribute) error {
			return PrintPriorityAttribute(w, a.(PriorityAttribute))
		},
	)
}

func ParsePriorityAttribute(r *read.BigEndian, l uint16) (PriorityAttribute, error) {
	return PriorityAttribute{r.Uint32()}, nil
}

func PrintPriorityAttribute(w *write.BigEndian, p PriorityAttribute) error {
	WriteTLVHeader(w, PriorityAttributeType, PriorityAttributeSize)
	w.Uint32(p.Priority)
	return nil
}

type PriorityAttribute struct {
	Priority uint32
}

func (p PriorityAttribute) Type() AttributeType {
	return PriorityAttributeType
}
