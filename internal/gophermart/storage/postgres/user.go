package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/storage"
	"github.com/sergeizaitcev/gophermart/pkg/monetary"
	"github.com/sergeizaitcev/gophermart/pkg/passwords"
)

func (s *Storage) CreateUser(
	ctx context.Context,
	login string,
	password string,
) (uuid.UUID, error) {
	var userID uuid.UUID

	hashedPassword, err := passwords.Hash(password)
	if err != nil {
		return uuid.Nil, fmt.Errorf("password hashing: %w", err)
	}

	query1 := `INSERT INTO users
		(login, hashed_password)
	VALUES ($1, $2)
	RETURNING id;`

	query2 := "INSERT INTO balance (user_id) VALUES ($1);"

	err = s.transaction(ctx, func(tx *sql.Tx) error {
		err := tx.QueryRowContext(ctx, query1, login, hashedPassword).Scan(&userID)
		if err != nil {
			return fmt.Errorf("creating a new user: %w", err)
		}

		_, err = tx.ExecContext(ctx, query2, userID)
		if err != nil {
			return fmt.Errorf("creating a user balance: %w", errorHandling(err))
		}

		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (s *Storage) GetUser(
	ctx context.Context,
	login string,
	password string,
) (uuid.UUID, error) {
	var (
		userID         uuid.UUID
		hashedPassword string
	)

	query := "SELECT id, hashed_password FROM users WHERE login = $1;"

	err := s.db.QueryRowContext(ctx, query, login).Scan(&userID, &hashedPassword)
	if err != nil {
		return uuid.Nil, fmt.Errorf("user search: %w", errorHandling(err))
	}

	if !passwords.Compare(hashedPassword, password) {
		return uuid.Nil, storage.ErrInvalidPassword
	}

	return userID, nil
}

func (s *Storage) GetLogin(ctx context.Context, userID uuid.UUID) (string, error) {
	var login string

	query := "SELECT login FROM users WHERE id = $1;"

	err := s.db.QueryRowContext(ctx, query, userID).Scan(&login)
	if err != nil {
		return "", fmt.Errorf("user search: %w", errorHandling(err))
	}

	return login, nil
}

func (s *Storage) GetBalance(ctx context.Context, userID uuid.UUID) (*storage.UserBalance, error) {
	var balance storage.UserBalance

	query := `SELECT
		user_id, amount, withdrawn
	FROM balance
	WHERE user_id = $1 AND deleted_at IS NULL;`

	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&balance.UserID,
		&balance.Amount,
		&balance.Withdrawn,
	)
	if err != nil {
		return nil, fmt.Errorf("balance search: %w", err)
	}

	return &balance, nil
}

func (s *Storage) UpdateBalance(ctx context.Context, userID uuid.UUID, amount monetary.Unit) error {
	if amount == 0 {
		return nil
	}
	if amount > 0 {
		return s.incrementBalance(ctx, userID, amount)
	}
	return s.decrementBalance(ctx, userID, -amount)
}

func (s *Storage) incrementBalance(
	ctx context.Context,
	userID uuid.UUID,
	amount monetary.Unit,
) error {
	query := `UPDATE balance
	SET amount = amount + $1, updated_at = now()
	WHERE user_id = $2;`

	_, err := s.db.ExecContext(ctx, query, amount, userID)
	if err != nil {
		return fmt.Errorf("increasing the balance: %w", errorHandling(err))
	}

	return nil
}

func (s *Storage) decrementBalance(
	ctx context.Context,
	userID uuid.UUID,
	amount monetary.Unit,
) error {
	query1 := "SELECT amount FROM balance WHERE user_id = $1;"

	query2 := `UPDATE balance
	SET amount = amount - $1, withdrawn = withdrawn + $1, updated_at = now()
	WHERE user_id = $2;`

	return s.transaction(ctx, func(tx *sql.Tx) error {
		var currentAmount monetary.Unit

		err := tx.QueryRowContext(ctx, query1, userID).Scan(&currentAmount)
		if err != nil {
			return fmt.Errorf("balance search: %w", err)
		}

		if currentAmount-amount < 0 {
			return storage.ErrBalanceBelowZero
		}

		_, err = tx.ExecContext(ctx, query2, amount, userID)
		if err != nil {
			return fmt.Errorf("decrementing the balance: %w", errorHandling(err))
		}

		return nil
	})
}
