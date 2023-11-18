package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sergeizaitcev/gophermart/deployments/accrual/migrations"
	"github.com/sergeizaitcev/gophermart/internal/accrual/config"
	"github.com/sergeizaitcev/gophermart/internal/accrual/storage"
)

var _ storage.Storage = (*Storage)(nil)

// Storage ...
type Storage struct {
	db *sql.DB
}

// NewStorage ...
func Connect(c *config.Config) (*Storage, error) {
	db, err := connect(c.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("connectiong to postgres: %w", err)
	}
	return &Storage{db: db}, nil
}

func connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("create a new connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("database ping: %w", err)
	}

	return db, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) Up(ctx context.Context) error {
	return migrations.Up(ctx, s.db)
}

func (s *Storage) Down(ctx context.Context) error {
	return migrations.Down(ctx, s.db)
}

func (s *Storage) transaction(
	ctx context.Context,
	fn func(*sql.Tx) error,
) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	err = fn(tx)
	if err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}