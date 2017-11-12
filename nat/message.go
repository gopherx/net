package nat

import (
	"fmt"
	"os"

	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/errors"

	crand "crypto/rand"
	mrand "math/rand"
)

const (
	MagicCookie = 0x2112A442

	MessageClassMask = 0x0110
	MessageTypeMask  = 0x3FFF

	MessageClassRequest         MessageClass = 0x00
	MessageClassIndication      MessageClass = 0x01
	MessageClassResponseSuccess MessageClass = 0x02
	MessageClassResponseError   MessageClass = 0x03

	// HeaderSize is the size of the STUN header.
	HeaderSize uint16 = 20

	// TLVHeaderSize is the size of a TLV (Type, Length Value) encoded attribute.
	TLVHeaderSize uint16 = 4
)

var (
	// EmptyMessage is the default empty Message instance.
	EmptyMessage Message = Message{}

	rnd *mrand.Rand
)

func init() {
	b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	n, err := crand.Reader.Read(b)
	if err != nil || n != 8 {
		fmt.Fprintf(os.Stderr, "failed to read crypto.rand.Reader; n=%d err=%+v", n, err)
		os.Exit(666)
	}

	seed := read.Int64(b)
	rnd = mrand.New(mrand.NewSource(seed))
}

// MessageType holds the STUN message type.
type MessageType uint16

// MessageTypeClass holds the STUN message class.
type MessageClass byte

func (c MessageClass) String() string {
	switch {
	case c == MessageClassRequest:
		return "REQUEST"

	case c == MessageClassIndication:
		return "INDICATION"

	case c == MessageClassResponseSuccess:
		return "RESPONSE-SUCCESS"

	case c == MessageClassResponseError:
		return "RESPONSE-ERROR"
	}

	return "UNKNOWN"
}

func (m MessageType) Method() uint16 {
	//...extract the encoded method bits
	encoded := uint16(m & 0xFEEF)

	//...lower 4 bits
	v := encoded & 0x000F
	//...mid 3 bits
	v = v | (encoded>>1)&0x0070
	//...upper 5 bits
	v = v | (encoded>>2)&0x0F80

	return v
}

func (m MessageType) Class() MessageClass {
	encoded := uint16(m & 0x0110)
	v := (encoded >> 4) & 0x0001
	v = v | (encoded>>7)&0x0002
	return MessageClass(v)
}

func NewMessageType(method uint16, class MessageClass) MessageType {
	upperM := (method << 2) & 0x3E00
	midM := (method << 1) & 0x00E0
	lowerM := method & 0x000F

	upperC := (uint16(class) << 7) & 0x0100
	lowerC := (uint16(class) << 4) & 0x0010

	return MessageType(upperM | midM | lowerM | upperC | lowerC)
}

// TransactionID is a STUN transaction ID.
type TransactionID struct {
	p0 uint32
	p1 uint32
	p2 uint32
}

// NewTransactionID returns a new TransactionID object.
func NewTransactionID() TransactionID {
	return TransactionID{rnd.Uint32(), rnd.Uint32(), rnd.Uint32()}
}

type Message struct {
	Type  MessageType
	TID   TransactionID
	Attrs map[AttributeType]Attribute
	Types []AttributeType
}

func (m Message) Software() (SoftwareAttribute, bool) {
	a, ok := m.Attrs[SoftwareAttributeType]
	if !ok {
		return SoftwareAttribute{}, false
	}

	swa, ok := a.(SoftwareAttribute)
	return swa, ok
}

func (m Message) Nonce() (NonceAttribute, bool) {
	a, ok := m.Attrs[NonceAttributeType]
	if !ok {
		return NonceAttribute{}, false
	}

	na, ok := a.(NonceAttribute)
	return na, ok
}

func (m Message) Realm() (RealmAttribute, bool) {
	a, ok := m.Attrs[RealmAttributeType]
	if !ok {
		return RealmAttribute{}, false
	}

	ra, ok := a.(RealmAttribute)
	return ra, ok
}

func (m Message) Username() (UsernameAttribute, bool) {
	a, ok := m.Attrs[UsernameAttributeType]
	if !ok {
		return UsernameAttribute{}, false
	}

	ua, ok := a.(UsernameAttribute)
	return ua, ok
}

func (m Message) MessageIntegrity() (MessageIntegrityAttribute, bool) {
	a, ok := m.Attrs[MessageIntegrityAttributeType]
	if !ok {
		return MessageIntegrityAttribute{}, false
	}

	mi, ok := a.(MessageIntegrityAttribute)
	return mi, ok
}

// NewMessage returns a new Message.
func NewMessage(method uint16, class MessageClass, attrs ...Attribute) Message {
	return MakeMessage(NewMessageType(method, class), NewTransactionID(), attrs)
}

// NewRequest creates a new request message.
func NewRequest(method uint16, attrs ...Attribute) Message {
	tID := NewTransactionID()
	return MakeMessage(NewMessageType(method, MessageClassRequest), tID, attrs)
}

// NewErrorResponse returns a new response Message with one ERROR-CODE attribute.
func NewErrorResponse(
	method uint16,
	tID TransactionID,
	class byte,
	number byte,
	reason string,
	attrs ...Attribute) Message {

	attrs = append(attrs, ErrorCodeAttribute{class, number, reason})

	return MakeMessage(NewMessageType(method, MessageClassResponseError), tID, attrs)
}

// MakeMessage is a low level convenience function with access to all details.
func MakeMessage(t MessageType, tID TransactionID, attrs []Attribute) Message {
	m := Message{t, tID, nil, nil}
	m.initAttrs(attrs)
	return m
}

func (m *Message) initAttrs(attrs []Attribute) {
	m.Attrs = map[AttributeType]Attribute{}
	for _, attr := range attrs {
		at := attr.Type()
		m.Attrs[at] = attr
		m.Types = append(m.Types, at)
	}
}

type AttributeType uint16

func (a AttributeType) String() string {
	return fmt.Sprintf("0x%x", int64(a))
}

type Attribute interface {
	Type() AttributeType
}

// Read127CharString reads a string that is at most 127 chars long (and 783 bytes)
func Read127CharString(r *read.BigEndian, l uint16) (string, error) {
	const maxBytes = 783
	const maxChars = 128

	if l > maxBytes {
		return "", errors.InvalidArgument(nil, fmt.Sprintf("too many bytes in text; max=%d current%d", maxBytes, l))
	}

	txt := string(r.Bytes(int(l)))
	if len(txt) > maxChars {
		return "", errors.InvalidArgument(nil, fmt.Sprintf("too many chars in text; max=%d current=%d", maxChars, len(txt)))
	}

	return txt, nil
}

// Check127CharString verifies that the string is at most 127 chars and 783 bytes long.
func Check127CharString(txt string) ([]byte, error) {
	const maxBytes = 783
	const maxChars = 128

	if len(txt) >= maxChars {
		return nil, errors.InvalidArgument(nil, "too many chars; max=%d current=%d", maxChars, len(txt))
	}

	bytes := []byte(txt)
	if len(bytes) > int(maxBytes) {
		return nil, errors.InvalidArgument(nil, "oo many bytes; max=%d current=%d", maxBytes, len(bytes))
	}

	return bytes, nil
}
