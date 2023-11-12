package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/sergeizaitcev/gophermart/internal/accrual/config"
	"github.com/sergeizaitcev/gophermart/internal/accrual/service"
	"github.com/sergeizaitcev/gophermart/internal/accrual/storage"
)

// Server определяет HTTP-сервер для accrual
type Server struct {
	config *config.Config
}

// New возвращает новый экземпляр Server
func New(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

// Run запускает сервер и блокируется до тех пор, пока не сработает контекст
// или функция не вернет ошибку
func (s *Server) Run(ctx context.Context) error {
	storage, err := storage.NewStorage(ctx, s.config)
	if err != nil {
		return fmt.Errorf("create a new storage: %w", err)
	}
	defer storage.Close()

	service := service.NewService(storage)
	mux := NewHandler(service)

	return listenAndServe(ctx, s.config.RunAddress, mux)
}

func listenAndServe(ctx context.Context, addr string, handler http.Handler) error {
	srv := &http.Server{
		Addr: addr,
		Handler: handler,
		BaseContext: func(net.Listener) context.Context {return ctx},
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