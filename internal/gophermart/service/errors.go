package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
)

const integrityConstraintViolationClass = "23"

func errorHandling(err error) error {
	var pqErr *pq.Error
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %s", domain.ErrNotFound, err)
	}
	if errors.As(err, &pqErr) {
		if pqErr.Code.Class() == integrityConstraintViolationClass {
			return fmt.Errorf("%w: %s", domain.ErrDuplicate, pqErr.Message)
		}
	}
	return err
}
