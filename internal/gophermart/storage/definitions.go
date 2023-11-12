package storage

import (
	"time"

	"github.com/google/uuid"
)

// Возможные статусы заказа
type Status uint8

const (
	New Status = iota
	Processing
	Processed
	Invalid
)

func (s Status) String() string {
	switch s {
	case New:
		return "new"
	case Processing:
		return "processing"
	case Processed:
		return "processed"
	case Invalid:
		return "invalid"
	}
	return "unknow"
}

// BalanceOperationState Возможные состояния балансовой операции
type BalanceOperationState uint8

const (
	Run BalanceOperationState = iota
	Done
	Error
)

func (b BalanceOperationState) String() string {
	switch b {
	case Run:
		return "run"
	case Done:
		return "done"
	case Error:
		return "error"
	}
	return "unknow"
}

// Структуры для Read операций

// OrderItem заказ
type OrderItem struct {
	ID      uuid.UUID `db:"id"`
	Order   string `db:"order_number"`
	Status  string `db:"status"`
	Accrual int64 `db:"accrual"`
	UserID  uuid.UUID `db:"user_created"`
	UploadedAt time.Time `db:"updated_at"`
}

// BalanceItem баланс пользователя
type BalanceItem struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	CurrentBalance int64
	Withdrawn int64
}

// BalanceOperationItem BalanceOperation операции с балансом по заказам
type BalanceOperationItem struct {
	ID           uuid.UUID `db:"id"`
	OrderID      uuid.UUID `db:"order_id"`
	BalanceID    uuid.UUID `db:"balance_id"`
	SumOperation int64     `db:"sum_operation"`
	UpdatedAt time.Time `db:"updated_at"`
}

//Структуры для Create / Update операций

// UpdateOrder обновление заказа
type UpdateOrder struct {
	Order   string
	Status  Status
	Accrual int64
}
