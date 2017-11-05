package handlers

import (
	"net"
	"testing"

	"github.com/gopherx/net/nat"
	"github.com/gopherx/net/nat/handlers/test"
)

func TestBindingHandler(t *testing.T) {
	req := NewBindingRequestMessage(newSoftwareAttribute(t))
	h := &BindingHandler{}
	writer := &test.Writer{}
	h.ServeSTUN(writer, &nat.Request{req, net.ParseIP("192.168.0.1"), 123234, ""})

}
