package nat

import (
	"hash/crc32"

	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/binary/write"
	"github.com/gopherx/base/errors"
)

const (
	FingerprintAttributeType    AttributeType = 0x8028
	FingerprintAttributeRfcName string        = "FINGERPRINT"
	FingerprintAttributeSize    uint16        = 4
)

func init() {
	RegisterFingerprintAttribute(DefaultParser)
}

func RegisterFingerprintAttribute(p *MessageParser) {
	p.Register(
		FingerprintAttributeType,
		FingerprintAttributeRfcName,
		func(r *read.BigEndian, l uint16) (Attribute, error) {
			return ParseFingerprintAttribute(r, l)
		},
		func(w *write.BigEndian, a Attribute) error {
			return errors.Unimplemented(nil, "Use PrintOptions instead")
		},
	)
}

func ParseFingerprintAttribute(r *read.BigEndian, l uint16) (FingerprintAttribute, error) {
	return FingerprintAttribute{r.Uint32()}, nil
}

func PrintFingerprintAttribute(w *write.BigEndian) error {
	raw := w.Dest[0:w.Offset]

	WriteTLVHeader(w, FingerprintAttributeType, FingerprintAttributeSize)

	hash := crc32.NewIEEE()
	hash.Write(raw)
	fp := hash.Sum32() ^ 0x5354554e

	w.Uint32(fp)

	return nil
}

type FingerprintAttribute struct {
	CRC32 uint32
}

func (f FingerprintAttribute) Type() AttributeType {
	return FingerprintAttributeType
}
