package domain

import (
	"context"
)

// AccrualClient описывает интерфейс клиента accrual.
//
//go:generate mockgen -source=contract.go -destination=mocks/mocks.go
type AccrualClient interface {
	// GetAccrualInfo возвращает информацию о расчёте начислений баллов
	// лояльности за совершённый заказ.
	GetAccrualInfo(ctx context.Context, number OrderNumber) (AccrualInfo, error)
}

// AuthService описывает интерфейс сервиса регистрации и аутентификации
// пользователя.
//
//go:generate mockgen -source=contract.go -destination=mocks/mocks.go
type AuthService interface {
	// Identify идентифицирует уникальный идентификатор пользователя.
	Identify(ctx context.Context, id UserID) error

	// SignIn выполняет вход пользователя и возвращает его уникальный
	// идентификатор.
	SignIn(ctx context.Context, auth Authentication) (UserID, error)

	// SignUp выполняет регистрацию нового пользователя и возвращает его
	// уникальный идентификатор.
	SignUp(ctx context.Context, auth Authentication) (UserID, error)
}

// UserService описывает интерфейс сервиса для работы с пользователем.
//
//go:generate mockgen -source=contract.go -destination=mocks/mocks.go
type UserService interface {
	// GetBalance возвращает баланс пользователя.
	GetBalance(ctx context.Context, id UserID) (UserBalance, error)
}

// OrderService описывает интерфейс сервиса обработки заказов пользователя.
//
//go:generate mockgen -source=contract.go -destination=mocks/mocks.go
type OrderService interface {
	// GetOrders возвращает список заказов пользователя.
	GetOrders(ctx context.Context, id UserID) ([]Order, error)

	// Process добавляет заказ пользователя в обработку.
	Process(ctx context.Context, order Order) error
}

// OperationService описывает интерфейс сервиса обработки балансовых операций.
//
//go:generate mockgen -source=contract.go -destination=mocks/mocks.go
type OperationService interface {
	// GetOperations возвращает все балансовые операции пользователя.
	GetOperations(ctx context.Context, id UserID) ([]Operation, error)

	// Perform выполняет балансовую операцию.
	Perform(ctx context.Context, operation Operation) error
}
