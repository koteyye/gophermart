package storage

import (
	"context"
	"errors"

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
	return uuid.Nil, errors.New("not implemented")
}

func (a *AuthPostgres) GetUser(ctx context.Context, login string, password string) (uuid.UUID, error) {
	return uuid.Nil, errors.New("not implemented")
}