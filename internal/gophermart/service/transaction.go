package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// transaction создает новую SQL-транзакцию.
func transaction(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	var tx *sql.Tx

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		rErr := tx.Rollback()
		if err == nil && rErr != nil && !errors.Is(rErr, sql.ErrTxDone) {
			err = fmt.Errorf("rollback transaction: %w", err)
		}
	}()

	err = fn(tx)
	if err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
