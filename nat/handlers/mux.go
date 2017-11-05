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
	m := r.Msg.Type.Method()
	h, ok := m.handlers[m]
	v := glog.V(11)
	if !ok {
		v.Info("handler missing; method: ", m)
		return
	}

	h.ServeSTUN(w, r)
}

func (m *MuxSTUN) Add(t nat.MessageType, h nat.Handler) {
	m.handlers[t] = h
}
