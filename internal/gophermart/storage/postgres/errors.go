package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/service"
)

const integrityConstraintViolationClass = "23"

func errorHandling(err error) error {
	var pqErr *pq.Error

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %s", service.ErrNotFound, err)
	}

	if errors.As(err, &pqErr) {
		if pqErr.Code.Class() == integrityConstraintViolationClass {
			return fmt.Errorf("%w: %s", service.ErrDuplicate, pqErr.Message)
		}
	}

	return fmt.Errorf("%s: %s", service.ErrOther, err)
}
