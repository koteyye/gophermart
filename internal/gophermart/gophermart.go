package gophermart

import (
	"context"
	"fmt"
	"os"

	"log/slog"

	"github.com/sergeizaitcev/gophermart/deployments/gophermart/migrations"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/clients/accrual"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/config"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/handler"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/service"
	"github.com/sergeizaitcev/gophermart/pkg/commands"
	"github.com/sergeizaitcev/gophermart/pkg/httpserver"
	"github.com/sergeizaitcev/gophermart/pkg/postgres"
	"github.com/sergeizaitcev/gophermart/pkg/sign"
)

// Run запускает gophermart и блокируется до тех пор, пока не сработает
// контекст или функция не вернёт ошибку.
func Run(ctx context.Context) error {
	cmd := commands.New("gophermart", runGophermart)
	return cmd.Execute(ctx)
}

func runGophermart(ctx context.Context, c *config.Config) error {
	configureLogger(c)

	signer, err := newSigner(c)
	if err != nil {
		return fmt.Errorf("creating a new signer: %w", err)
	}

	db, err := postgres.Connect(c.DatabaseURI)
	if err != nil {
		return fmt.Errorf("connecting to the postgres: %w", err)
	}
	defer db.Close()

	err = migrations.Up(ctx, db)
	if err != nil {
		return fmt.Errorf("migration up: %w", err)
	}

	accrual := accrual.NewClient(c.AccrualSystemAddress, nil)

	orders := service.NewOrders(db, accrual)
	defer orders.Close()

	handler := handler.New(handler.HandlerOptions{
		Auth:       service.NewAuth(db),
		Orders:     orders,
		Users:      service.NewUsers(db),
		Operations: service.NewOperations(db),
		Signer:     signer,
	})

	return httpserver.ListenAndServe(ctx, c.RunAddress, handler)
}

func newSigner(c *config.Config) (sign.Signer, error) {
	secretKey, err := c.SecretKey()
	if err != nil {
		return nil, fmt.Errorf("getting a secret key: %w", err)
	}
	return sign.New(secretKey, sign.WithTTL(c.TokenTTL)), nil
}

func configureLogger(c *config.Config) {
	opts := &slog.HandlerOptions{Level: c.Level}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))
}
