package domain

import (
	"time"

	"github.com/sergeizaitcev/gophermart/pkg/monetary"
)

// Operation определяет балансовую операцию.
type Operation struct {
	UserID      UserID        `json:"-"`
	OrderNumber OrderNumber   `json:"order"`
	Sum         monetary.Unit `json:"sum"`
	ProcessedAt time.Time     `json:"processed_at,omitempty"`
}
