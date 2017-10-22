package nat

import (
	"github.com/gopherx/base/read"
)

const (
	IceControlledAttributeType    AttributeType = 0x8029
	IceControlledAttributeRfcName string        = "ICE-CONTROLLED"
)

func init() {
	RegisterIceControlledAttribute(DefaultParser)
}

func RegisterIceControlledAttribute(p *MessageParser) {
	p.Register(
		IceControlledAttributeType,
		IceControlledAttributeRfcName,
		func(b []byte) (Attribute, error) {
			return ParseIceControlledAttribute(b)
		})
}

func ParseIceControlledAttribute(b []byte) (Attribute, error) {
	return IceControlledAttribute{read.Uint64(b)}, nil
}

type IceControlledAttribute struct {
	TieBreaker uint64
}
