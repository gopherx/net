package nat

import (
	"hash/crc32"

	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/binary/write"
)

const (
	FingerprintAttributeType    AttributeType = 0x8028
	FingerprintAttributeRfcName string        = "FINGERPRINT"
	FingerprintAttributeSize    uint16        = 4
)

func init() {
	RegisterFingerprintAttribute(DefaultRegistry)
}

func RegisterFingerprintAttribute(r AttributeRegistry) {
	r.Register(
		FingerprintAttributeType,
		FingerprintAttributeRfcName,
		func(r *read.BigEndian, l uint16) (Attribute, error) {
			return ParseFingerprintAttribute(r, l)
		},
	)
}

func ParseFingerprintAttribute(r *read.BigEndian, l uint16) (FingerprintAttribute, error) {
	return FingerprintAttribute{r.Uint32()}, nil
}

func (f FingerprintAttribute) Print(w *write.BigEndian) error {
	panic("don't call this!!! fingerprinting is implemented in the MessagePrinter")
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
