package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/sergeizaitcev/gophermart/internal/accrual/models"
)

func errorHandle(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return fmt.Errorf("%w: %s", models.ErrDuplicate, pgErr.Message)
		}
		return fmt.Errorf("%w: %s", models.ErrOther, pgErr.Message)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %s", models.ErrNotFound, err)
	}

	return fmt.Errorf("%s: %s", models.ErrOther, err)
}