package main

import (
	"etcdtry/echo"
	"etcdtry/echo-server/server"
	"etcdtry/etcd"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	echo.RegisterEchoServer(s, &server.EchoService{})

	etcd.CusServiceRegister("echo-service", fmt.Sprintf("localhost:%d", *port))

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
