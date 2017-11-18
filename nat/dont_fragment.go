package nat

import (
	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/binary/write"
)

const (
	DontFragmentAttributeType    AttributeType = 0x001A
	DontFragmentAttributeRfcName string        = "DONT-FRAGMENT"
	DontFragmentAttributeSize    uint16        = 0
)

func init() {
	RegisterDontFragmentAttribute(DefaultRegistry)
}

func RegisterDontFragmentAttribute(r AttributeRegistry) {
	r.Register(
		DontFragmentAttributeType,
		DontFragmentAttributeRfcName,
		func(r *read.BigEndian, l uint16) (Attribute, error) {
			return ParseDontFragmentAttribute(r, l)
		},
	)
}

func ParseDontFragmentAttribute(r *read.BigEndian, l uint16) (DontFragmentAttribute, error) {
	return DontFragmentAttribute{}, nil
}

func (a DontFragmentAttribute) Print(w *write.BigEndian) error {
	WriteTLVHeader(w, DontFragmentAttributeType, DontFragmentAttributeSize)
	return nil
}

type DontFragmentAttribute struct{}

func (r DontFragmentAttribute) Type() AttributeType {
	return DontFragmentAttributeType
}
