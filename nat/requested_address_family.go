package nat

import (
	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/binary/write"
)

const (
	RequestedAddressFamiliyAttributeType    AttributeType = 0x0017
	RequestedAddressFamiliyAttributeRfcName string        = "REQUESTED-ADDRESS-FAMILY"
	RequestedAddressFamiliyAttributeSize    uint16        = 4
	RequestedAddressFamiliyAttributeIPv4    byte          = 0x01
	RequestedAddressFamiliyAttributeIPv6    byte          = 0x02
)

func init() {
	RegisterRequestedAddressFamiliyAttribute(DefaultParser)
}

func RegisterRequestedAddressFamiliyAttribute(p *MessageParser) {
	p.Register(
		RequestedAddressFamiliyAttributeType,
		RequestedAddressFamiliyAttributeRfcName,
		func(r *read.BigEndian, l uint16) (Attribute, error) {
			return ParseRequestedAddressFamiliyAttribute(r, l)
		},
		func(w *write.BigEndian, a Attribute) error {
			return PrintRequestedAddressFamiliyAttribute(w, a.(RequestedAddressFamiliyAttribute))
		},
	)
}

func ParseRequestedAddressFamiliyAttribute(r *read.BigEndian, l uint16) (RequestedAddressFamiliyAttribute, error) {
	v := r.Byte()
	r.Byte()
	r.Uint16()
	return RequestedAddressFamiliyAttribute{v}, nil
}

func PrintRequestedAddressFamiliyAttribute(w *write.BigEndian, a RequestedAddressFamiliyAttribute) error {
	WriteTLVHeader(w, RequestedAddressFamiliyAttributeType, RequestedAddressFamiliyAttributeSize)
	w.Byte(a.Family)
	//...can't use padding to fill due to spec requirement of using zeros
	w.Byte(0x00)
	w.Uint16(0x0000)
	return nil
}

type RequestedAddressFamiliyAttribute struct {
	Family byte
}

func (r RequestedAddressFamiliyAttribute) Type() AttributeType {
	return RequestedAddressFamiliyAttributeType
}
