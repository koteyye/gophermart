package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sergeizaitcev/gophermart/pkg/luhn"
	"github.com/sergeizaitcev/gophermart/pkg/monetary"
	"github.com/sergeizaitcev/gophermart/pkg/strutil"
)

// Order определяет пользовательский заказ.
type Order struct {
	UserID     UserID        `json:"-"`
	Number     OrderNumber   `json:"number"`
	Status     OrderStatus   `json:"status"`
	Accrual    monetary.Unit `json:"accrual,omitempty"`
	UploadedAt time.Time     `json:"uploaded_at,omitempty"`
}

// IsEmpty возвращает true, если заказ пользователя пуст.
func (o Order) IsEmpty() bool {
	return o.Equal(Order{})
}

// Equal возвращает true, если заказ пользователя равен x.
func (o Order) Equal(x Order) bool {
	return o.UserID == x.UserID && o.Number == x.Number && o.Status == x.Status &&
		o.Accrual == x.Accrual &&
		o.UploadedAt.Equal(x.UploadedAt)
}

// OrderNumber определяет номер заказа.
type OrderNumber string

// NewOrderNumber конвертирует строку в номер заказа и возвращает его.
func NewOrderNumber(s string) (OrderNumber, error) {
	number := OrderNumber(s)
	err := number.Validate()
	if err != nil {
		return "", fmt.Errorf("order number validation: %w", err)
	}
	return number, nil
}

func (o OrderNumber) Validate() error {
	str := string(o)
	if !strutil.OnlyDigits(str) {
		return fmt.Errorf("order number must contain only digits: %q", str)
	}
	if !luhn.Check(str) {
		return fmt.Errorf("hash amount of the order number is invalid: %q", str)
	}
	return nil
}

var (
	_ json.Marshaler   = (*OrderStatus)(nil)
	_ json.Unmarshaler = (*OrderStatus)(nil)

	_ driver.Valuer = (*OrderStatus)(nil)
)

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
	"UNKNOWN",
	"NEW",
	"PROCESSING",
	"PROCESSED",
	"INVALID",
}

// NewOrderStatus конвертирует строку в статус заказа и возвращает его.
func NewOrderStatus(s string) (OrderStatus, error) {
	var status OrderStatus
	upper := strings.ToUpper(s)
	for i, v := range orderValues {
		if v == upper {
			status = OrderStatus(i)
			break
		}
	}
	if status == OrderStatusUnknown {
		return OrderStatusUnknown, fmt.Errorf("unknown order status: %q", s)
	}
	return status, nil
}

func (s OrderStatus) String() string {
	if s > 0 && int(s) < len(orderValues) {
		return orderValues[s]
	}
	return orderValues[0]
}

func (s OrderStatus) MarshalJSON() ([]byte, error) {
	str := s.String()
	b := make([]byte, 0, len(`""`)+len(str))
	b = append(append(b, '"'), str...)
	b = append(b, '"')
	return b, nil
}

func (s *OrderStatus) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("data must be a string")
	}
	str := string(data[1 : len(data)-1])
	status, err := NewOrderStatus(str)
	if err != nil {
		return err
	}
	*s = status
	return nil
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
	str := string(b)
	status, err := NewOrderStatus(str)
	if err != nil {
		return err
	}
	*s = status
	return nil
}

func (s OrderStatus) Value() (driver.Value, error) {
	if s == OrderStatusUnknown {
		return nil, nil
	}
	return s.String(), nil
}
