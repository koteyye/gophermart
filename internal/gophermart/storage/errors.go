package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
)

// ошибки storage
var (
	ErrDuplicate = errors.New("duplicate value")
	ErrNotFound  = errors.New("value not found")
	ErrOther     = errors.New("other storage error")
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
			return ErrDuplicate
		default:
			return ErrOther
		}
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	slog.Error(fmt.Sprintf("other storage error: %s", err))
	return ErrOther
}
