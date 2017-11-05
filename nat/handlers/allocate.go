package handlers

import (
	"net"

	"github.com/golang/glog"
	"github.com/gopherx/net/nat"
)

const (
	MessageTypeAllocateRequest = nat.MessageType(0x0003)
)

// AllocateHandler handles a STUN Allocate request.
type AllocateHandler struct {
}

func (b *AllocateHandler) ServeSTUN(w nat.ResponseWriter, r *nat.Request) {
}
