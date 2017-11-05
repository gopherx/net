package nat

import (
	"time"

	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/binary/write"
)

const (
	LifetimeAttributeType    AttributeType = 0x000D
	LifetimeAttributeRfcName string        = "LIFETIME"
	LifetimeAttributeSize    uint16        = 4
)

func init() {
	RegisterLifetimeAttribute(DefaultParser)
}

func RegisterLifetimeAttribute(p *MessageParser) {
	p.Register(
		LifetimeAttributeType,
		LifetimeAttributeRfcName,
		func(r *read.BigEndian, l uint16) (Attribute, error) {
			return ParseLifetimeAttribute(r, l)
		},
		func(w *write.BigEndian, a Attribute) error {
			return PrintLifetimeAttribute(w, a.(LifetimeAttribute))
		},
	)
}

func ParseLifetimeAttribute(r *read.BigEndian, l uint16) (LifetimeAttribute, error) {
	s := r.Uint32()
	return LifetimeAttribute{time.Duration(s) * time.Second}, nil
}

func PrintLifetimeAttribute(w *write.BigEndian, a LifetimeAttribute) error {
	WriteTLVHeader(w, LifetimeAttributeType, LifetimeAttributeSize)
	w.Uint32(uint32(a.Lifetime.Seconds()))
	return nil
}

type LifetimeAttribute struct {
	Lifetime time.Duration
}

func (r LifetimeAttribute) Type() AttributeType {
	return LifetimeAttributeType
}
