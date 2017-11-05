package handlers

import (
	"github.com/golang/glog"
	"github.com/gopherx/net/nat"
)

type LongTermCredentialsHandler struct {
	Next nat.Handler
}

func (l *LongTermCredentialsHandler) ServeSTUN(w nat.ResponseWriter, r *nat.Request) {
	auth, ok := findAuth(r.Msg)
	if !ok {
		//...first request; start handshake.
		glog.Infof("[%s:%d] initial auth handshake", r.IP, r.Port)
		return
	}

	glog.Info(auth, ok)

	l.Next.ServeSTUN(w, r)
}

func RequireLongTermCreds(h nat.Handler) nat.Handler {
	return &LongTermCredentialsHandler{h}
}

type auth struct {
	Username  nat.UsernameAttribute
	Integrity nat.MessageIntegrityAttribute
	Realm     nat.RealmAttribute
	Nonce     nat.NonceAttribute
}

func findAuth(msg nat.Message) (auth, bool) {
	u, ok := msg.Attrs[nat.UsernameAttributeType]
	if !ok {
		return auth{}, false
	}

	i, ok := msg.Attrs[nat.MessageIntegrityAttributeType]
	if !ok {
		return auth{}, false
	}

	r, ok := msg.Attrs[nat.RealmAttributeType]
	if !ok {
		return auth{}, false
	}

	n, ok := msg.Attrs[nat.NonceAttributeType]
	if !ok {
		return auth{}, false
	}

	return auth{
		u.(nat.UsernameAttribute),
		i.(nat.MessageIntegrityAttribute),
		r.(nat.RealmAttribute),
		n.(nat.NonceAttribute),
	}, true
}
