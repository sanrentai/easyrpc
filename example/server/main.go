package main

import (
	"net"
	"time"

	"github.com/sanrentai/easyrpc"
	"github.com/sanrentai/easyrpc/registry"
)

type Foo int

type Args struct{ Num1, Num2 int }

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func (f Foo) Sleep(args Args, reply *int) error {
	time.Sleep(time.Second * time.Duration(args.Num1))
	*reply = args.Num1 + args.Num2
	return nil
}

func main() {
	registryAddr := "http://localhost:9999/_easyrpc_/registry"
	var foo Foo
	l, _ := net.Listen("tcp", ":0")
	server := easyrpc.NewServer()
	_ = server.Register(&foo)
	registry.Heartbeat(registryAddr, "tcp@"+l.Addr().String(), 0)
	server.Accept(l)
}
