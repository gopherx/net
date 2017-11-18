package nat

import (
	"bytes"
	"io"

	"github.com/golang/glog"

	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/errors"
)

var (
	DefaultParser *MessageParser = &MessageParser{
		DefaultRegistry,
	}
)

const (
	initialAttrCount = 16
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

func (p *MessageParser) parseHeader(r *read.BigEndian) (MessageType, uint16, TransactionID, error) {
	tid := TransactionID{}
	mt := MessageType(r.Uint16())
	if 0xC0&mt == 0xC0 {
		//...first two bits not zero.
		return mt, 0, tid, errors.InvalidArgument(nil, "first bits not zero", 0xC0&mt)
	}

	ml := r.Uint16()
	cookie := r.Uint32()
	if cookie != MagicCookie {
		return mt, 0, tid, errors.InvalidArgument(nil, "invalid cookie", cookie)
	}

	tid.p0, tid.p1, tid.p2 = r.Uint32x3()

	if r.Err != nil {
		return mt, 0, tid, errors.InvalidArgument(r.Err, "reader failure")
	}

	return mt, ml, tid, nil
}

func (p *MessageParser) parseAttribute(at AttributeType, al uint16, r *read.BigEndian) (Attribute, uint16, error) {
	reg, ok := p.Registry[at]
	name := ""
	var parse AttributeParserFunc
	if ok {
		name = reg.Name
		parse = reg.Parse
	} else {
		parse = func(r *read.BigEndian, l uint16) (Attribute, error) {
			return UnknownAttribute{at, r.Bytes(int(l))}, nil
		}
	}

	if v := glog.V(11); v {
		v.Infof("attribute: [%v] %s (size: %d)", at, name, al)
	}

	r0 := len(r.Read)
	a, err := parse(r, al)
	if err != nil {
		return nil, 0, errors.InvalidArgument(err, "failed to parse attribute", at, name)
	}

	readBytes := len(r.Read) - r0
	if uint16(readBytes) != al {
		return nil, 0, errors.InvalidArgument(nil, "read does not match length", readBytes, al)
	}

	if r.Err != nil {
		return nil, 0, r.Err
	}

	//...finally remove any padding!
	padding := uint16(0)
	if readBytes%4 > 0 {
		pad := (4 - readBytes%4)
		for i := 0; i < pad; i++ {
			r.Byte()
			padding++
		}
	}

	return a, padding, nil
}

func (p *MessageParser) parse(r *read.BigEndian) (Message, error) {
	mt, ml, tid, err := p.parseHeader(r)
	if err != nil {
		return EmptyMessage, err
	}

	msg := Message{mt, tid, nil}
	remaining := ml
	fpFound := false
	miaFound := false

	for remaining > 0 {
		//...read the TLV encoded attribute (Type, Length, Value)
		at := AttributeType(r.Uint16())

		// TODO(d):
		// - fingerprint must be last; fail parse if not last
		// - mi may be last; ignore subsequent bytes if last is not fingerprint

		// If we already found a MESSAGE-INTEGRITY attribute then the current attribute must be a FINGERPRINT;
		// if current attrbute is not a FINGERPRINT then all other attributes must be ignored.
		// If we already found a FINGERPRINT then all other attributes must be ignored.
		ignoreRemaining := miaFound && at != FingerprintAttributeType
		if ignoreRemaining {
			remaining -= 2 // adjust the remaining for the attribute type just read
			if v := glog.V(11); v {
				v.Infof("ignoring 0x%x bytes (MESSAGE-INTEGRITY:%v FINGERPRINT:%v)", remaining, miaFound, fpFound)
			}

			for i := uint16(0); i < remaining; i++ {
				r.Byte()
			}

			break
		}

		if fpFound {
			return EmptyMessage, errors.InvalidArgument(nil, "FINGERPRINT must be the last attribute", remaining)
		}

		al := r.Uint16()
		remaining -= TLVHeaderSize

		a, padding, err := p.parseAttribute(at, al, r)
		if err != nil {
			return EmptyMessage, err
		}

		remaining -= al
		remaining -= padding

		if a.Type() != at {
			panic("parsed attribute don't match wire format!!!")
		}

		// MESSAGE-INTEGRITY (MI) and FINGERPRINT (FP) must be handled with care! They must:
		// - Both MI and FP are optional
		// - MI must be last or second to last. If second to last then FP must be last
		// - FP must be last
		// - All attributes after MI or FP must be ignored
		switch {
		case at == FingerprintAttributeType:
			fpFound = true

		case at == MessageIntegrityAttributeType:
			miaFound = true
		}

		msg.Attrs = append(msg.Attrs, AttrRef{at, a})
	}

	return msg, nil
}

func (p *MessageParser) Register(t AttributeType, name string, parse AttributeParserFunc, print AttributePrinterFunc) {
	p.Registry[t] = Registration{name, parse}
}

// ParseMessage parses a byte array into a Message object.
func ParseMessage(b []byte) (Message, error) {
	return DefaultParser.Parse(b)
}
