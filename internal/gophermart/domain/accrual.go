package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/sergeizaitcev/gophermart/pkg/monetary"
)

// AccrualInfo определяет информацию о начислении баллов лояльности
// за совершённый заказ.
type AccrualInfo struct {
	OrderNumber OrderNumber   `json:"order"`
	Status      AccrualStatus `json:"status"`
	Accrual     monetary.Unit `json:"accrual"`
}

var (
	_ json.Marshaler   = (*AccrualStatus)(nil)
	_ json.Unmarshaler = (*AccrualStatus)(nil)
)

// AccrualStatus определяет статус начисления баллов лояльности за совершённый
// заказ.
type AccrualStatus uint8

const (
	AccrualStatusUnknown AccrualStatus = iota

	AccrualStatusRegistered
	AccrualStatusInvalid
	AccrualStatusProcessing
	AccrualStatusProcessed
)

var accrualStatusValues = []string{
	"UNKNOWN",
	"REGISTERED",
	"INVALID",
	"PROCESSING",
	"PROCESSED",
}

// NewAccrualStatus конвертирует строку в статус начисления баллов лояльности
// и возвращает его.
func NewAccrualStatus(s string) (AccrualStatus, error) {
	var status AccrualStatus
	upper := strings.ToUpper(s)
	for i, v := range accrualStatusValues {
		if v == upper {
			status = AccrualStatus(i)
			break
		}
	}
	if status == AccrualStatusUnknown {
		return AccrualStatusUnknown, fmt.Errorf("unknown accrual status: %q", s)
	}
	return status, nil
}

func (s AccrualStatus) String() string {
	if s > 0 && int(s) < len(accrualStatusValues) {
		return accrualStatusValues[s]
	}
	return accrualStatusValues[0]
}

func (s AccrualStatus) MarshalJSON() ([]byte, error) {
	str := s.String()
	b := make([]byte, 0, len(`""`)+len(str))
	b = append(append(b, '"'), str...)
	b = append(b, '"')
	return b, nil
}

func (s *AccrualStatus) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("data must be a string")
	}
	str := string(data[1 : len(data)-1])
	status, err := NewAccrualStatus(str)
	if err != nil {
		return err
	}
	*s = status
	return nil
}
