package nat

import (
	"fmt"

	"github.com/gopherx/base/errors"
)

const (
	SoftwareAttributeType     AttributeType = 0x8022
	SoftwareAttributeRfcName  string        = "SOFTWARE"
	SoftwareAttributeMaxChars int           = 128
	SoftwareAttributeMaxBytes int           = 763
)

func init() {
	RegisterSoftwareAttribute(DefaultParser)
}

func RegisterSoftwareAttribute(p *MessageParser) {
	p.Register(
		SoftwareAttributeType,
		SoftwareAttributeRfcName,
		func(b []byte) (Attribute, error) {
			return ParseSoftwareAttribute(b)
		})
}

// ParseSoftwareAttribute parses the bytes into a SoftwareAttribute instance.
func ParseSoftwareAttribute(b []byte) (SoftwareAttribute, error) {
	res := SoftwareAttribute{}
	if len(b) > SoftwareAttributeMaxBytes {
		return res, errors.InvalidArgument(nil, fmt.Sprintf("too many bytes in text; max=%d current=%d", SoftwareAttributeMaxBytes, len(b)))
	}

	txt := string(b)
	if len(txt) >= SoftwareAttributeMaxChars {
		return res, errors.InvalidArgument(nil, "too many chars in text; max=%d current=%d", SoftwareAttributeMaxChars, len(txt))
	}

	res.Text = txt
	return res, nil
}

// SoftwareAttribute is a STUN SOFTWARE attribute.
type SoftwareAttribute struct {
	Text string
}
