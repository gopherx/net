package nat

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/binary/write"
)

const (
	MessageIntegrityAttributeType     AttributeType = 0x0008
	MessageIntegrityAttributeRfcName  string        = "MESSAGE-INTEGRITY"
	MessageIntegrityAttributeMaxBytes int           = 20
	MessageIntegrityAttributeSize     uint16        = 20
)

func init() {
	RegisterParseMessageIntegrityAttribute(DefaultRegistry)
}

func RegisterParseMessageIntegrityAttribute(r AttributeRegistry) {
	r.Register(
		MessageIntegrityAttributeType,
		MessageIntegrityAttributeRfcName,
		func(r *read.BigEndian, l uint16) (Attribute, error) {
			return ParseMessageIntegrityAttribute(r, l)
		},
	)
}

func ParseMessageIntegrityAttribute(r *read.BigEndian, l uint16) (MessageIntegrityAttribute, error) {
	res := MessageIntegrityAttribute{}
	res.HMAC = r.Bytes(int(MessageIntegrityAttributeSize))
	res.raw = r.Read[0 : uint16(len(r.Read))-TLVHeaderSize-MessageIntegrityAttributeSize]
	return res, nil
}

func (a MessageIntegrityAttribute) IsCorrupt(username, password, realm string) bool {
	// We use different keys for long vs short - term credentials. If we have a valid realm then
	// this request uses long term credentials.
	key := []byte(password)
	if len(realm) > 0 {
		//...long term credentials; calculate the key
		h := md5.New()
		tmp := fmt.Sprintf("%s:%s:%s", username, realm, password)
		io.WriteString(h, tmp)
		key = h.Sum(nil)
	}

	// Change the size value directly in the buffer since the RFC requires it to include
	// at most the MESSAGE-INTEGRITY attribute itself and not any following attributes.
	// (and remember that the size of the message doesn't include the message header)
	sizeBytes := a.raw[2:4]
	currSize := read.Uint16(sizeBytes)
	// Since HeaderSize and MessageIntegrityAttribute size are equal we can just skip
	// them from the size calculation since HeaderSize would be removed and
	// MessageIntegrityAttibuteSize would be added.
	// The real calc should look like this:
	// size = len(r.raw) - HeaderSize + TLVHeaderSize + MessageIntegrityAttributeSize
	tmpSize := uint16(len(a.raw)) + TLVHeaderSize
	write.Uint16(sizeBytes, tmpSize)
	defer write.Uint16(sizeBytes, currSize)

	mac := hmac.New(sha1.New, key)
	mac.Write(a.raw)
	expected := mac.Sum(nil)

	return !hmac.Equal(a.HMAC, expected)
}

func (m MessageIntegrityAttribute) Print(w *write.BigEndian) error {
	panic("don't call this; message integrity is handled by the printer")
}

func PrintMessageIntegrityAttribute(w *write.BigEndian, key []byte) error {
	//...message integrity is calculated from the start of the message to the end of the
	// message integrity attribute; we therefore write the length of the message as if there
	// is no attribute beyond the message integrity attribute. If there is a fingerprint
	// attribute then the size needs adjustment.
	tmpSize := uint16(w.Offset) - HeaderSize + TLVHeaderSize + MessageIntegrityAttributeSize

	w.Uint16At(2, tmpSize)

	// There is no need to undo the 'temporary' write of size since it will be overwritten
	// by the normal print flow.

	raw := w.Dest[0:w.Offset]
	WriteTLVHeader(w, MessageIntegrityAttributeType, MessageIntegrityAttributeSize)
	mac := hmac.New(sha1.New, key)
	mac.Write(raw)
	w.Bytes(mac.Sum(nil))

	return nil
}

type MessageIntegrityAttribute struct {
	HMAC []byte
	raw  []byte
}

func (m MessageIntegrityAttribute) Type() AttributeType {
	return MessageIntegrityAttributeType
}
