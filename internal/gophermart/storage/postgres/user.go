package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/service"
	"github.com/sergeizaitcev/gophermart/pkg/monetary"
	"github.com/sergeizaitcev/gophermart/pkg/passwords"
)

func (s *Storage) CreateUser(
	ctx context.Context,
	u service.User,
) (uuid.UUID, error) {
	var userID uuid.UUID

	hashedPassword, err := passwords.Hash(u.Password)
	if err != nil {
		return uuid.Nil, fmt.Errorf("password hashing: %w", err)
	}

	query1 := `INSERT INTO users
		(login, hashed_password)
	VALUES ($1, $2)
	RETURNING id;`

	query2 := "INSERT INTO balance (user_id) VALUES ($1);"

	err = s.transaction(ctx, func(tx *sql.Tx) error {
		err := tx.QueryRowContext(ctx, query1, u.Login, hashedPassword).Scan(&userID)
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

func (s *Storage) UserID(
	ctx context.Context,
	u service.User,
) (uuid.UUID, error) {
	var (
		userID         uuid.UUID
		hashedPassword string
	)

	query := "SELECT id, hashed_password FROM users WHERE login = $1;"

	err := s.db.QueryRowContext(ctx, query, u.Login).Scan(&userID, &hashedPassword)
	if err != nil {
		return uuid.Nil, fmt.Errorf("user search: %w", errorHandling(err))
	}

	if !passwords.Compare(hashedPassword, u.Password) {
		return uuid.Nil, service.ErrInvalidPassword
	}

	return userID, nil
}

func (s *Storage) UserExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	var count int

	query := "SELECT COUNT(*) FROM users WHERE id = $1;"

	err := s.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("user search: %w", errorHandling(err))
	}

	return count == 1, nil
}

func (s *Storage) Balance(ctx context.Context, userID uuid.UUID) (*service.Balance, error) {
	var balance service.Balance

	query := `SELECT
		amount, withdrawn
	FROM balance
	WHERE user_id = $1 AND deleted_at IS NULL;`

	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&balance.Current,
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
			return service.ErrBalanceBelowZero
		}

		_, err = tx.ExecContext(ctx, query2, amount, userID)
		if err != nil {
			return fmt.Errorf("decrementing the balance: %w", errorHandling(err))
		}

		return nil
	})
}
