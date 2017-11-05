package nat

import (
	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/binary/write"
)

const (
	RequestedTransportAttributeType    AttributeType = 0x0019
	RequestedTransportAttributeRfcName string        = "REQUESTED-TRANSPORT"
	RequestedTransportAttributeSize    uint16        = 4
)

func init() {
	RegisterRequestedTransportAttribute(DefaultParser)
}

func RegisterRequestedTransportAttribute(p *MessageParser) {
	p.Register(
		RequestedTransportAttributeType,
		RequestedTransportAttributeRfcName,
		func(r *read.BigEndian, l uint16) (Attribute, error) {
			return ParseRequestedTransportAttribute(r, l)
		},
		func(w *write.BigEndian, a Attribute) error {
			return PrintRequestedTransportAttribute(w, a.(RequestedTransportAttribute))
		},
	)
}

func ParseRequestedTransportAttribute(r *read.BigEndian, l uint16) (RequestedTransportAttribute, error) {
	p := r.Byte()
	//...read the RFFU fields but ignore them since they carry no information.
	r.Byte()
	r.Byte()
	r.Byte()
	return RequestedTransportAttribute{p}, nil
}

func PrintRequestedTransportAttribute(w *write.BigEndian, a RequestedTransportAttribute) error {
	WriteTLVHeader(w, RequestedTransportAttributeType, RequestedTransportAttributeSize)

	w.Byte(a.Protocol)
	//...write the RFFU fields as zero for now
	w.Byte(0x00)
	w.Uint16(0x0000)
	return nil
}

type RequestedTransportAttribute struct {
	Protocol byte
}

func (r RequestedTransportAttribute) Type() AttributeType {
	return RequestedTransportAttributeType
}
