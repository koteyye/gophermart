package accrual

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	
	"github.com/sergeizaitcev/gophermart/internal/accrual/config"
	"github.com/sergeizaitcev/gophermart/internal/accrual/server"
	"github.com/sergeizaitcev/gophermart/internal/accrual/service"
	"github.com/sergeizaitcev/gophermart/internal/accrual/storage/postgres"
	"github.com/sergeizaitcev/gophermart/pkg/commands"
	"github.com/sergeizaitcev/gophermart/pkg/httpserver"
)


func Run(ctx context.Context) error {
	cmd := commands.New("accrual", runAccrual)
	return cmd.Execute(ctx)
}

func runAccrual(ctx context.Context, c *config.Config) error {
	storage, err := postgres.Connect(c)
	if err != nil {
		return fmt.Errorf("create a new storage: %w", err)
	}
	defer storage.Close()

	err = storage.Up(ctx)
	if err != nil {
		return fmt.Errorf("migration up: %w", err)
	}

	logger := newLogger(c)

	service := service.NewService(storage)
	handler := server.NewHandler(logger, service)

	return httpserver.ListenAndServe(ctx, c.RunAddress, handler)
}

func newLogger(c *config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{Level: c.Level}
	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}