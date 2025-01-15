package storages

import (
	"context"
	"fmt"
)

var DuplicateRowError = fmt.Errorf("duplicate key value violates unique constraint")
var NoRowError = fmt.Errorf("no rows in result set")

type Storage interface {
	CreateUser(ctx context.Context, user NewUser) (int, error)
	GetUserInfo(ctx context.Context, login string) (UserInfo, error)
	GetToken(ctx context.Context, userID int) (string, error)
	UpdateToken(ctx context.Context, token UserToken) error
}

type NewUser struct {
	Login    string
	Password string
	Email    string
}

type UserInfo struct {
	ID       int
	Password string
}

type UserToken struct {
	Login string
	Token string
}
