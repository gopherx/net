package handlers

import (
	"net"

	"github.com/golang/glog"
	"github.com/gopherx/net/nat"
)

const (
	MethodBinding uint16 = 0x01
)

// BindingHandler handles a STUN Binding request.
type BindingHandler struct {
}

func (b *BindingHandler) ServeSTUN(w nat.ResponseWriter, r *nat.Request) {
	logf := glog.Infof

	ua, _ := r.Msg.Software()
	logf("[BindingHandler] ServeSTUN Software:%q", ua.Text)

	err := w.Write(NewBindingResponseMessage(r.Msg.ID, r.IP, r.Port))
	if err != nil {
		logf = glog.Errorf
	}

	logf("[BindingHandler] err:%+v", err)
}

func NewBindingRequestMessage(attrs ...nat.Attribute) nat.Message {
	return nat.NewMessage(MessageTypeBindingRequest, attrs...)
}

func NewBindingResponseMessage(rid nat.TransactionID, ip net.IP, port int) nat.Message {
	return nat.NewMessage(MessageTypeBindingResponseSuccess, nat.XorMappedAddressAttribute{ip, port})
}
