package domain

import (
	"errors"
	"time"
)

// Ошибки, возвращаемые сервисом accrual.
var (
	// ErrOrderNotRegistered возвращается, если номер заказа не зарегистрирован.
	ErrOrderNotRegistered = errors.New("order is not registered")

	// ErrInternalServerError возвращается, если сервер вернул 500 код ответа.
	ErrInternalServerError = errors.New("internal server error")
)

// ResourceExhaustedError возвращается, если клиент превысил лимит запросов
// в минуту.
type ResourceExhaustedError struct {
	Message    string
	RetryAfter time.Duration
}

func (err *ResourceExhaustedError) Error() string {
	return err.Message
}

// Ошибки, возвращаемые хранилищем gophermart.
var (
	// ErrDuplicate возвращается, когда заказ уже был загружен.
	ErrDuplicate = errors.New("duplicate value")

	// ErrDuplicateOtherUser возвращается, когда заказ уже был загружен другим
	// пользователем.
	ErrDuplicateOtherUser = errors.New("duplicate value from other user")

	// ErrNotFound возвращается, если значение не найдено.
	ErrNotFound = errors.New("value not found")

	// ErrBalanceBelowZero возвращается, когда баланс после проведения операции
	// будет ниже нуля.
	ErrBalanceBelowZero = errors.New("balance can't be below zero")

	// ErrInvalidPassword возвращается, когда пользователь передал не верный
	// пароль.
	ErrInvalidPassword = errors.New("invalid password")
)
