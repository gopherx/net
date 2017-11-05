package nat

import (
	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/binary/write"
	"github.com/gopherx/base/errors"
)

const (
	UsernameAttributeType     AttributeType = 0x0006
	UsernameAttributeRfcName  string        = "USERNAME"
	UsernameAttributeMaxBytes uint16        = 512
)

func init() {
	RegisterUsernameAttribute(DefaultParser)
}

func RegisterUsernameAttribute(p *MessageParser) {
	p.Register(
		UsernameAttributeType,
		UsernameAttributeRfcName,
		func(r *read.BigEndian, l uint16) (Attribute, error) {
			return ParseUsernameAttribute(r, l)
		},
		func(w *write.BigEndian, a Attribute) error {
			return PrintUsernameAttribute(w, a.(UsernameAttribute))
		},
	)
}

func ParseUsernameAttribute(r *read.BigEndian, l uint16) (UsernameAttribute, error) {
	res := UsernameAttribute{}

	if l > UsernameAttributeMaxBytes {
		return res, errors.InvalidArgument(nil, "too many bytes", l, UsernameAttributeMaxBytes)
	}

	res.Username = string(r.Bytes(int(l)))
	return res, nil
}

func PrintUsernameAttribute(w *write.BigEndian, u UsernameAttribute) error {
	unb := []byte(u.Username)
	if len(unb) > int(UsernameAttributeMaxBytes) {
		return errors.InvalidArgument(nil, "too many bytes in username", len(u.Username), len(unb))
	}

	WriteTLVHeader(w, UsernameAttributeType, uint16(len(unb)))
	w.Bytes(unb)
	WriteTLVPadding(w, uint16(len(unb)))

	//return TLVHeaderSize + uint16(len(unb)) + p, nil
	return nil
}

type UsernameAttribute struct {
	Username string
}

func (s UsernameAttribute) Type() AttributeType {
	return UsernameAttributeType
}
