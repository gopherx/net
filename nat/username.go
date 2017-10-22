package nat

import (
	"github.com/gopherx/base/errors"
)

const (
	UsernameAttributeType     AttributeType = 0x0006
	UsernameAttributeRfcName  string        = "USERNAME"
	UsernameAttributeMaxBytes int           = 512
)

func init() {
	RegisterUsernameAttribute(DefaultParser)
}

func RegisterUsernameAttribute(p *MessageParser) {
	p.Register(
		UsernameAttributeType,
		UsernameAttributeRfcName,
		func(b []byte) (Attribute, error) {
			return ParseUsernameAttribute(b)
		})
}

func ParseUsernameAttribute(b []byte) (UsernameAttribute, error) {
	res := UsernameAttribute{}

	if len(b) > UsernameAttributeMaxBytes {
		return res, errors.InvalidArgument(nil, "too many bytes", len(b), UsernameAttributeMaxBytes)
	}

	res.Username = string(b)
	return res, nil
}

type UsernameAttribute struct {
	Username string
}
