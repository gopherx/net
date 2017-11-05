package handlers

import (
	"github.com/golang/glog"
	"github.com/gopherx/net/nat"
)

const (
	MethodAllocate uint16 = 0x03
)

// AllocateHandler handles a STUN Allocate request.
type AllocateHandler struct {
}

func (b *AllocateHandler) ServeSTUN(w nat.ResponseWriter, r *nat.Request) {
	glog.Info("allocate")
}
