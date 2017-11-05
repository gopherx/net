package nat

import (
	"net"
)

const (
	DefaultTURNPort = 3478
)

// Handler to handle STUN requests and indications.
type Handler interface {
	ServeSTUN(ResponseWriter, *Request)
}

// Request is a STUN request or indication.
type Request struct {
	Msg  Message
	IP   net.IP
	Port int
	Zone string
}

// ResponseWriter is used to respond to requests.
type ResponseWriter interface {
	Write(m Message) error
}
