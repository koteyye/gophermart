package service

import (
	"log/slog"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/storage"
	"github.com/sergeizaitcev/gophermart/pkg/sign"
)

// Service определяет бизнес-логику gophermart.
type Service struct {
	Auth *userService

	logger *slog.Logger
}

// NewService возвращает новый экземпляр Service.
func NewService(logger *slog.Logger, storage storage.Storage, signer sign.Signer) *Service {
	return &Service{
		Auth:   newUserService(storage, signer),
		logger: logger,
	}
}
