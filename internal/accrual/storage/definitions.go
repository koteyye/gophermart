package storage

import "github.com/google/uuid"

// OrderStatus тип статуса заказа
type OrderStatus uint8

// 
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
	Price int
}

// Order структура для обновления записи в таблице orders
type Order struct {
	OrderID uuid.UUID
	Status OrderStatus
	Accrual int
}
