package nat

import (
	"github.com/gopherx/base/errors"
)

const (
	MessageIntegrityType     AttributeType = 0x0008
	MessageIntegrityRfcName  string        = "MESSAGE-INTEGRITY"
	MessageIntegrityMaxBytes int           = 20
)

func init() {
	RegisterParseMessageIntegrityAttribute(DefaultParser)
}

func RegisterParseMessageIntegrityAttribute(p *MessageParser) {
	p.Register(
		MessageIntegrityType,
		MessageIntegrityRfcName,
		func(b []byte) (Attribute, error) {
			return ParseMessageIntegrityAttribute(b)
		})
}

func ParseMessageIntegrityAttribute(b []byte) (MessageIntegrityAttribute, error) {
	res := MessageIntegrityAttribute{}

	if len(b) > MessageIntegrityMaxBytes {
		return res, errors.InvalidArgument(nil, "too many bytes", len(b), MessageIntegrityMaxBytes)
	}

	res.HMAC = b
	return res, nil
}

type MessageIntegrityAttribute struct {
	HMAC []byte
}
