package nat

import (
	"net"
)

type MappedAddressAttribute struct {
	Address net.IP
	Port    uint16
}

type XorMappedAddressAttribute struct {
	MappedAddressAttribute
}
