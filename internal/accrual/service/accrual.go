package service

import (
	"github.com/sergeizaitcev/gophermart/internal/accrual/storage"
)

// Accrual определяет сервис расчета вознаграждений
type Accrual struct {
	storage storage.Accrual
}

// NewAccrual возвращает новый экземпляр Accrual
func NewAccrual(accrual storage.Accrual) *Accrual {
	return &Accrual{
		storage: accrual,
	}
}
