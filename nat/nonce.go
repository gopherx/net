package nat

import (
	"fmt"

	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/binary/write"
	"github.com/gopherx/base/errors"
)

const (
	NonceAttributeType     AttributeType = 0x0015
	NonceAttributeRfcName  string        = "NONCE"
	NonceAttributeMaxChars uint16        = 128
	NonceAttributeMaxBytes uint16        = 763
)

func init() {
	RegisterNonceAttribute(DefaultParser)
}

func RegisterNonceAttribute(p *MessageParser) {
	p.Register(
		NonceAttributeType,
		NonceAttributeRfcName,
		func(r *read.BigEndian, l uint16) (Attribute, error) {
			return ParseNonceAttribute(r, l)
		},
		func(w *write.BigEndian, a Attribute) error {
			return PrintNonceAttribute(w, a.(NonceAttribute))
		},
	)
}

func ParseNonceAttribute(r *read.BigEndian, l uint16) (NonceAttribute, error) {
	Nonce := NonceAttribute{}
	if l > NonceAttributeMaxBytes {
		return Nonce, errors.InvalidArgument(nil, fmt.Sprintf("too many bytes in Nonce; max=%d current=%d", NonceAttributeMaxBytes, l))
	}

	txt := string(r.Bytes(int(l)))
	if uint16(len(txt)) > NonceAttributeMaxChars {
		return Nonce, errors.InvalidArgument(nil, fmt.Sprintf("too many chars in Nonce; max=%d current=%d", NonceAttributeMaxChars, len(txt)))
	}

	Nonce.Nonce = txt
	return Nonce, nil
}

func PrintNonceAttribute(w *write.BigEndian, a NonceAttribute) error {
	bytes := []byte(a.Nonce)
	WriteTLVHeader(w, NonceAttributeType, uint16(len(bytes)))
	w.Bytes(bytes)
	WriteTLVPadding(w, uint16(len(bytes)))
	return nil
}

type NonceAttribute struct {
	Nonce string
}

func (f NonceAttribute) Type() AttributeType {
	return NonceAttributeType
}
