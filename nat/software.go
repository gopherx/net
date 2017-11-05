package nat

import (
	"fmt"

	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/binary/write"
	"github.com/gopherx/base/errors"
)

const (
	SoftwareAttributeType     AttributeType = 0x8022
	SoftwareAttributeRfcName  string        = "SOFTWARE"
	SoftwareAttributeMaxChars uint16        = 128
	SoftwareAttributeMaxBytes uint16        = 763
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
		func(w *write.BigEndian, a Attribute) error {
			return PrintSoftwareAttribute(w, a.(SoftwareAttribute))
		},
	)
}

// ParseSoftwareAttribute parses the bytes into a SoftwareAttribute instance.
func ParseSoftwareAttribute(r *read.BigEndian, l uint16) (SoftwareAttribute, error) {
	res := SoftwareAttribute{}
	if l > SoftwareAttributeMaxBytes {
		return res, errors.InvalidArgument(nil, fmt.Sprintf("too many bytes in text; max=%d current=%d", SoftwareAttributeMaxBytes, l))
	}

	txt := string(r.Bytes(int(l)))
	if len(txt) >= int(SoftwareAttributeMaxChars) {
		return res, errors.InvalidArgument(nil, "too many chars in text; max=%d current=%d", SoftwareAttributeMaxChars, len(txt))
	}

	res.Text = txt
	return res, nil
}

func PrintSoftwareAttribute(w *write.BigEndian, sa SoftwareAttribute) error {
	if len(sa.Text) >= int(SoftwareAttributeMaxChars) {
		return errors.InvalidArgument(nil, "Too many chars; max=%d current=%d", SoftwareAttributeMaxChars, len(sa.Text))
	}

	bytes := []byte(sa.Text)
	if len(bytes) > int(SoftwareAttributeMaxBytes) {
		return errors.InvalidArgument(nil, "Too many bytes; max=%d current=%d", SoftwareAttributeMaxBytes, len(bytes))
	}

	bl := uint16(len(bytes))
	WriteTLVHeader(w, SoftwareAttributeType, bl)
	w.Bytes(bytes)
	WriteTLVPadding(w, bl)

	return nil
}

// SoftwareAttribute is a STUN SOFTWARE attribute.
type SoftwareAttribute struct {
	Text string
}

func (s SoftwareAttribute) Type() AttributeType {
	return SoftwareAttributeType
}
