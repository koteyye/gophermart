package service

import (
	"github.com/sergeizaitcev/gophermart/internal/gophermart/storage"
	"github.com/sergeizaitcev/gophermart/pkg/sign"
)

// Service определяет бизнес-логику gophermart.
type Service struct {
	Auth *Auth
}

// NewService возвращает новый экземпляр Service.
func NewService(storage *storage.Storage, signer sign.Signer) *Service {
	return &Service{
		Auth: NewAuth(storage.Auth(), signer),
	}
}
