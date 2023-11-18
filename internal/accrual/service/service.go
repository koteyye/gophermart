package service

import (
	"github.com/sergeizaitcev/gophermart/internal/accrual/storage"
)

// Service определяет бизнес-логику accrual
type Service struct {
	Accrual *Accrual
}

// NewService возвращает экземпляр Service
func NewService(s storage.Storage) *Service {
	return &Service{
		Accrual: NewAccrual(s),
	}
}
