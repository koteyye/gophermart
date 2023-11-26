package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"

	"log/slog"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var fsys embed.FS

func init() {
	logger := slog.NewLogLogger(slog.NewTextHandler(os.Stdout, nil), slog.LevelInfo)

	goose.SetLogger(logger)
	goose.SetBaseFS(fsys)

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Fatalf("migrations: set dialect: %s", err)
	}
}

// Up запускает миграцию в БД.
func Up(ctx context.Context, db *sql.DB) error {
	if err := goose.UpContext(ctx, db, "."); err != nil {
		return fmt.Errorf("migrations up: %w", err)
	}
	return nil
}

// Down откатывает миграцию в БД.
func Down(ctx context.Context, db *sql.DB) error {
	if err := goose.DownContext(ctx, db, "."); err != nil {
		return fmt.Errorf("migrations down: %w", err)
	}
	return nil
}
