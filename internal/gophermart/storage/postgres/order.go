package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/service"
	"github.com/sergeizaitcev/gophermart/pkg/monetary"
)

func (s *Storage) CreateOrder(ctx context.Context, userID uuid.UUID, order string) error {
	query1 := "SELECT user_created FROM orders WHERE number = $1;"
	query2 := "INSERT INTO orders (number, user_created) VALUES ($1, $2);"

	// Запускаем транзакцию, чтобы сначала проверить наличие в БД добавляемого
	// номера заказа и кто его добавил, а затем добавляем запись, если ее нет.
	return s.transaction(ctx, func(tx *sql.Tx) error {
		var creatorUserID uuid.UUID

		err := tx.QueryRowContext(ctx, query1, order).Scan(&creatorUserID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("order search: %w", errorHandling(err))
			}
		}

		// Если текущий пользователь или другой пользователь ранее загрузил
		// номер заказа, то возвращаем ошибку.
		if creatorUserID == userID {
			return service.ErrDuplicate
		} else if creatorUserID != uuid.Nil && creatorUserID != userID {
			return service.ErrDuplicateOtherUser
		}

		_, err = tx.ExecContext(ctx, query2, order, userID)
		if err != nil {
			return fmt.Errorf("creating a new order: %w", errorHandling(err))
		}

		return nil
	})
}

func (s *Storage) OrderStatus(ctx context.Context, order string) (service.OrderStatus, error) {
	var status service.OrderStatus

	query := "SELECT status FROM orders WHERE number = $1;"

	err := s.db.QueryRowContext(ctx, query, order).Scan(&status)
	if err != nil {
		return service.OrderStatusUnknown, fmt.Errorf("order search: %w", errorHandling(err))
	}

	return status, nil
}

func (s *Storage) Orders(ctx context.Context, userID uuid.UUID) ([]service.Order, error) {
	query := `SELECT
	number, status, accrual, created_at
	FROM orders
	WHERE user_created = $1;`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("orders search: %w", errorHandling(err))
	}
	defer rows.Close()

	var orders []service.Order

	for rows.Next() {
		var order service.Order

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
		return nil, err
	}

	if len(orders) == 0 {
		return nil, service.ErrNotFound
	}

	return orders, nil
}

func (s *Storage) UpdateOrder(
	ctx context.Context,
	order string,
	status service.OrderStatus,
	accrual monetary.Unit,
) error {
	query := `UPDATE orders
	SET status = $1, accrual = $2, updated_at = now()
	WHERE number = $3;`

	_, err := s.db.ExecContext(ctx, query, status, accrual, order)
	if err != nil {
		return fmt.Errorf("updating an order: %w", errorHandling(err))
	}

	return nil
}

func (s *Storage) UpdateOrderStatus(
	ctx context.Context,
	order string,
	status service.OrderStatus,
) error {
	query := `UPDATE orders
	SET status = $1, updated_at = now()
	WHERE number = $2;`

	_, err := s.db.ExecContext(ctx, query, status, order)
	if err != nil {
		return fmt.Errorf("updating an order status: %w", errorHandling(err))
	}

	return nil
}
