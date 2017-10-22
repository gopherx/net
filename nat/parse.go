package nat

import (
	"github.com/gopherx/base/errors"
	"github.com/gopherx/base/read"
)

const (
	// HeaderSize is the size of the STUN header.
	HeaderSize int = 20
)

var (
	DefaultParser *MessageParser = &MessageParser{map[AttributeType]registration{}}

	unknownAttributeRegistration registration = registration{
		"UNKNOWN",
		func(b []byte) (Attribute, error) {
			return UnknownAttribute{b}, nil
		},
	}
)

// AttributeParser parses bytes into Attribute instances.
type AttributeParserFunc func(b []byte) (Attribute, error)

// MessageParser is used to parse bytes into STUN messages.
type MessageParser struct {
	registry map[AttributeType]registration
}

// Parse parses bytes into a STUN message.
func (p *MessageParser) Parse(b []byte) (Message, error) {
	if len(b) < HeaderSize {
		return EmptyMessage, errors.InvalidArgument(nil, "buffer too small; size:", b)
	}

	mt := MessageType(read.Uint16(b))
	if 0xC0&mt == 0xC0 {
		//...first two bits not zero.
		return EmptyMessage, errors.InvalidArgument(nil, "first bits not zero", 0xC0&mt)
	}

	ml := read.Uint16(b[2:4])
	cookie := read.Uint32(b[4:8])
	if cookie != MagicCookie {
		return EmptyMessage, errors.InvalidArgument(nil, "invalid cookie", cookie)
	}

	p0, p1, p2 := read.Uint32x3(b[8:20])
	tid := TransactionID{p0, p1, p2}

	remd := b[20:]
	var attrs []Attribute = nil

	msg := Message{mt, tid, ml, attrs}

	for len(remd) > 0 {
		//...read the TLV encoded attribute (Type, Length, Value)
		at := AttributeType(read.Uint16(remd))
		al := read.Uint16(remd[2:4])
		ad := remd[4 : 4+al]

		end := 4 + al
		//...data is always be padded to a 32-bit boundry
		if end%4 > 0 {
			end += (4 - end%4)
		}

		remd = remd[end:]

		reg, ok := p.registry[at]
		if !ok {
			reg = unknownAttributeRegistration
		}

		a, err := reg.parse(ad)
		if err != nil {
			return EmptyMessage, errors.InvalidArgument(err, "failed to parse attribute", at, reg.name)
		}

		msg.Attrs = append(msg.Attrs, a)
	}

	return msg, nil
}

func (p *MessageParser) Register(t AttributeType, name string, parse AttributeParserFunc) {
	p.registry[t] = registration{name, parse}
}

type registration struct {
	name  string
	parse AttributeParserFunc
}

// ParseMessage parses a byte array into a Message object.
func ParseMessage(b []byte) (Message, error) {
	return DefaultParser.Parse(b)
}

type UnknownAttribute struct {
	Data []byte
}
