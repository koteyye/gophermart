package service

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/sergeizaitcev/gophermart/pkg/monetary"
	"github.com/sergeizaitcev/gophermart/pkg/queue"
)

// OperationQueue определяет очередь с балансовыми операциями.
type OperationQueue = queue.FIFO[Operation]

// Balance определяет баланс пользователя.
type Balance struct {
	Current   monetary.Unit `json:"current"`
	Withdrawn monetary.Unit `json:"withdrawn"`
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
	"UNKNOWN",
	"RUN",
	"DONE",
	"ERROR",
}

func (s OperationStatus) String() string {
	if s > 0 && int(s) < len(operationStatusValues) {
		return operationStatusValues[s]
	}
	return operationStatusValues[0]
}

func (s OperationStatus) MarshalJSON() ([]byte, error) {
	str := s.String()
	b := make([]byte, len(`""`)+len(str))
	b = append(append(b, '"'), str...)
	b = append(b, '"')
	return b, nil
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
	low := strings.ToUpper(string(b))

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
	ID          uuid.UUID       `json:"-"`
	UserID      uuid.UUID       `json:"-"`
	Order       string          `json:"order"`
	Sum         monetary.Unit   `json:"sum"`
	Status      OperationStatus `json:"-"`
	ProcessedAt time.Time       `json:"processed_at,omitempty"`
}

// Balance возвращает текущий баланс пользователя.
func (s *Service) Balance(ctx context.Context, userID uuid.UUID) (*Balance, error) {
	b, err := s.storage.Balance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("balance search: %w", err)
	}
	return b, nil
}

// Withdrawals возвращает список списаний пользователя.
func (s *Service) Withdrawals(ctx context.Context, userID uuid.UUID) ([]Operation, error) {
	operations, err := s.storage.Operations(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("withdrawals search: %w", err)
	}
	return operations, nil
}

// Withdrawn списывает средства с баланса пользователя по номеру заказа.
func (s *Service) Withdraw(
	ctx context.Context,
	userID uuid.UUID,
	order string,
	amount monetary.Unit,
) error {
	operationID, err := s.storage.CreateOperation(ctx, userID, order, amount)
	if err != nil {
		return fmt.Errorf("withdrawal of the amount from the balance: %w", err)
	}

	s.wg.Add(1)

	go func() {
		ctx, cancel := s.withCancel()
		defer cancel()
		s.operations.Enqueue(ctx, Operation{ID: operationID, UserID: userID, Order: order})
		s.wg.Done()
	}()

	return nil
}

func (s *Service) operationProcessing() {
	ctx, cancel := s.withCancel()
	defer cancel()

	for {
		err := s.operationTransaction(ctx, func(operation Operation) (bool, error) {
			err := s.storage.PerformOperation(ctx, operation.ID)
			return err != nil, err
		})
		if err != nil {
			break
		}
	}
}

func (s *Service) operationTransaction(
	ctx context.Context,
	fn func(operation Operation) (bool, error),
) error {
	operation, err := s.operations.Dequeue(ctx)
	if err != nil {
		return err
	}

	rollback, err := fn(operation)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return err
		}
		s.logger.Error(err.Error())
	}

	if rollback {
		return s.operations.Enqueue(ctx, operation)
	}

	return nil
}
