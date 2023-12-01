package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
	"github.com/sergeizaitcev/gophermart/pkg/monetary"
)

var _ domain.OperationService = (*Operations)(nil)

// Operations определяет сервис обработки балансовых операций.
type Operations struct {
	db *sql.DB
}

// NewOperations возвращает новый экземпляр Operation.
func NewOperations(db *sql.DB) *Operations {
	return &Operations{db: db}
}

// Perform реализует интерфейс domain.OperationService.
func (o *Operations) Perform(ctx context.Context, operation domain.Operation) error {
	err := performOperation(ctx, o.db, operation)
	if err != nil {
		return fmt.Errorf("performing a balance operation: %w", err)
	}
	return nil
}

// GetOperations реализует интерфейс domain.OperationService.
func (o *Operations) GetOperations(
	ctx context.Context,
	id domain.UserID,
) ([]domain.Operation, error) {
	operations, err := getOperations(ctx, o.db, id)
	if err != nil {
		return nil, fmt.Errorf("operations search: %w", err)
	}
	return operations, nil
}

func performOperation(ctx context.Context, db *sql.DB, operation domain.Operation) error {
	query1 := "SELECT current_balance FROM users WHERE id = $1;"

	query2 := `INSERT INTO operations
		(user_created, order_number, amount)
	VALUES ($1, $2, $3);`

	query3 := `UPDATE users
	SET current_balance = current_balance - $1, withdrawn_balance = withdrawn_balance + $1
	WHERE id = $2;`

	return transaction(ctx, db, func(tx *sql.Tx) error {
		var currentBalance monetary.Unit

		err := tx.QueryRowContext(ctx, query1, operation.UserID).Scan(&currentBalance)
		if err != nil {
			return fmt.Errorf("user search: %w", errorHandling(err))
		}

		if currentBalance-operation.Sum < 0 {
			return domain.ErrBalanceBelowZero
		}

		_, err = tx.ExecContext(ctx, query2, operation.UserID, operation.OrderNumber, operation.Sum)
		if err != nil {
			return fmt.Errorf("creating a new balance operation: %w", errorHandling(err))
		}

		_, err = tx.ExecContext(ctx, query3, operation.Sum, operation.UserID)
		if err != nil {
			return fmt.Errorf("updating the user balance: %w", errorHandling(err))
		}

		return nil
	})
}

func getOperations(ctx context.Context, db *sql.DB, id domain.UserID) ([]domain.Operation, error) {
	query := `SELECT
		o.order_number, o.amount, o.created_at
	FROM operations AS o
		INNER JOIN users AS u ON o.user_created = u.id
	WHERE u.id = $1
	ORDER BY o.created_at;`

	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("operations search: %w", errorHandling(err))
	}
	defer rows.Close()

	var operations []domain.Operation

	for rows.Next() {
		operation := domain.Operation{UserID: id}

		err = rows.Scan(
			&operation.OrderNumber,
			&operation.Sum,
			&operation.ProcessedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("copying operation fields: %w", errorHandling(err))
		}

		operations = append(operations, operation)
	}

	err = rows.Err()
	if err != nil {
		return nil, errorHandling(err)
	}

	if len(operations) == 0 {
		return nil, domain.ErrNotFound
	}

	return operations, nil
}
