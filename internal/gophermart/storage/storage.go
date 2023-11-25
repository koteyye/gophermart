package storage

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"

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
//go:generate mockgen -source=storage.go -destination=mocks/mock.go
type Storage interface {
	Auth
	Balance
	Orders
	Operations
}

// Auth описывает интерфейс для регистрации и авторизации пользователя.
type Auth interface {
	// CreateUser регистрирует нового пользователя по связке логин-пароль
	// и возвращает его ID.
	CreateUser(ctx context.Context, login string, password string) (uuid.UUID, error)

	// GetUser возвращает ID пользователя по связке логин-пароль.
	GetUser(ctx context.Context, login string, password string) (uuid.UUID, error)

	// GetLogin возвращает логин пользователя по ID.
	GetLogin(ctx context.Context, userID uuid.UUID) (string, error)
}

// Balance описывает интерфейс ведения баланса пользователя.
type Balance interface {
	// GetBalance возвращает баланс пользователя.
	GetBalance(ctx context.Context, userID uuid.UUID) (*UserBalance, error)

	// UpdateBalance обновляет баланс пользователя.
	UpdateBalance(ctx context.Context, userID uuid.UUID, amount monetary.Unit) error
}

// Orders описывает интерфейс обработки заказов пользователя.
type Orders interface {
	// CreateOrder создает новый заказ пользователя.
	CreateOrder(ctx context.Context, userID uuid.UUID, order string) error

	// GetOrder возвращает заказ пользователя.
	GetOrder(ctx context.Context, order string) (*Order, error)

	// GetOrders возвращает все заказы пользователя.
	GetOrders(ctx context.Context, userID uuid.UUID) ([]Order, error)

	// UpdateOrder обновляет заказ пользователя.
	UpdateOrder(ctx context.Context, order string, status OrderStatus, accrual monetary.Unit) error

	// UpdateOrderStatus обновляет статус заказа пользователя.
	UpdateOrderStatus(ctx context.Context, order string, status OrderStatus) error

	// DeleteOrder удаляет заказ пользователя.
	DeleteOrder(ctx context.Context, order string) error
}

// Operations описывает интерфейс обработки балансовых операций.
type Operations interface {
	// CreateOperation создает новую балансовую операцию.
	CreateOperation(ctx context.Context, userID uuid.UUID, order string, amount monetary.Unit) error

	// GetOperation вовзращает балансовую операцию по номеру заказа.
	GetOperation(ctx context.Context, order string) (*Operation, error)

	// GetOperations возвращает все балансовые операции пользователя.
	GetOperations(ctx context.Context, userID uuid.UUID) ([]Operation, error)

	// UpdateOperationStatus изменяет статус балансовой операции по номеру заказа.
	UpdateOperationStatus(ctx context.Context, order string, status OperationStatus) error

	// DeleteOperation удаляет балансовую операцию по номеру заказа.
	DeleteOperation(ctx context.Context, order string) error
}

// UserBalance определяет баланс пользователя.
type UserBalance struct {
	UserID    uuid.UUID
	Amount    monetary.Unit
	Withdrawn monetary.Unit
}

// OrderStatus определяет статус заказа.
type OrderStatus uint8

const (
	OrderStatusUnknown OrderStatus = iota

	OrderStatusNew        // Заказ создан.
	OrderStatusProcessing // Заказ в обработке.
	OrderStatusProcessed  // Заказ обработан.
	OrderStatusInvalid    // Заказ недействителен.
)

var orderValues = []string{
	"unknown",
	"new",
	"processing",
	"processed",
	"invalid",
}

func (s OrderStatus) String() string {
	if s > 0 && int(s) < len(orderValues) {
		return orderValues[s]
	}
	return orderValues[0]
}

func (s *OrderStatus) Scan(value any) error {
	if value == nil {
		*s = OrderStatusUnknown
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("unsupported type: %T", value)
	}

	*s = OrderStatusUnknown
	str := strings.ToLower(string(b))

	for i, v := range orderValues {
		if v == str {
			*s = OrderStatus(i)
			break
		}
	}

	return nil
}

func (s OrderStatus) Value() (driver.Value, error) {
	if s == OrderStatusUnknown {
		return nil, nil
	}
	return s.String(), nil
}

// Order определяет пользовательский заказ.
type Order struct {
	UserID    uuid.UUID
	Number    string
	Status    OrderStatus
	Accrual   monetary.Unit
	UpdatedAt time.Time
}

// OperationStatus определяет статус операции.
type OperationStatus uint8

const (
	OperationStatusUnknown OperationStatus = iota

	OperationStatusRun   // Операция в обработке.
	OperationStatusDone  // Операция выполнена.
	OperationStatusError // Ошибка во время выполнения операции.
)

var operationStatusValues = []string{
	"unknown",
	"run",
	"done",
	"error",
}

func (s OperationStatus) String() string {
	if s > 0 && int(s) < len(operationStatusValues) {
		return operationStatusValues[s]
	}
	return operationStatusValues[0]
}

func (s *OperationStatus) Scan(value any) error {
	if value == nil {
		*s = OperationStatusUnknown
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("unsupported type: %T", value)
	}

	*s = OperationStatusUnknown
	low := strings.ToLower(string(b))

	for i, v := range operationStatusValues {
		if v == low {
			*s = OperationStatus(i)
			break
		}
	}

	return nil
}

func (s OperationStatus) Value() (driver.Value, error) {
	if s == OperationStatusUnknown {
		return nil, nil
	}
	return s.String(), nil
}

// Operation определяет балансовую операцию.
type Operation struct {
	UserID      uuid.UUID
	OrderNumber string
	Amount      monetary.Unit
	Status      OperationStatus
	UpdatedAt   time.Time
}
