package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"log/slog"

	"github.com/google/uuid"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
	"github.com/sergeizaitcev/gophermart/pkg/queue"
)

var _ domain.OrderService = (*Orders)(nil)

// Orders определяет сервис обработки заказов пользователя.
type Orders struct {
	db      *sql.DB
	accrual domain.AccrualClient

	orders queue.FIFO[domain.Order]
	wg     *sync.WaitGroup
	termCh chan struct{}
}

// NewOrders возвращает новый экземпляр Order.
func NewOrders(db *sql.DB, accrual domain.AccrualClient) *Orders {
	o := &Orders{
		db:      db,
		accrual: accrual,
		wg:      &sync.WaitGroup{},
		termCh:  make(chan struct{}),
	}

	for i := 0; i < 2; i++ {
		o.wg.Add(1)
		go func() {
			defer o.wg.Done()
			o.processing()
		}()
	}

	return o
}

// Close сигнализирует всем процессам о завершении работы и блокируется
// до тех пор, пока они не будут завершены.
func (o *Orders) Close() {
	if o.closed() {
		return
	}
	close(o.termCh)
	o.wg.Wait()
}

func (o *Orders) closed() bool {
	select {
	case <-o.termCh:
		return true
	default:
		return false
	}
}

// GetOrders реализует интерфейс domain.OrderService.
func (o *Orders) GetOrders(ctx context.Context, id domain.UserID) ([]domain.Order, error) {
	orders, err := getOrdersByUser(ctx, o.db, id)
	if err != nil {
		return nil, fmt.Errorf("orders search: %w", err)
	}
	return orders, nil
}

// Process реализует интерфейс domain.OrderService.
func (o *Orders) Process(ctx context.Context, order domain.Order) error {
	err := createOrder(ctx, o.db, order)
	if err != nil {
		return fmt.Errorf("creating a new order: %w", err)
	}

	o.wg.Add(1)

	go func() {
		ctx, cancel := o.withCancel()
		defer cancel()

		_ = o.orders.Enqueue(ctx, order)
		o.wg.Done()
	}()

	return nil
}

func (o *Orders) withCancel() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	o.wg.Add(1)
	go func() {
		select {
		case <-ctx.Done():
		case <-o.termCh:
			cancel()
		}
		o.wg.Done()
	}()
	return ctx, cancel
}

func (o *Orders) processing() {
	ctx, cancel := o.withCancel()
	defer cancel()

	for {
		order, err := o.orders.Dequeue(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				break
			}
			slog.Debug(
				err.Error(),
				slog.String("scope", "queue"),
				slog.String("method", "dequeue"),
			)
			continue
		}

		order, err = o.tryProcessOrder(ctx, order)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				break
			}
			slog.Debug(err.Error(), slog.String("scope", "updating order"))
		}
		if order.IsEmpty() {
			continue
		}

		err = o.orders.Enqueue(ctx, order)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				break
			}
			slog.Debug(
				err.Error(),
				slog.String("scope", "queue"),
				slog.String("method", "enqueu"),
			)
		}
	}
}

func (o *Orders) tryProcessOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	info, err := o.accrual.GetAccrualInfo(ctx, order.Number)
	if err != nil {
		return order, err
	}

	switch info.Status {
	case domain.AccrualStatusUnknown:
		return order, nil
	case domain.AccrualStatusInvalid:
		order.Status = domain.OrderStatusInvalid
		return domain.Order{}, updateOrderStatus(ctx, o.db, order)
	case domain.AccrualStatusRegistered, domain.AccrualStatusProcessing:
		if order.Status == domain.OrderStatusNew {
			order.Status = domain.OrderStatusProcessing
			return order, updateOrderStatus(ctx, o.db, order)
		}
		return order, nil
	case domain.AccrualStatusProcessed:
		order.Status = domain.OrderStatusProcessed
		order.Accrual = info.Accrual
	}

	err = processOrder(ctx, o.db, order)
	if err != nil {
		return order, err
	}

	return domain.Order{}, nil
}

func getOrdersByUser(ctx context.Context, db *sql.DB, id domain.UserID) ([]domain.Order, error) {
	query := `SELECT
		order_number, status, accrual, created_at
	FROM orders
	WHERE user_created = $1;`

	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("orders search: %w", errorHandling(err))
	}
	defer rows.Close()

	var orders []domain.Order

	for rows.Next() {
		order := domain.Order{UserID: id}

		err = rows.Scan(
			&order.Number,
			&order.Status,
			&order.Accrual,
			&order.UploadedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("copying order fields: %w", errorHandling(err))
		}

		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		return nil, errorHandling(err)
	}

	if len(orders) == 0 {
		return nil, domain.ErrNotFound
	}

	return orders, nil
}

func createOrder(ctx context.Context, db *sql.DB, order domain.Order) error {
	query1 := "SELECT user_created FROM orders WHERE order_number = $1;"
	query2 := "INSERT INTO orders (order_number, user_created) VALUES ($1, $2);"

	// Запускаем транзакцию, чтобы сначала проверить наличие в БД добавляемого
	// номера заказа и кто его добавил, а затем добавляем запись, если ее нет.
	return transaction(ctx, db, func(tx *sql.Tx) error {
		var creatorUserID domain.UserID

		err := tx.QueryRowContext(ctx, query1, order.Number).Scan(&creatorUserID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("order search: %w", errorHandling(err))
			}
		}

		// Если текущий пользователь или другой пользователь ранее загрузил
		// номер заказа, то возвращаем ошибку.
		if creatorUserID == order.UserID {
			return domain.ErrDuplicate
		} else if creatorUserID != uuid.Nil && creatorUserID != order.UserID {
			return domain.ErrDuplicateOtherUser
		}

		_, err = tx.ExecContext(ctx, query2, order.Number, order.UserID)
		if err != nil {
			return fmt.Errorf("creating a new order: %w", errorHandling(err))
		}

		return nil
	})
}

func processOrder(ctx context.Context, db *sql.DB, order domain.Order) error {
	query1 := "UPDATE users SET current_balance = current_balance + $1 WHERE id = $2;"
	query2 := `UPDATE orders
	SET status = $1, accrual = $2, updated_at = now()
	WHERE order_number = $3;`

	return transaction(ctx, db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, query1, order.Accrual, order.UserID)
		if err != nil {
			return fmt.Errorf("updating a balance: %w", errorHandling(err))
		}

		_, err = tx.ExecContext(ctx, query2, order.Status, order.Accrual, order.Number)
		if err != nil {
			return fmt.Errorf("updating an order: %w", errorHandling(err))
		}

		return nil
	})
}

func updateOrderStatus(ctx context.Context, db *sql.DB, order domain.Order) error {
	query := "UPDATE orders SET status = $1, updated_at = now() WHERE order_number = $2;"

	_, err := db.ExecContext(ctx, query, order.Status, order.Number)
	if err != nil {
		return fmt.Errorf("updating an order status: %w", errorHandling(err))
	}

	return nil
}
