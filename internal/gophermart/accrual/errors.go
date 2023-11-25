package accrual

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrOrderNotRegistered возвращается, если номер заказа не зарегистрирован.
	ErrOrderNotRegistered = errors.New("order is not registered")

	// ErrInternalServerError возвращается, если сервер вернул 500 код ответа.
	ErrInternalServerError = errors.New("internal server error")
)

// ResourceExhaustedError возвращается, если клиент превысил лимит запросов
// в минуту.
type ResourceExhaustedError struct {
	msg        string
	retryAfter time.Duration
}

func (err *ResourceExhaustedError) Error() string {
	return err.msg
}

func (err *ResourceExhaustedError) RetryAfter() time.Duration {
	return err.retryAfter
}

func prepareError(res *http.Response) error {
	switch res.StatusCode {
	case http.StatusNoContent:
		return ErrOrderNotRegistered
	case http.StatusTooManyRequests:
		retryAfter, err := strconv.ParseInt(res.Header.Get("Retry-After"), 10, 64)
		if err != nil {
			return fmt.Errorf("parse int64: %w", err)
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("reading a request body: %w", err)
		}

		return &ResourceExhaustedError{
			msg:        strings.ToLower(string(data)),
			retryAfter: time.Duration(retryAfter) * time.Second,
		}
	default: // http.StatusInternalServerError
		return ErrInternalServerError
	}
}
