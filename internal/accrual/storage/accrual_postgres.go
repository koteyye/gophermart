package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccrualPostgres struct {
	db *pgxpool.Pool
}

func NewAccrualPostgres(db *pgxpool.Pool) *AccrualPostgres {
	return &AccrualPostgres{db: db}
}

func (a *AccrualPostgres) CreateOrder(ctx context.Context, order string) (uuid.UUID, error) {
	return uuid.Nil, errors.New("not implemented")
}