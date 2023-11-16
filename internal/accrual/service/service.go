package service

import "github.com/sergeizaitcev/gophermart/internal/accrual/storage"

// Service определяет бизнес-логику accrual
type Service struct {
	accrual *Accrual
}

func NewService(s storage.Storage) *Service {
	return &Service{
		accrual: NewAccrual(s),
	}
}
