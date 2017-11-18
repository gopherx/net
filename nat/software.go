package nat

import (
	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/binary/write"
)

const (
	SoftwareAttributeType    AttributeType = 0x8022
	SoftwareAttributeRfcName string        = "SOFTWARE"
)

func init() {
	RegisterSoftwareAttribute(DefaultRegistry)
}

func RegisterSoftwareAttribute(r AttributeRegistry) {
	r.Register(
		SoftwareAttributeType,
		SoftwareAttributeRfcName,
		func(r *read.BigEndian, l uint16) (Attribute, error) {
			return ParseSoftwareAttribute(r, l)
		},
	)
}

// ParseSoftwareAttribute parses the bytes into a SoftwareAttribute instance.
func ParseSoftwareAttribute(r *read.BigEndian, l uint16) (SoftwareAttribute, error) {
	sw, err := Read127CharString(r, l)
	return SoftwareAttribute{sw}, err
}

func (sa SoftwareAttribute) Print(w *write.BigEndian) error {
	bytes, err := Check127CharString(sa.Text)
	if err != nil {
		return err
	}

	WriteTLVHeader(w, SoftwareAttributeType, uint16(len(bytes)))
	w.Bytes(bytes)
	WriteTLVPadding(w, uint16(len(bytes)))

	return nil
}

// SoftwareAttribute is a STUN SOFTWARE attribute.
type SoftwareAttribute struct {
	Text string
}

func (s SoftwareAttribute) Type() AttributeType {
	return SoftwareAttributeType
}
