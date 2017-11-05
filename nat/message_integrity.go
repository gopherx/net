package nat

import (
	"crypto/hmac"
	"crypto/sha1"

	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/binary/write"
	"github.com/gopherx/base/errors"
)

const (
	MessageIntegrityAttributeType     AttributeType = 0x0008
	MessageIntegrityAttributeRfcName  string        = "MESSAGE-INTEGRITY"
	MessageIntegrityAttributeMaxBytes int           = 20
	MessageIntegrityAttributeSize     uint16        = 20
)

func init() {
	RegisterParseMessageIntegrityAttribute(DefaultParser)
}

func RegisterParseMessageIntegrityAttribute(p *MessageParser) {
	p.Register(
		MessageIntegrityAttributeType,
		MessageIntegrityAttributeRfcName,
		func(r *read.BigEndian, l uint16) (Attribute, error) {
			return ParseMessageIntegrityAttribute(r, l)
		},
		func(w *write.BigEndian, a Attribute) error {
			return errors.Unimplemented(nil, "Use PrintOptions instead")
		},
	)
}

func ParseMessageIntegrityAttribute(r *read.BigEndian, l uint16) (MessageIntegrityAttribute, error) {
	res := MessageIntegrityAttribute{}

	res.HMAC = r.Bytes(int(MessageIntegrityAttributeSize))
	return res, nil
}

func PrintMessageIntegrityAttribute(w *write.BigEndian, key []byte) error {
	raw := w.Dest[0:w.Offset]
	WriteTLVHeader(w, MessageIntegrityAttributeType, MessageIntegrityAttributeSize)
	mac := hmac.New(sha1.New, key)
	mac.Write(raw)
	w.Bytes(mac.Sum(nil))

	return nil
}

type MessageIntegrityAttribute struct {
	HMAC []byte
}

func (m MessageIntegrityAttribute) Type() AttributeType {
	return MessageIntegrityAttributeType
}
