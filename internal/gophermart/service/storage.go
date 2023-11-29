package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/sergeizaitcev/gophermart/pkg/monetary"
)

var (
	ErrDuplicate          = errors.New("duplicate value")
	ErrDuplicateOtherUser = errors.New("duplicate value from other user")
	ErrNotFound           = errors.New("value not found")
	ErrOther              = errors.New("other storage error")
	ErrBalanceBelowZero   = errors.New("balance can't be below zero")
	ErrInvalidPassword    = errors.New("invalid password")
)

// Storage описывает интерфейс хранилища gophermart.
//
//go:generate mockgen -source=storage.go -destination=mocks/storage.go
type Storage interface {
	Users
	Orders
	Operations
}

// Users описывает интерфейс для взаимодействия с пользователем.
type Users interface {
	// CreateUser регистрирует нового пользователя и возвращает его ID.
	CreateUser(ctx context.Context, u User) (uuid.UUID, error)

	// User возвращает ID пользователя.
	UserID(ctx context.Context, u User) (uuid.UUID, error)

	// UserExists возвращает true, если пользователем с таким ID существует.
	UserExists(ctx context.Context, userID uuid.UUID) (bool, error)

	// Balance возвращает баланс пользователя.
	Balance(ctx context.Context, userID uuid.UUID) (*Balance, error)
}

// Orders описывает интерфейс обработки заказов пользователя.
type Orders interface {
	// CreateOrder создает новый заказ пользователя.
	CreateOrder(ctx context.Context, userID uuid.UUID, order string) error

	// OrderStatus возвращает статус заказа пользователя.
	OrderStatus(ctx context.Context, order string) (OrderStatus, error)

	// Orders возвращает все заказы пользователя.
	Orders(ctx context.Context, userID uuid.UUID) ([]Order, error)

	// ProccessOrder обрабатывает заказ.
	ProcessOrder(ctx context.Context, order string, accrual monetary.Unit) error

	// UpdateOrderStatus обновляет статус заказа пользователя.
	UpdateOrderStatus(ctx context.Context, order string, status OrderStatus) error
}

// Operations описывает интерфейс обработки балансовых операций.
type Operations interface {
	// CreateOperation создает новую балансовую операцию и возвращает её ID.
	CreateOperation(
		ctx context.Context,
		userID uuid.UUID,
		order string,
		amount monetary.Unit,
	) (uuid.UUID, error)

	// Operations возвращает все балансовые операции пользователя.
	Operations(ctx context.Context, userID uuid.UUID) ([]Operation, error)

	// UpdateOperationStatus изменяет статус балансовой операции по номеру заказа.
	UpdateOperationStatus(ctx context.Context, order string, status OperationStatus) error

	// PerformOperation выполняет операцию над балансом пользователя.
	PerformOperation(ctx context.Context, operationID uuid.UUID) error
}
