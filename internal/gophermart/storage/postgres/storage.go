package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/sergeizaitcev/gophermart/deployments/gophermart/migrations"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/config"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/storage"
)

var _ storage.Storage = (*Storage)(nil)

type Storage struct {
	db *sql.DB
}

func Connect(c *config.Config) (*Storage, error) {
	db, err := newConn(c.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("creating a new database connection: %w", err)
	}
	return &Storage{db: db}, nil
}

func newConn(dsn string) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("connection to the database: %w", err)
	}
	defer func() {
		if err != nil {
			_ = db.Close()
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping connections: %w", err)
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

func (s *Storage) transaction(ctx context.Context, fn func(*sql.Tx) error) (err error) {
	var tx *sql.Tx

	tx, err = s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	err = fn(tx)
	if err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
