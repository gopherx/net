package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/golang/glog"

	"github.com/gopherx/net/nat"

	nath "github.com/gopherx/net/nat/handlers"
)

func main() {
	flag.Parse()

	glog.Info("hello turnd")

	handler := nath.Mux()
	handler.Add(nath.MethodBinding, &nath.BindingHandler{})
	handler.Add(nath.MethodAllocate, nath.RequireLongTermCreds(&nath.AllocateHandler{}))

	parser := &nat.MessageParser{nat.DefaultRegistry}

	udpsrv := &nat.UDPServer{
		nil,
		parser,
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
