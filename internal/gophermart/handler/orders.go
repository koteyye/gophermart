package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"log/slog"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
)

// orderProcess добавляет заказ авторизованного пользователя в обработку.
func (h *handler) orderProcess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := userFromContext(ctx)
	if userID == domain.EmptyUserID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !strings.Contains(r.Header.Get("Content-Type"), "text/plain") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orderNumber, err := domain.NewOrderNumber(string(b))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	order := domain.Order{
		UserID: userID,
		Number: orderNumber,
		Status: domain.OrderStatusNew,
	}

	err = h.orders.Process(ctx, order)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicate) {
			w.WriteHeader(http.StatusOK)
		} else if errors.Is(err, domain.ErrDuplicateOtherUser) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// getOrders возвращает все заказы авторизованного пользователя.
func (h *handler) getOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := userFromContext(ctx)
	if userID == domain.EmptyUserID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	orders, err := h.orders.GetOrders(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		slog.Error(err.Error())
	}
}
