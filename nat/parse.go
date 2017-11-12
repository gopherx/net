package nat

import (
	"bytes"
	"io"

	"github.com/golang/glog"

	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/errors"
)

var (
	DefaultParser *MessageParser = &MessageParser{DefaultRegistry}

	unknownAttributeRegistration = Registration{
		"UNKNOWN",
		func(r *read.BigEndian, l uint16) (Attribute, error) {
			return UnknownAttribute{r.Bytes(int(l))}, nil
		},
		nil,
	}
)

// AttributeParserFunc parses bytes into Attribute instances.
type AttributeParserFunc func(r *read.BigEndian, l uint16) (Attribute, error)

// MessageParser is used to parse bytes into STUN messages.
type MessageParser struct {
	Registry AttributeRegistry
}

// Parse parses bytes into a STUN message.
func (p *MessageParser) Parse(b []byte) (Message, error) {
	return p.ParseFrom(bytes.NewBuffer(b))
}

// ParseFrom parses a Message from the bytes.
func (p *MessageParser) ParseFrom(r io.Reader) (Message, error) {
	return p.parse(read.NewBigEndian(r))
}

func (p *MessageParser) parse(r *read.BigEndian) (Message, error) {
	mt := MessageType(r.Uint16())
	if 0xC0&mt == 0xC0 {
		//...first two bits not zero.
		return EmptyMessage, errors.InvalidArgument(nil, "first bits not zero", 0xC0&mt)
	}

	ml := r.Uint16()
	cookie := r.Uint32()
	if cookie != MagicCookie {
		return EmptyMessage, errors.InvalidArgument(nil, "invalid cookie", cookie)
	}

	p0, p1, p2 := r.Uint32x3()
	tid := TransactionID{p0, p1, p2}

	if r.Err != nil {
		return EmptyMessage, errors.InvalidArgument(r.Err, "reader failure")
	}

	msg := Message{mt, tid, map[AttributeType]Attribute{}, nil}
	remaining := ml

	for remaining > 0 {
		//...read the TLV encoded attribute (Type, Length, Value)
		at := AttributeType(r.Uint16())
		al := r.Uint16()
		remaining -= TLVHeaderSize

		reg, ok := p.Registry[at]
		if !ok {
			reg = unknownAttributeRegistration
		}

		if v := glog.V(11); v {
			v.Infof("attribute: [%v] %s (size: %d)", at, reg.Name, al)
		}

		r.ReadBytes = 0
		a, err := reg.Parse(r, al)
		if err != nil {
			return EmptyMessage, errors.InvalidArgument(err, "failed to parse attribute", at, reg.Name)
		}

		if uint16(r.ReadBytes) != al {
			return EmptyMessage, errors.InvalidArgument(nil, "read does not match length", r.ReadBytes, al)
		}

		if r.Err != nil {
			return EmptyMessage, r.Err
		}

		remaining -= al

		if r.ReadBytes%4 > 0 {
			pad := (4 - r.ReadBytes%4)
			for i := 0; i < pad; i++ {
				r.Byte()
				remaining--
			}
		}

		msg.Attrs[at] = a
		msg.Types = append(msg.Types, at)
	}

	return msg, nil
}

func (p *MessageParser) Register(t AttributeType, name string, parse AttributeParserFunc, print AttributePrinterFunc) {
	p.Registry[t] = Registration{name, parse, print}
}

// ParseMessage parses a byte array into a Message object.
func ParseMessage(b []byte) (Message, error) {
	return DefaultParser.Parse(b)
}

type UnknownAttribute struct {
	Data []byte
}

func (u UnknownAttribute) Type() AttributeType {
	panic("never ever call this")
}
