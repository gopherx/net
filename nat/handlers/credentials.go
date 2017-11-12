package handlers

import (
	"github.com/golang/glog"

	"github.com/gopherx/base/errors"
	"github.com/gopherx/net/nat"
)

type LongTermCredentialsHandler struct {
	Realm string
	Next  nat.Handler
}

func (l *LongTermCredentialsHandler) ServeSTUN(w nat.ResponseWriter, r *nat.Request) {
	auth, ok := findAuth(r.Msg)
	if !ok {
		//...first request; start handshake.
		method := r.Msg.Type.Method()
		tID := r.Msg.TID
		glog.Infof("[%s:%d] initial auth handshake (method:%d)", r.IP, r.Port, method)
		resp, err := l.newHandshakeResponse(method, tID)
		if err != nil {
			glog.Errorf("[%s:%d] failed to create auth handshake response", r.IP, r.Port)
			resp = NewServerErrorResponse(method, tID, err)
			return
		}

		w.Write(resp, nil)
		return
	}

	if v := glog.V(11); v {
		v.Infof("[%s:%d] auth ok", auth, ok)
	}

	panic("check auth!")

	l.Next.ServeSTUN(w, r)
}

func (l *LongTermCredentialsHandler) newHandshakeResponse(method uint16, tID nat.TransactionID) (nat.Message, error) {
	nonce, err := nat.NewNonceAttribute()
	if err != nil {
		return nat.Message{}, errors.Internal(err, "failed to generate Nonce")
	}

	realm := nat.RealmAttribute{l.Realm}

	return nat.NewErrorResponse(method, tID, 4, 1, "", nonce, realm), nil
}

func RequireLongTermCreds(realm string, h nat.Handler) nat.Handler {
	return &LongTermCredentialsHandler{realm, h}
}

type Auth struct {
	Username  nat.UsernameAttribute
	Integrity nat.MessageIntegrityAttribute
	Realm     nat.RealmAttribute
	Nonce     nat.NonceAttribute
}

func findAuth(msg nat.Message) (Auth, bool) {
	u, ok := msg.Username()
	if !ok {
		return Auth{}, false
	}

	i, ok := msg.MessageIntegrity()
	if !ok {
		return Auth{}, false
	}

	r, ok := msg.Realm()
	if !ok {
		return Auth{}, false
	}

	n, ok := msg.Nonce()
	if !ok {
		return Auth{}, false
	}

	return Auth{u, i, r, n}, true
}
