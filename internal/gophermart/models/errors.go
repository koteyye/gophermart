package models

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ошибки storage
var (
	ErrDuplicate          = errors.New("duplicate value")
	ErrDuplicateOtherUser = errors.New("duplicate value from other user")
	ErrNotFound           = errors.New("value not found")
	ErrOther              = errors.New("other storage error")
	ErrBalanceBelowZero   = errors.New("balance can't be below zero")
)

// обрабатываемые ошибки pgx
const (
	PqDuplicateErr = "23505"
)

func MapStorageErr(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case PqDuplicateErr:
			return fmt.Errorf("%w: %s", ErrDuplicate, err)
		default:
			return fmt.Errorf("%w: %s", ErrOther, err)
		}
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: %s", ErrNotFound, err)
	}
	return fmt.Errorf("%s: %s", ErrOther, err)
}
