package storage

import (

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


//Структуры для Read операций
// Order номер заказа
type OrderItem struct {
	ID uuid.UUID `db:"id"`
	OrderNumber int64 `db:"order_number"`
	Status string `db:"status"`
	Accrual int64 `db:"accrual"`
	CreatedAt string `db:"created_at"`
	UploadAt string `db:"updated_at"`
	UserID uuid.UUID `db:"user_created"`
}

// Balance баланс пользователя
type BalanceItem struct {
	ID uuid.UUID `db:"id"`
	UserID uuid.UUID `db:"user_id"`
	CurrentBalance int64 `db:"current_balance"`
}

// BalanceOperation операции с балансом по заказам
type BalanceOperationItem struct {
	ID uuid.UUID `db:"id"`
	OrderID uuid.UUID `db:"order_id"`
	BalanceID uuid.UUID `db:"balance_id"`
	SumOperation int64 `db:"sum_operation"`
}

//Структуры для Create / Update операций
// UpdateOrder обновление заказа
type UpdateOrder struct {
	Number int64
	Status string
	Accrual int64
}


