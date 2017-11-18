package test

import (
	"net"

	"github.com/gopherx/net/nat"
)

var (
	IP   = net.ParseIP("192.168.0.1")
	Port = 123234
	Zone = "unknown"

	SwAttr = nat.SoftwareAttribute{"unit-test"}

	Method = uint16(0x123)

	TID = nat.NewTransactionID()
)
