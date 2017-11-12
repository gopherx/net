package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/golang/glog"

	"github.com/gopherx/net/nat"

	nath "github.com/gopherx/net/nat/handlers"
)

var (
	realm = flag.String("realm", "example.com", "The realm to use")
)

func main() {
	flag.Parse()

	glog.Info("hello turnd")

	handler := nath.Mux()
	handler.Add(nath.MethodBinding, &nath.BindingHandler{})
	handler.Add(nath.MethodAllocate, nath.RequireLongTermCreds(*realm, &nath.AllocateHandler{}))

	parser := &nat.MessageParser{nat.DefaultRegistry}
	printer := &nat.MessagePrinter{nat.DefaultRegistry, 1024}

	udpsrv := &nat.UDPServer{
		nil,
		parser,
		printer,
		handler,
	}

	err := udpsrv.Start()
	if err != nil {
		glog.Error("failed to start UDPServer; err:", err)
	}

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt)

	glog.Info("ctrl-c to terminate")
	<-sigchan
	glog.Info("terminating")
}
