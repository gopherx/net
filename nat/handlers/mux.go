package handlers

import (
	"github.com/golang/glog"

	"github.com/gopherx/net/nat"
)

func Mux() *MuxSTUN {
	return &MuxSTUN{map[uint16]nat.Handler{}}
}

type MuxSTUN struct {
	handlers map[uint16]nat.Handler
}

func (m *MuxSTUN) ServeSTUN(w nat.ResponseWriter, r *nat.Request) {
	method := r.Msg.Type.Method()
	h, ok := m.handlers[method]
	v := glog.V(11)
	if !ok {
		v.Info("handler missing; method: ", method)
		return
	}

	h.ServeSTUN(w, r)
}

func (m *MuxSTUN) Add(method uint16, h nat.Handler) {
	m.handlers[method] = h
}
