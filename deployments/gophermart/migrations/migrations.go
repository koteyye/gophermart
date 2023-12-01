package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"log/slog"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var fsys embed.FS

func lazyInit() error {
	logger := slog.NewLogLogger(slog.Default().Handler(), slog.LevelInfo)

	goose.SetLogger(logger)
	goose.SetBaseFS(fsys)

	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("dialect could not be set: %w", err)
	}

	return nil
}

// Up запускает миграцию в БД.
func Up(ctx context.Context, db *sql.DB) error {
	err := lazyInit()
	if err != nil {
		return fmt.Errorf("initializing the migrator: %w", err)
	}
	if err := goose.UpContext(ctx, db, "."); err != nil {
		return fmt.Errorf("migration up: %w", err)
	}
	return nil
}

// Down откатывает миграцию в БД.
func Down(ctx context.Context, db *sql.DB) error {
	err := lazyInit()
	if err != nil {
		return fmt.Errorf("initializing the migrator: %w", err)
	}
	if err := goose.DownContext(ctx, db, "."); err != nil {
		return fmt.Errorf("migration down: %w", err)
	}
	return nil
}
