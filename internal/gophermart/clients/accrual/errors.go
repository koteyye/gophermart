package accrual

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
)

func prepareError(res *http.Response) error {
	switch res.StatusCode {
	case http.StatusNoContent:
		return domain.ErrOrderNotRegistered
	case http.StatusTooManyRequests:
		retryAfter, err := strconv.ParseInt(res.Header.Get("Retry-After"), 10, 64)
		if err != nil {
			return fmt.Errorf("parse int64: %w", err)
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("reading a request body: %w", err)
		}

		return &domain.ResourceExhaustedError{
			Message:    strings.ToLower(string(data)),
			RetryAfter: time.Duration(retryAfter) * time.Second,
		}
	default: // http.StatusInternalServerError
		return domain.ErrInternalServerError
	}
}
