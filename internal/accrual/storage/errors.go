package storage

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sergeizaitcev/gophermart/internal/accrual/models"
)

// обрабатываемые ошибки pgx
const (
	PqDuplicateErr = "23505"
)

func mapStorageErr(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case PqDuplicateErr:
			return fmt.Errorf("%w: %s", models.ErrDuplicate, err)
		default:
			return fmt.Errorf("%w: %s", models.ErrOther, err)
		}
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: %s", models.ErrNotFound, err)
	}
	return fmt.Errorf("%s: %s", models.ErrOther, err)
}
