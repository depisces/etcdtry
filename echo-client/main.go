package main

import (
	"context"
	"etcdtry/echo"
	"etcdtry/etcd"
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
// addr = flag.String("addr", "localhost:50051", "")
)

func main() {
	flag.Parse()
	etcd.CusLoadService("echo-service")
	addr := etcd.CusServiceDiscovery("echo-service")

	fmt.Println(addr)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := echo.NewEchoClient(conn)
	CallUnaryEcho(c)
}

func CallUnaryEcho(c echo.EchoClient) {
	ctx := context.Background()
	in := echo.EchoMessage{
		Mess: "client say hello",
	}
	res, err := c.UnaryEcho(ctx, &in)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("client recv", res.Mess)
}
