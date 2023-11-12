package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/models"
	"github.com/sergeizaitcev/gophermart/pkg/passwords"
)

type AuthPostgres struct {
	db *pgxpool.Pool
}

func NewAuthPostgres(db *pgxpool.Pool) *AuthPostgres {
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

	tx, err := a.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return uuid.Nil, fmt.Errorf("transaction error: %w", err)
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, "insert into users (user_name, hashed_password) values ($1, $2) returning id;", login, hashedPassword).
		Scan(&userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create user error: %w", mapStorageErr(err))
	}

	_, err = tx.Exec(ctx, "insert into balance (user_id) values ($1)", userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create user balance error: %w", mapStorageErr(err))
	}
	tx.Commit(ctx)

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
		return uuid.Nil, models.ErrInvalidPassword
	}

	return userID, nil
}

func (a *AuthPostgres) GetUserByID(ctx context.Context, userID uuid.UUID) (string, error) {
	var login string

	err := a.db.QueryRow(ctx, "select user_name from users where id = $1", userID).Scan(&login)
	if err != nil {
		return "", fmt.Errorf("select login by userID err: %w", mapStorageErr(err))
	}

	return login, nil
}