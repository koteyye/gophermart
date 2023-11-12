package service

import "github.com/sergeizaitcev/gophermart/internal/accrual/storage"

// Service определяет бизнес-логику accrual
type Service struct {
	Accrual *Accrual
}

func NewService(storage *storage.Storage) *Service {
	return &Service{
		Accrual: NewAccrual(storage.Accrual()),
	}
}