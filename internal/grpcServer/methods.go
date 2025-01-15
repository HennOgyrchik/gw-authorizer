package grpcServer

import (
	"context"
	"errors"
	pb "github.com/HennOgyrchik/proto-jwt-auth/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *handler) CreateUser(ctx context.Context, user *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	userID, err := h.auth.CreateUser(ctx, user)

	switch {
	case errors.Is(err, UserAlreadyExistsErr):
		return userID, status.Errorf(codes.AlreadyExists, err.Error())
	case err != nil:
		return userID, status.Errorf(codes.Internal, "internal error")
	default:
		return userID, nil
	}
}

func (h *handler) Login(ctx context.Context, credentials *pb.LoginRequest) (*pb.Token, error) {
	token, err := h.auth.Login(ctx, credentials)

	switch {
	case errors.Is(err, InvalidCredentialsErr):
		return token, status.Errorf(codes.InvalidArgument, err.Error())
	case err != nil:
		return token, status.Errorf(codes.Internal, "internal error")
	default:
		return token, nil
	}
}

func (h *handler) VerifyToken(ctx context.Context, token *pb.Token) (*pb.VerifyTokenResponse, error) {
	ok, err := h.auth.VerifyToken(ctx, token)
	if err != nil {
		err = status.Errorf(codes.Internal, "internal error")
	}

	return ok, err
}
