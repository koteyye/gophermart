package storage

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pressly/goose/v3"
	"github.com/sergeizaitcev/gophermart/deployments/gophermart/migrations"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Тесты для проверки CRUD на тестовой базе

const test_dsn = "postgresql://localhost:5433/gophermart?user=postgres&password=postgres&sslmode=disable"

var testAuthDB *AuthPostgres

func TestMain(m *testing.M) {
	ctx := context.Background()

	db, err := pgx.Connect(ctx, test_dsn)
	if err != nil {
		slog.Error("can't connect db: %w", err)
	}
	defer db.Close(ctx)

	if err := db.Ping(ctx); err != nil {
		slog.Error("can't ping db: %w", err)
	}

	sql, err := goose.OpenDBWithDriver("pgx", test_dsn)
	if err != nil {
		slog.Error("can't open db with driver: %w", err)
	}

	if err := migrations.Up(ctx, sql); err != nil {
		slog.Error("can't migration: %w", err)
	}

	testAuthDB = NewAuthPostgres(db)

	_, err = db.Exec(ctx, "truncate table users cascade;")
	if err != nil {
		slog.Error("can't truncate table test db")
	}
	os.Exit(m.Run())
}
