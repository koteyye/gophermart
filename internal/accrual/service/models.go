package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sergeizaitcev/gophermart/pkg/monetary"
)

// ошибки service
var (
	ErrOrderRegistered = errors.New("order registered")
)

// workerOrder структура заказа для обработчика
type workerOrder struct {
	orderID uuid.UUID
	goods   []*workerGoods
	accrual monetary.Unit
}

// workerGoods структура для элемента содержимого заказа
type workerGoods struct {
	goodID     uuid.UUID
	matchID    uuid.UUID
	price      monetary.Unit
	reward     float64
	rewardType string
	accrual    monetary.Unit
}

const (
	percent = "PERCENT"
	natural = "NATURAL"
)
