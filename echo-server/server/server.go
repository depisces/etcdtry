package server

import (
	"context"
	"etcdtry/echo"
	"fmt"
)

type EchoService struct {
	echo.UnimplementedEchoServer
}

func (EchoService) UnaryEcho(ctx context.Context, in *echo.EchoMessage) (*echo.EchoMessage, error) {
	fmt.Println("server recv", in.Mess)
	return &echo.EchoMessage{
		Mess: "server send " + "hello client",
	}, nil
}
