package service

import (
	"github.com/google/uuid"
	"github.com/sergeizaitcev/gophermart/pkg/monetary"
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
	reward     int
	rewardType string
	accrual    monetary.Unit
}

const (
	percent = "percent"
	natural = "natural"
)
