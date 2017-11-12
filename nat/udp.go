package nat

import (
	"net"

	"github.com/golang/glog"
)

var (
	DefaultUDPAddr = &net.UDPAddr{net.ParseIP("127.0.0.1"), DefaultTURNPort, ""}
)

type UDPServer struct {
	Addr    *net.UDPAddr
	Parser  *MessageParser
	Printer *MessagePrinter
	Handler Handler
}

func (u *UDPServer) Start() error {
	if u.Addr == nil {
		glog.Info("address not set; using default")
		u.Addr = DefaultUDPAddr
	}

	glog.Info("listen udp; ", u.Addr)
	conn, err := net.ListenUDP("udp4", u.Addr)
	if err != nil {
		return err
	}

	go u.read(conn)

	return nil
}

func (u *UDPServer) read(conn *net.UDPConn) {
	glog.Info("dispatching;", conn.LocalAddr())

	for {
		buff := make([]byte, 1500)
		n, src, err := conn.ReadFromUDP(buff)
		if err != nil {
			glog.Error("read failed; err: ", err)
			continue
		}

		go u.dispatch(conn, buff[0:n], src)
	}
}

type writer struct {
	conn *net.UDPConn
	u    *UDPServer
	src  *net.UDPAddr
}

func (w writer) Write(msg Message, opts *PrintOptions) error {
	bytes, err := w.u.Printer.Print(msg, opts)
	if err != nil {
		return err
	}

	if v := glog.V(11); v {
		v.Infof("[%v] %v (opts:%v err:%v)", w.src, msg, opts, err)
	}

	_, err = w.conn.WriteToUDP(bytes, w.src)
	return err
}

func (u *UDPServer) dispatch(conn *net.UDPConn, buff []byte, src *net.UDPAddr) {
	m, err := u.Parser.Parse(buff)
	if v := glog.V(11); v {
		v.Infof("[%v] %v (err:%v)", src, m, err)
	}

	u.Handler.ServeSTUN(writer{conn, u, src}, &Request{m, src.IP, src.Port, src.Zone})
}
