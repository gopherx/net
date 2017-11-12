package handlers

import (
	"testing"

	"github.com/gopherx/net/nat"
	"github.com/gopherx/net/nat/handlers/test"
)

func TestBindingHandler(t *testing.T) {
	req := NewBindingRequestMessage(test.SwAttr)
	h := &BindingHandler{}
	writer := &test.Writer{}
	h.ServeSTUN(writer, &nat.Request{req, test.IP, test.Port, test.Zone})
}
