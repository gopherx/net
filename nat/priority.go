package nat

import (
	"github.com/gopherx/base/read"
)

const (
	PriorityAttributeType    AttributeType = 0x0024
	PriorityAttributeRfcName string        = "PRIORITY"
)

func init() {
	RegisterPriorityAttribute(DefaultParser)
}

func RegisterPriorityAttribute(p *MessageParser) {
	p.Register(
		PriorityAttributeType,
		PriorityAttributeRfcName,
		func(b []byte) (Attribute, error) {
			return ParsePriorityAttribute(b)
		})
}

func ParsePriorityAttribute(b []byte) (PriorityAttribute, error) {
	return PriorityAttribute{read.Uint32(b)}, nil
}

type PriorityAttribute struct {
	Priority uint32
}
