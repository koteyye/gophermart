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

// OrderQueue определяет очередь заказов по принципу FIFO.
type OrderQueue = queue.FIFO[string]

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
	str := strings.ToUpper(string(b))

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
	Number     string        `json:"number"`
	Status     OrderStatus   `json:"status"`
	Accrual    monetary.Unit `json:"accrual,omitempty"`
	UploadedAt time.Time     `json:"uploaded_at"`
}

// Orders возвращает заказы пользователя.
func (s *Service) Orders(ctx context.Context, userID uuid.UUID) ([]Order, error) {
	values, err := s.storage.Orders(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("orders search: %w", err)
	}
	return values, nil
}

// AddOrder добавляет заказ пользователя в обработку.
func (s *Service) AddOrder(ctx context.Context, userID uuid.UUID, order string) error {
	err := s.storage.CreateOrder(ctx, userID, order)
	if err != nil {
		return fmt.Errorf("creating a new order: %w", err)
	}

	s.wg.Add(1)

	go func() {
		ctx, cancel := s.withCancel()
		defer cancel()
		s.orders.Enqueue(ctx, order)
		s.wg.Done()
	}()

	return nil
}

func (s *Service) orderProcessing() {
	ctx, cancel := s.withCancel()
	defer cancel()

	for {
		err := s.orderTransaction(ctx, func(order string) (bool, error) {
			status, err := s.storage.OrderStatus(ctx, order)
			if err != nil {
				if errors.Is(err, ErrNotFound) {
					return false, err
				}
				return true, err
			}

			fsm := orderFSM{order: order, service: s}

			rollback, err := fsm.Event(ctx, status)
			if err != nil {
				var exhausted *ResourceExhaustedError

				if errors.As(err, &exhausted) {
					s.logger.Debug(exhausted.Message)
					select {
					case <-ctx.Done():
						return false, ctx.Err()
					case <-time.After(exhausted.RetryAfter):
					}
				}
			}

			return rollback, err
		})
		if err != nil {
			break
		}
	}
}

func (s *Service) orderTransaction(ctx context.Context, fn func(order string) (bool, error)) error {
	order, err := s.orders.Dequeue(ctx)
	if err != nil {
		return err
	}

	rollback, err := fn(order)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return err
		}
		s.logger.Error(err.Error())
	}

	if rollback {
		return s.orders.Enqueue(ctx, order)
	}

	return nil
}

type orderFSM struct {
	order   string
	service *Service
}

func (fsm orderFSM) Event(ctx context.Context, status OrderStatus) (rollback bool, err error) {
	switch status {
	case OrderStatusNew:
		return fsm.processing(ctx)
	case OrderStatusProcessing:
		return fsm.continueProcessing(ctx)
	}
	return false, nil
}

func (fsm orderFSM) processing(ctx context.Context) (rollback bool, err error) {
	info, err := fsm.service.accrual.OrderInfo(ctx, fsm.order)
	if err != nil {
		return true, err
	}

	fsm.service.logger.Info(fmt.Sprintf("%+v", info))

	switch info.Status {
	case AccrualOrderStatusInvalid:
		err := fsm.service.storage.UpdateOrderStatus(ctx, fsm.order, OrderStatusInvalid)
		return false, err
	case AccrualOrderStatusRegistered, AccrualOrderStatusProcessing:
		err := fsm.service.storage.UpdateOrderStatus(ctx, fsm.order, OrderStatusProcessing)
		return true, err
	}
	err = fsm.service.storage.UpdateOrder(ctx, fsm.order, OrderStatusProcessed, info.Accrual)
	if err != nil {
		return true, err
	}

	err = fsm.service.storage.BalanceIncrement(ctx, fsm.order)
	if err != nil {
		return true, err
	}

	return false, nil
}

func (fsm orderFSM) continueProcessing(ctx context.Context) (rollback bool, err error) {
	info, err := fsm.service.accrual.OrderInfo(ctx, fsm.order)
	if err != nil {
		return true, err
	}

	switch info.Status {
	case AccrualOrderStatusInvalid:
		err := fsm.service.storage.UpdateOrderStatus(ctx, fsm.order, OrderStatusInvalid)
		return false, err
	case AccrualOrderStatusRegistered, AccrualOrderStatusProcessing:
		return true, nil
	}

	err = fsm.service.storage.UpdateOrder(ctx, fsm.order, OrderStatusProcessed, info.Accrual)
	if err != nil {
		return true, err
	}

	return false, nil
}
