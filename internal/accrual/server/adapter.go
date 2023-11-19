package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/sergeizaitcev/gophermart/internal/accrual/models"
	"github.com/sergeizaitcev/gophermart/pkg/luhn"
)

var (
	ErrOrderEmpty   = errors.New("order is empty")
	ErrOrderInvalid = errors.New("order number invalid")
	ErrNoGoods      = errors.New("order doesnt contain goods")
)

// parseOrder парсит запрос на регистрацию заказа и валидирует его
func parseOrder(r io.Reader) (models.Order, error) {
	var o models.Order

	err := json.NewDecoder(r).Decode(&o)
	if err != nil {
		return models.Order{}, fmt.Errorf("decoding the order: %w", err)
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

// parseMatch парсит запрос на создание вознаграждения за товар и валидирует его
func parseMatch(r io.Reader) (models.Match, error) {
	var m models.Match

	err := json.NewDecoder(r).Decode(&m)
	if err != nil {
		return models.Match{}, fmt.Errorf("decoding the match: %w", err)
	}
	if m.MatchName == "" && m.RewardType == "" {
		return models.Match{}, fmt.Errorf("match is empty")
	}
	if m.RewardType != "%" && m.RewardType != "pt" {
		return models.Match{}, fmt.Errorf("reward type invalid")
	}
	return m, nil
}

// mapErrorToResponse маппит ошибку на соответствующий код ответа
func mapErrorToResponse(w http.ResponseWriter, err error) {
	if errors.Is(err, models.ErrDuplicate) {
		w.WriteHeader(http.StatusConflict)
	}
	if errors.Is(err, models.ErrNotFound) {
		w.WriteHeader(http.StatusBadRequest)
	}
	if errors.Is(err, models.ErrOther) {
		w.WriteHeader(http.StatusInternalServerError)
	}
}