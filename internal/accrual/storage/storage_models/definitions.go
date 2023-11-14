package storage_models

import "github.com/google/uuid"

// OrderStatus тип статуса заказа
type OrderStatus uint8

const (
	Registered = iota
	Invalid
	Processing
	Processed
)

func (o OrderStatus) String() string {
	switch o {
	case Registered:
		return "registered"
	case Invalid:
		return "invalid"
	case Processing:
		return "processing"
	case Processed:
		return "processed"
	}
	return "unknow"
}

// Goods структура для создания записи в таблице goods
type Goods struct {
	MatchID uuid.UUID
	Price   int
}

// Order структура для обновления записи в таблице orders
type Order struct {
	OrderID uuid.UUID
	Status  OrderStatus
	Accrual int
}

// Match структура для создания записи в таблице matches
type Match struct {
	MatchName string
	Reward    int
	Type      RewardType
}

type RewardType uint8

const (
	Percent = iota
	Natural
)

func (r RewardType) String() string {
	switch r {
	case Percent:
		return "percent"
	case Natural:
		return "natural"
	}
	return "unknow"
}
