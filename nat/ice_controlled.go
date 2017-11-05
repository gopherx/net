package nat

import (
	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/binary/write"
)

const (
	IceControlledAttributeType    AttributeType = 0x8029
	IceControlledAttributeRfcName string        = "ICE-CONTROLLED"
	IceControlledAttributeSize    uint16        = 0x08
)

func init() {
	RegisterIceControlledAttribute(DefaultParser)
}

func RegisterIceControlledAttribute(p *MessageParser) {
	p.Register(
		IceControlledAttributeType,
		IceControlledAttributeRfcName,
		func(r *read.BigEndian, l uint16) (Attribute, error) {
			return ParseIceControlledAttribute(r, l)
		},
		func(w *write.BigEndian, a Attribute) error {
			return PrintIceControlledAttribute(w, a.(IceControlledAttribute))
		},
	)
}

func ParseIceControlledAttribute(r *read.BigEndian, l uint16) (Attribute, error) {
	return IceControlledAttribute{r.Uint64()}, nil
}

func PrintIceControlledAttribute(w *write.BigEndian, i IceControlledAttribute) error {
	WriteTLVHeader(w, IceControlledAttributeType, IceControlledAttributeSize)
	w.Uint64(i.TieBreaker)
	return nil
}

type IceControlledAttribute struct {
	TieBreaker uint64
}

func (i IceControlledAttribute) Type() AttributeType {
	return IceControlledAttributeType
}
