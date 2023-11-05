package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AuthPostgres struct {
	db *pgx.Conn
}

func NewAuthPostgres(db *pgx.Conn) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (a *AuthPostgres) CreateUser(ctx context.Context, login string, password string) (uuid.UUID, error) {
	var userID uuid.UUID

	err := a.db.QueryRow(ctx, "insert into users (user_name, user_password) values ($1, $2) returning id;", login, password).Scan(&userID)
	if err != nil {
		return uuid.Nil, mapStorageErr(err)
	}

	return userID, nil
}

func (a *AuthPostgres) GetUser(ctx context.Context, login string, password string) (uuid.UUID, error) {
	var userID uuid.UUID

	err := a.db.QueryRow(ctx, "select id from users where user_name = $1 and user_password = $2", login, password).Scan(&userID)
	if err != nil {
		return uuid.Nil, mapStorageErr(err)
	}

	return userID, nil
}
