package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type authPostgres struct {
	db *pgx.Conn
}

func NewAuthPostgres(db *pgx.Conn) *authPostgres {
	return &authPostgres{db: db}
}

func (a *authPostgres) CreateUser(ctx context.Context, login string, password string) (uuid.UUID, error) {
	return uuid.Nil, errors.New("not implemented")
}

func (a *authPostgres) GetUser(ctx context.Context, login string, password string) (uuid.UUID, error) {
	return uuid.Nil, errors.New("not implemented")
}