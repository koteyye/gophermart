package storage

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"

	"github.com/sergeizaitcev/gophermart/deployments/gophermart/migrations"
)

// Тесты для проверки CRUD на тестовой базе

const test_dsn = "postgresql://postgres:postgres@localhost:5433/gophermart?sslmode=disable"

func testDB(t *testing.T) (*pgxpool.Pool, func()) {
	ctx := context.Background()

	db, err := pgxpool.New(ctx, test_dsn)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	require.NoError(t, db.Ping(ctx))

	sql, err := goose.OpenDBWithDriver("pgx", test_dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, sql.Close()) })

	require.NoError(t, migrations.Up(ctx, sql))

	return db, func() { require.NoError(t, migrations.Down(ctx, sql)) }
}
