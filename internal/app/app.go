package app

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/HennOgyrchik/proto-jwt-auth/auth"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"gw-authorizer/internal/grpcServer"
	"gw-authorizer/internal/storages"
	"gw-authorizer/pkg/logs"
	"strconv"
)

func New(logger *logs.Log, storage storages.Storage, costEncoding int, secretKey string) *App {
	return &App{
		log:       logger,
		storage:   storage,
		cost:      costEncoding,
		secretKey: []byte(secretKey),
	}
}

func (a *App) CreateUser(ctx context.Context, user *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	const op = "App CreateUser"

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), a.cost)
	if err != nil {
		a.log.Err(op, err)
		return nil, err
	}

	userID, err := a.storage.CreateUser(ctx, storages.NewUser{
		Login:    user.Username,
		Password: string(hashPassword),
		Email:    user.Email,
	})

	switch {
	case errors.Is(err, storages.DuplicateRowError):
		return nil, grpcServer.UserAlreadyExistsErr
	case err != nil:
		return nil, fmt.Errorf(op, err)
	default:
		return &pb.CreateUserResponse{UserId: strconv.Itoa(userID)}, nil
	}

}

func (a *App) Login(ctx context.Context, credentials *pb.LoginRequest) (*pb.TokenResponse, error) {
	const op = "App Login"

	userInfo, err := a.storage.GetUserInfo(ctx, credentials.Username)
	switch {
	case errors.Is(err, storages.NoRowError):
		return nil, grpcServer.InvalidCredentialsErr
	case err != nil:
		a.log.Err(op, err)
		return nil, fmt.Errorf(op, err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(credentials.Password)); err != nil {
		return nil, grpcServer.InvalidCredentialsErr
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"id": userInfo.ID})
	signedToken, err := token.SignedString(a.secretKey)
	if err != nil {
		a.log.Err(op, err)
		return nil, fmt.Errorf(op, err)
	}

	err = a.storage.UpdateToken(ctx, storages.UserToken{
		Login: credentials.Username,
		Token: signedToken,
	})
	if err != nil {
		a.log.Err(op, err)
		return nil, fmt.Errorf(op, err)
	}

	return &pb.TokenResponse{Value: signedToken}, nil
}

func (a *App) VerifyToken(ctx context.Context, request *pb.TokenReuest) (*pb.VerifyTokenResponse, error) {
	const op = "App VerifyToken"

	id, err := strconv.Atoi(request.UserId)
	if err != nil {
		return nil, grpcServer.InvalidCredentialsErr
	}

	token, err := a.storage.GetToken(ctx, id)

	switch {
	case errors.Is(err, storages.NoRowError):
		return nil, grpcServer.InvalidCredentialsErr
	case err != nil:
		a.log.Err(op, err)
		return nil, fmt.Errorf(op, err)
	}

	var response pb.VerifyTokenResponse

	if response.Ok = request.Token == token; !response.Ok {
		return &response, grpcServer.InvalidCredentialsErr
	}

	return &response, nil
}
