package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/service"
	"github.com/sergeizaitcev/gophermart/pkg/monetary"
)

func (s *Storage) CreateOperation(
	ctx context.Context,
	userID uuid.UUID,
	order string,
	amount monetary.Unit,
) (uuid.UUID, error) {
	query1 := "SELECT id, amount FROM balance WHERE user_id = $1;"

	query2 := `INSERT INTO operations
		(order_number, balance_id, amount)
	VALUES ($1, $2, $3)
	RETURNING id;`

	var id uuid.UUID

	err := s.transaction(ctx, func(tx *sql.Tx) error {
		var balanceID uuid.UUID
		var currentAmount monetary.Unit

		err := tx.QueryRowContext(ctx, query1, userID).Scan(&balanceID, &currentAmount)
		if err != nil {
			return fmt.Errorf("balance search: %w", errorHandling(err))
		}

		if currentAmount-amount < 0 {
			return service.ErrBalanceBelowZero
		}

		err = tx.QueryRowContext(ctx, query2, order, balanceID, amount).Scan(&id)
		if err != nil {
			return fmt.Errorf("creating a new operation: %w", errorHandling(err))
		}

		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *Storage) Operations(
	ctx context.Context,
	userID uuid.UUID,
) ([]service.Operation, error) {
	query := `SELECT
		o.order_number, o.amount, o.status, o.updated_at
	FROM operations AS o INNER JOIN balance AS b
		ON o.balance_id = b.id
	WHERE b.user_id = $1 AND o.status = $2
	ORDER BY o.updated_at;`

	rows, err := s.db.QueryContext(ctx, query, userID, service.OperationStatusDone)
	if err != nil {
		return nil, fmt.Errorf("operations search: %w", errorHandling(err))
	}
	defer rows.Close()

	var operations []service.Operation

	for rows.Next() {
		var operation service.Operation

		err = rows.Scan(
			&operation.Order,
			&operation.Sum,
			&operation.Status,
			&operation.ProcessedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("copying operation fields: %w", errorHandling(err))
		}

		operations = append(operations, operation)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if len(operations) == 0 {
		return nil, service.ErrNotFound
	}

	return operations, nil
}

func (s *Storage) UpdateOperationStatus(
	ctx context.Context,
	order string,
	status service.OperationStatus,
) error {
	query := `UPDATE operations
	SET status = $1, updated_at = now()
	WHERE order_number = $2;`

	_, err := s.db.ExecContext(ctx, query, status, order)
	if err != nil {
		return fmt.Errorf("updating an operation state: %w", errorHandling(err))
	}

	return nil
}

func (s *Storage) PerformOperation(ctx context.Context, operationID uuid.UUID) error {
	query1 := `SELECT
		b.id, b.amount, o.amount
	FROM operations AS o INNER JOIN balance AS b
		ON o.balance_id = b.id
	WHERE id = $1;`

	query2 := `UPDATE balance
	SET amount = amount - $1, withdrawn = withdrawn + $1, updated_at = now()
	WHERE id = $2;`

	query3 := `UPDATE operations
	SET status = $1, updated_at = now()
	WHERE id = $2;`

	return s.transaction(ctx, func(tx *sql.Tx) error {
		var balanceID uuid.UUID
		var currentAmount monetary.Unit
		var amount monetary.Unit

		err := tx.QueryRowContext(ctx, query1, operationID).
			Scan(&balanceID, &currentAmount, &amount)
		if err != nil {
			return fmt.Errorf("balance search: %w", errorHandling(err))
		}

		if currentAmount-amount < 0 {
			return service.ErrBalanceBelowZero
		}

		_, err = tx.ExecContext(ctx, query2, balanceID)
		if err != nil {
			return fmt.Errorf("updating a balance: %w", errorHandling(err))
		}

		_, err = tx.ExecContext(ctx, query3, service.OperationStatusDone, operationID)
		if err != nil {
			return fmt.Errorf("updating a balance: %w", errorHandling(err))
		}

		return nil
	})
}

func (s *Storage) BalanceIncrement(ctx context.Context, order string) error {
	query1 := "select accrual, user_created from orders where number = $1"
	query2 := "update balance set amount = amount + $1 where user_id = $2"

	return s.transaction(ctx, func(tx *sql.Tx) error {
		var orderAmount monetary.Unit
		var userID uuid.UUID

		err := tx.QueryRowContext(ctx, query1, order).Scan(&orderAmount, &userID)
		if err != nil {
			return fmt.Errorf("order search: %w", errorHandling(err))
		}

		_, err = tx.ExecContext(ctx, query2, orderAmount, userID)
		if err != nil {
			return fmt.Errorf("updating a balance: %w", errorHandling(err))
		}

		return nil
	})
}