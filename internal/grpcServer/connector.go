package grpcServer

import (
	"fmt"
	pb "github.com/HennOgyrchik/proto-jwt-auth/auth"
	"google.golang.org/grpc"
	"net"
	"time"
)

func New(addr string, timeout time.Duration, hndlr AuthorizationServer) *Authorizer {
	opts := []grpc.ServerOption{grpc.ConnectionTimeout(timeout)}
	grpcSrv := grpc.NewServer(opts...)

	pb.RegisterAuthorizationServer(grpcSrv, &handler{auth: hndlr})

	return &Authorizer{
		server: grpcSrv,
		addr:   addr,
	}
}

func (a *Authorizer) Run() error {
	const op = "gRPC run"

	listener, err := net.Listen("tcp", a.addr)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = a.server.Serve(listener)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *Authorizer) Stop() {
	a.server.GracefulStop()
}
