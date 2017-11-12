package handlers

import (
	"github.com/gopherx/net/nat"
)

const (
	MethodAllocate uint16 = 0x03
)

type AllocateHandler struct {
}

func (h *AllocateHandler) ServeSTUN(w nat.ResponseWriter, r *nat.Request) {
}
