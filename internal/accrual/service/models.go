package service

import (
	"github.com/google/uuid"
)

// workerOrder структура заказа для обработчика
type workerOrder struct {
	orderID uuid.UUID
	goods   []*workerGoods
	accrual float64
}

// workerGoods структура для элемента содержимого заказа
type workerGoods struct {
	goodID     uuid.UUID
	matchID    uuid.UUID
	price      float64
	reward     float64
	rewardType string
	accrual    float64
}

const (
	percent = "percent"
	natural = "natural"
)
