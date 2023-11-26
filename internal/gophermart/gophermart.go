package gophermart

import (
	"context"
	"fmt"
	"os"

	"log/slog"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/config"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/handlers"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/service"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/storage/postgres"
	"github.com/sergeizaitcev/gophermart/pkg/commands"
	"github.com/sergeizaitcev/gophermart/pkg/httpserver"
	"github.com/sergeizaitcev/gophermart/pkg/sign"
)

// Run запускает gophermart и блокируется до тех пор, пока не сработает
// контекст или функция не вернёт ошибку.
func Run(ctx context.Context) error {
	cmd := commands.New("gophermart", runGophermart)
	return cmd.Execute(ctx)
}

func runGophermart(ctx context.Context, c *config.Config) error {
	signer, err := newSigner(c)
	if err != nil {
		return fmt.Errorf("creating a new signer: %w", err)
	}

	storage, err := postgres.Connect(c)
	if err != nil {
		return fmt.Errorf("creating a new storage: %w", err)
	}
	defer storage.Close()

	err = storage.Up(ctx)
	if err != nil {
		return fmt.Errorf("migration up: %w", err)
	}

	logger := newLogger(c)

	service := service.NewService(service.ServiceOptions{
		Logger:  logger,
		Accrual: nil,
		Storage: storage,
		Signer:  signer,
	})
	defer service.Close()

	handler := handlers.NewHandler(logger, service)

	return httpserver.ListenAndServe(ctx, c.RunAddress, handler)
}

func newSigner(c *config.Config) (sign.Signer, error) {
	secretKey, err := c.SecretKey()
	if err != nil {
		return nil, fmt.Errorf("getting a secret key: %w", err)
	}
	return sign.New(secretKey, sign.WithTTL(c.TokenTTL)), nil
}

func newLogger(c *config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{Level: c.Level}
	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}
