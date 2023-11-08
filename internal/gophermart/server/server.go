package server

import (
	"context"

	"log/slog"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/config"
)

type Server struct {
	logger *slog.Logger
}

func New(config *config.Config) *Server {
	return nil
}

func (s *Server) Run(ctx context.Context) error {
	return nil
}
