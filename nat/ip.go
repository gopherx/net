package nat

import (
	"net"

	"github.com/gopherx/base/binary/write"
)

const (
	XorMappedAddressAttributeType    AttributeType = 0x0020
	XorMappedAddressAttributeRfcName string        = "XOR-MAPPED-ADDRESS"
)

type MappedAddressAttribute struct {
	Address net.IP
	Port    int
}

type XorMappedAddressAttribute struct {
	Address net.IP
	Port    int
}

func (x XorMappedAddressAttribute) Type() AttributeType {
	return XorMappedAddressAttributeType
}

func (x XorMappedAddressAttribute) Print(w *write.BigEndian) error {
	panic("implement me!!!!")
}
