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
	RegisterLifetimeAttribute(DefaultRegistry)
}

func RegisterLifetimeAttribute(r AttributeRegistry) {
	r.Register(
		LifetimeAttributeType,
		LifetimeAttributeRfcName,
		func(r *read.BigEndian, l uint16) (Attribute, error) {
			return ParseLifetimeAttribute(r, l)
		},
	)
}

func ParseLifetimeAttribute(r *read.BigEndian, l uint16) (LifetimeAttribute, error) {
	s := r.Uint32()
	return LifetimeAttribute{time.Duration(s) * time.Second}, nil
}

func (a LifetimeAttribute) Print(w *write.BigEndian) error {
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
