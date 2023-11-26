package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/sergeizaitcev/gophermart/pkg/monetary"
)

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

// AccrualClient определяет интерфейс клиента accrual.
//
//go:generate mockgen -source=storage.go -destination=mocks/accrual_client.go
type AccrualClient interface {
	OrderInfo(ctx context.Context, order string) (*AccrualOrderInfo, error)
}

// AccrualOrderInfo определяет формат ответа на получение информации о расчёте
// начислений баллов лояльности за совершённый заказ.
type AccrualOrderInfo struct {
	Order   string             `json:"order"`
	Status  AccrualOrderStatus `json:"status"`
	Accrual monetary.Unit      `json:"accrual"`
}

var (
	_ json.Marshaler   = (*AccrualOrderStatus)(nil)
	_ json.Unmarshaler = (*AccrualOrderStatus)(nil)
)

// OrderStatus определяет статус заказа в accrual.
type AccrualOrderStatus uint8

const (
	AccrualOrderStatusUnknown AccrualOrderStatus = iota

	AccrualOrderStatusRegistered
	AccrualOrderStatusInvalid
	AccrualOrderStatusProcessing
	AccrualOrderStatusProcessed
)

var accrualOrderStatusValues = []string{
	"UNKNOWN",
	"REGISTERED",
	"INVALID",
	"PROCESSING",
	"PROCESSED",
}

func (s AccrualOrderStatus) String() string {
	if s > 0 && int(s) < len(accrualOrderStatusValues) {
		return accrualOrderStatusValues[s]
	}
	return accrualOrderStatusValues[0]
}

func (s AccrualOrderStatus) MarshalJSON() ([]byte, error) {
	str := s.String()
	b := make([]byte, 0, len(`""`)+len(str))
	b = append(append(b, '"'), str...)
	return append(b, '"'), nil
}

func (s *AccrualOrderStatus) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("data must be a string")
	}
	target := string(data[1 : len(data)-1])
	for i := 0; i < len(accrualOrderStatusValues); i++ {
		if accrualOrderStatusValues[i] == target {
			*s = AccrualOrderStatus(i)
			return nil
		}
	}
	*s = AccrualOrderStatusUnknown
	return nil
}
