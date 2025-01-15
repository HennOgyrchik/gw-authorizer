package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"gw-authorizer/internal/storages"
)

func (p *PSQL) CreateUser(ctx context.Context, user storages.NewUser) (int, error) {
	const op = "PSQL CreateUser"

	ctxWithTimeout, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var userID int
	err := p.pool.QueryRow(ctxWithTimeout,
		"insert into users(login, password, email) values ($1,$2,$3) returning id",
		user.Login,
		user.Password,
		user.Email).Scan(&userID)

	var pgErr *pgconn.PgError
	switch {
	case errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation:
		err = storages.DuplicateRowError
	case err != nil:
		err = fmt.Errorf(op, err)
	}

	return userID, err
}

func (p *PSQL) GetUserInfo(ctx context.Context, login string) (storages.UserInfo, error) {
	const op = "PSQL GetUserInfo"

	ctxWithTimeout, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var user storages.UserInfo
	err := p.pool.QueryRow(ctxWithTimeout, "select id, password from users where login = $1", login).Scan(&user.ID, &user.Password)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return storages.UserInfo{}, storages.NoRowError
	case err != nil:
		return storages.UserInfo{}, fmt.Errorf(op, err)
	default:
		return user, nil
	}

}

func (p *PSQL) UpdateToken(ctx context.Context, token storages.UserToken) error {
	const op = "PSQL UpdateToken"

	ctxWithTimeout, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	_, err := p.pool.Exec(ctxWithTimeout, "update users set token = $1 where login = $2", token.Token, token.Login)
	if err != nil {
		err = fmt.Errorf(op, err)
	}

	return err
}
