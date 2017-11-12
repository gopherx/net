package handlers

import (
	"testing"

	"github.com/gopherx/net/nat"
	"github.com/gopherx/net/nat/handlers/test"
)

func TestLongTermCredentialsHandshake(t *testing.T) {
	th := &testHandler{}
	ch := &LongTermCredentialsHandler{"unittest.com", th}

	writer := &test.Writer{t, nil}

	allocate := nat.NewRequest(MethodAllocate)

	ch.ServeSTUN(writer, &nat.Request{allocate, test.IP, test.Port, test.Zone})
	resp := writer.Pop()
	nonce, ok := resp.Nonce()
	if !ok || len(nonce.Nonce) != 127 {
		t.Fatal("nonce missing or invalid", ok, nonce)
	}

	realm, ok := resp.Realm()
	if !ok || realm.Realm != "unittest.com" {
		t.Fatal("realm missing of invalid", ok, realm)
	}
}

type testHandler struct {
}

func (t *testHandler) ServeSTUN(w nat.ResponseWriter, r *nat.Request) {
}
