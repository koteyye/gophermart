package accrual

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/sergeizaitcev/gophermart/internal/accrual/config"
	"github.com/sergeizaitcev/gophermart/internal/accrual/server"
	"github.com/sergeizaitcev/gophermart/internal/accrual/service"
	"github.com/sergeizaitcev/gophermart/internal/accrual/storage/postgres"
	"github.com/sergeizaitcev/gophermart/pkg/commands"
)


func Run(ctx context.Context) error {
	cmd := commands.New("accrual", runServer)
	return cmd.Execute(ctx)
}

func runServer(ctx context.Context, c *config.Config) error {
	storage, err := postgres.Connect(c)
	if err != nil {
		return fmt.Errorf("create a new storage: %w", err)
	}
	defer storage.Close()

	err = storage.Up(ctx)
	if err != nil {
		return fmt.Errorf("migration up: %w", err)
	}

	service := service.NewService(storage)
	mux := server.NewHandler(service)

	return listenAndServe(ctx, c.RunAddress, mux)
}

func listenAndServe(ctx context.Context, addr string, handler http.Handler) error {
	srv := &http.Server{
		Addr:        addr,
		Handler:     handler,
		BaseContext: func(net.Listener) context.Context { return ctx },
	}

	errc := make(chan error)
	go func() { errc <- srv.ListenAndServe() }()

	select {
	case <-ctx.Done():
	case err := <-errc:
		return fmt.Errorf("unexpected listening error: %w", err)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer shutdownCancel()

	err := srv.Shutdown(shutdownCtx)
	<-errc

	if err != nil {
		return fmt.Errorf("graceful shutdown of the server: %w", err)
	}

	return nil
}