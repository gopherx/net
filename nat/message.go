package nat

const (
	MagicCookie = 0x2112A442
)

var (
	// EmptyMessage is the default empty Message instance.
	EmptyMessage Message = Message{}
)

// MessageType holds the STUN message type.
type MessageType uint16

// TransactionID is a STUN transaction ID.
type TransactionID struct {
	p0 uint32
	p1 uint32
	p2 uint32
}

type Message struct {
	Type   MessageType
	ID     TransactionID
	length uint16
	Attrs  []Attribute
}

type AttributeType uint16

type Attribute interface {
}
