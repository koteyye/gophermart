package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sergeizaitcev/gophermart/internal/accrual/models"
	"github.com/sergeizaitcev/gophermart/pkg/luhn"
	"io"
)

var (
	ErrOrderEmpty   = errors.New("order is empty")
	ErrOrderInvalid = errors.New("order number invalid")
	ErrNoGoods      = errors.New("order doesnt contain goods")
)

func parseOrder(r io.Reader) (models.Order, error) {
	var o models.Order

	err := json.NewDecoder(r).Decode(&o)
	if err != nil {
		return models.Order{}, fmt.Errorf("decoding the user: %w", err)
	}
	if o.Number == "" {
		return models.Order{}, errors.New("order is empty")
	}
	if !luhn.Check(o.Number) {
		return models.Order{}, errors.New("order number invalid")
	}
	if len(o.Goods) == 0 {
		return models.Order{}, errors.New("order doesnt contain goods")
	}

	return o, nil
}
