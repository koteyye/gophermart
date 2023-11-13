package accrual

import (
	"encoding/json"
	"errors"

	"github.com/sergeizaitcev/gophermart/pkg/monetary"
)

// OrderInfo определяет формат ответа на получение информации о расчёте
// начислений баллов лояльности за совершённый заказ.
type OrderInfo struct {
	Order   string        `json:"order"`
	Status  OrderStatus   `json:"status"`
	Accrual monetary.Unit `json:"accrual"`
}

// IsEmpty возвращает true, если информация о расчёте начисления пуста.
func (o OrderInfo) IsEmpty() bool {
	return o.Order == "" && o.Status == StatusUnknown && o.Accrual == 0
}

var (
	_ json.Marshaler   = (*OrderStatus)(nil)
	_ json.Unmarshaler = (*OrderStatus)(nil)
)

// OrderStatus определяет статус заказа.
type OrderStatus uint8

const (
	StatusUnknown OrderStatus = iota

	StatusRegistered
	StatusInvalid
	StatusPrecessing
	StatusProcessed
)

var statusValues = []string{
	"UNKNOWN",
	"REGISTERED",
	"INVALID",
	"PROCESSING",
	"PROCESSED",
}

func (s OrderStatus) String() string {
	if s >= 0 && int(s) < len(statusValues) {
		return statusValues[s]
	}
	return statusValues[0]
}

func (s OrderStatus) MarshalJSON() ([]byte, error) {
	str := s.String()
	b := make([]byte, 0, len(`""`)+len(str))
	b = append(append(b, '"'), str...)
	return append(b, '"'), nil
}

func (s *OrderStatus) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("data must be a string")
	}
	target := string(data[1 : len(data)-1])
	for i := 0; i < len(statusValues); i++ {
		if statusValues[i] == target {
			*s = OrderStatus(i)
			return nil
		}
	}
	*s = StatusUnknown
	return nil
}
