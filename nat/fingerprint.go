package nat

import (
	"github.com/gopherx/base/read"
)

const (
	FingerprintAttributeType    AttributeType = 0x8028
	FingerprintAttributeRfcName string        = "FINGERPRINT"
)

func init() {
	RegisterFingerprintAttribute(DefaultParser)
}

func RegisterFingerprintAttribute(p *MessageParser) {
	p.Register(
		FingerprintAttributeType,
		FingerprintAttributeRfcName,
		func(b []byte) (Attribute, error) {
			return ParseFingerprintAttribute(b)
		})
}

func ParseFingerprintAttribute(b []byte) (FingerprintAttribute, error) {
	return FingerprintAttribute{read.Uint32(b)}, nil
}

type FingerprintAttribute struct {
	CRC32 uint32
}
