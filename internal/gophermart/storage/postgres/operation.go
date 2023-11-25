package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/storage"
	"github.com/sergeizaitcev/gophermart/pkg/monetary"
)

func (s *Storage) CreateOperation(
	ctx context.Context,
	userID uuid.UUID,
	order string,
	amount monetary.Unit,
) error {
	query1 := "SELECT id FROM balance WHERE user_id = $1;"
	query2 := `INSERT INTO operations
		(order_number, balance_id, amount)
	VALUES ($1, $2, $3);`

	return s.transaction(ctx, func(tx *sql.Tx) error {
		var balanceID uuid.UUID

		err := tx.QueryRowContext(ctx, query1, userID).Scan(&balanceID)
		if err != nil {
			return fmt.Errorf("balance search: %w", errorHandling(err))
		}

		_, err = tx.ExecContext(ctx, query2, order, balanceID, amount)
		if err != nil {
			return fmt.Errorf("creating a new operation: %w", errorHandling(err))
		}

		return nil
	})
}

func (s *Storage) GetOperation(ctx context.Context, order string) (*storage.Operation, error) {
	var operation storage.Operation

	query := `SELECT
		b.user_id, o.order_number, o.amount, o.status, o.updated_at
	FROM operations AS o INNER JOIN balance AS b
		ON o.balance_id = b.id
	WHERE o.order_number = $1 AND o.deleted_at IS NULL;`

	err := s.db.QueryRowContext(ctx, query, order).Scan(
		&operation.UserID,
		&operation.OrderNumber,
		&operation.Amount,
		&operation.Status,
		&operation.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("operation search: %w", errorHandling(err))
	}

	return &operation, nil
}

func (s *Storage) GetOperations(
	ctx context.Context,
	userID uuid.UUID,
) ([]storage.Operation, error) {
	query := `SELECT
		b.user_id, o.order_number, o.amount, o.status, o.updated_at
	FROM operations AS o INNER JOIN balance AS b
		ON o.balance_id = b.id
	WHERE b.user_id = $1 AND b.deleted_at IS NULL
	ORDER BY o.updated_at;`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("operations search: %w", errorHandling(err))
	}
	defer rows.Close()

	var operations []storage.Operation

	for rows.Next() {
		var operation storage.Operation

		err = rows.Scan(
			&operation.UserID,
			&operation.OrderNumber,
			&operation.Amount,
			&operation.Status,
			&operation.UpdatedAt,
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
		return nil, storage.ErrNotFound
	}

	return operations, nil
}

func (s *Storage) UpdateOperationStatus(
	ctx context.Context,
	order string,
	status storage.OperationStatus,
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

func (s *Storage) DeleteOperation(ctx context.Context, order string) error {
	query := `UPDATE operations
	SET deleted_at = now()
	WHERE order_number = $1;`

	_, err := s.db.ExecContext(ctx, query, order)
	if err != nil {
		return fmt.Errorf("deleting an operation: %w", errorHandling(err))
	}

	return nil
}
