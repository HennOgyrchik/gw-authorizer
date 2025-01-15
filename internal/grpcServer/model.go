package grpcServer

import (
	"context"
	"fmt"
	pb "github.com/HennOgyrchik/proto-jwt-auth/auth"
	"google.golang.org/grpc"
)

var UserAlreadyExistsErr = fmt.Errorf("already exists")
var InvalidCredentialsErr = fmt.Errorf("invalid credentials")

type AuthorizationServer interface {
	CreateUser(ctx context.Context, user *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	Login(ctx context.Context, credentials *pb.LoginRequest) (*pb.TokenResponse, error)
	VerifyToken(ctx context.Context, token *pb.TokenReuest) (*pb.VerifyTokenResponse, error)
}

type Authorizer struct {
	server *grpc.Server
	addr   string
}

type handler struct {
	pb.UnimplementedAuthorizationServer
	auth AuthorizationServer
}
