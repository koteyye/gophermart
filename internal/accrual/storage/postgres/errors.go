package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/sergeizaitcev/gophermart/internal/accrual/storage"
)

func errorHandle(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return fmt.Errorf("%w: %s", storage.ErrDuplicate, pgErr.Message)
		}
		return fmt.Errorf("%w: %s", storage.ErrOther, pgErr.Message)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %s", storage.ErrNotFound, err)
	}

	return fmt.Errorf("%s: %s", storage.ErrOther, err)
}
