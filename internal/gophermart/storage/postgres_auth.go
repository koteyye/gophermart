package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/sergeizaitcev/gophermart/pkg/passwords"
)

type AuthPostgres struct {
	db *pgx.Conn
}

func NewAuthPostgres(db *pgx.Conn) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (a *AuthPostgres) CreateUser(
	ctx context.Context,
	login string,
	password string,
) (uuid.UUID, error) {
	var userID uuid.UUID

	hashedPassword, err := passwords.Hash(password)
	if err != nil {
		return uuid.Nil, fmt.Errorf("password hashing: %w", err)
	}

	err = a.db.QueryRow(ctx, "insert into users (user_name, hashed_password) values ($1, $2) returning id;", login, hashedPassword).
		Scan(&userID)
	if err != nil {
		return uuid.Nil, mapStorageErr(err)
	}

	return userID, nil
}

func (a *AuthPostgres) GetUser(
	ctx context.Context,
	login string,
	password string,
) (uuid.UUID, error) {
	var (
		userID         uuid.UUID
		hashedPassword string
	)

	err := a.db.QueryRow(ctx, "select id, hashed_password from users where user_name = $1;", login).
		Scan(&userID, &hashedPassword)
	if err != nil {
		return uuid.Nil, mapStorageErr(err)
	}

	if !passwords.Compare(hashedPassword, password) {
		return uuid.Nil, errors.New("invalid password")
	}

	return userID, nil
}
