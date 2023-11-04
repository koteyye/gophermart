package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -source=storage.go -destination=mocks/mock.go

// Auth методы для CRUD с users, используется в signIn и signUp
type Auth interface {
	CreateUser(ctx context.Context, login string, password string) (uuid.UUID, error)
	GetUser(ctx context.Context, login string, passwrod string) (uuid.UUID, error)
}

// GophermartDB CRUD операции с БД
type GophermartDB interface {
	Orders
	Balance
}

// Orders - CRUD с заказами
type Orders interface {
	CreateOrder(ctx context.Context, orderNumber int64) (uuid.UUID, error)
	UpdateOrder(ctx context.Context, order *UpdateOrder) error
	UpdateOrderStatus(ctx context.Context, orderNumber int64, orderStatus Status) error
	DeleteOrderByNumber(ctx context.Context, orderNumber int64) error
	GetOrderByNumber(ctx context.Context, orderNumber int64) (*OrderItem, error)
}

// Balance - CRUD с балансом
type Balance interface {
	// CRUD BalanceOperation
	CreateBalanceOperation(ctx context.Context, operation *BalanceOperationItem) (uuid.UUID, error)
	UpdateBalanceOperation(ctx context.Context, operation *BalanceOperationItem) error
	DeleteBalanceOperation(ctx context.Context, operationID uuid.UUID) error
	GetBalanceOperationByOrderID(ctx context.Context, orderID uuid.UUID) (*BalanceOperationItem, error)
	GetBalanceOperationByBalanceID(ctx context.Context, balanceID uuid.UUID) ([]*BalanceOperationItem, error)
	// CRUD Balance
	GetBalanceByUserID(ctx context.Context, userID uuid.UUID) (*BalanceItem, error)
	UpdateBalance(ctx context.Context, userID uuid.UUID, currentSum int64) (error)
	IncrementBalance(ctx context.Context, userID uuid.UUID, currentSum int64) (error)
	DecrementBalance(ctx context.Context, userID uuid.UUID, currentSum int64) (error)
}


type Storage struct {
	Auth
	GophermartDB
}

func NewStorage(db *pgx.Conn) *Storage {
	return &Storage{
		Auth: NewAuthPostgres(db),
		GophermartDB: NewGophermartPostgres(db),
	}
}