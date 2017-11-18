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
	RegisterPriorityAttribute(DefaultRegistry)
}

func RegisterPriorityAttribute(r AttributeRegistry) {
	r.Register(
		PriorityAttributeType,
		PriorityAttributeRfcName,
		func(r *read.BigEndian, l uint16) (Attribute, error) {
			return ParsePriorityAttribute(r, l)
		},
	)
}

func ParsePriorityAttribute(r *read.BigEndian, l uint16) (PriorityAttribute, error) {
	return PriorityAttribute{r.Uint32()}, nil
}

func (p PriorityAttribute) Print(w *write.BigEndian) error {
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
