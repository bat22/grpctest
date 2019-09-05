package main

import (
	"context"
	"github.com/bat22/grpctest/internal/rpc"
	"google.golang.org/grpc"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":55000")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	rpc.RegisterTestSvcServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

type server struct {
}

func (i *server) Add(ctx context.Context, in *rpc.AddIntMessage) (*rpc.AddIntReply, error) {
	return &rpc.AddIntReply{
		C: in.A + in.B,
	}, nil
}
