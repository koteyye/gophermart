package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"log/slog"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
)

// operationPerform выполняет балансовую операцию авторизованного пользователя.
func (h *handler) operationPerform(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := userFromContext(ctx)
	if userID == domain.EmptyUserID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var operation domain.Operation

	err := json.NewDecoder(r.Body).Decode(&operation)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	err = operation.OrderNumber.Validate()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		slog.Error(err.Error())
		return
	}

	operation.UserID = userID

	err = h.operations.Perform(ctx, operation)
	if err != nil {
		if errors.Is(err, domain.ErrBalanceBelowZero) {
			w.WriteHeader(http.StatusPaymentRequired)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		slog.Error(err.Error())
	}
}

// getOperations возвращает все балансовые операции авторизованного пользователя.
func (h *handler) getOperations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := userFromContext(ctx)
	if userID == domain.EmptyUserID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	operations, err := h.operations.GetOperations(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		slog.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(operations)
	if err != nil {
		slog.Error(err.Error())
	}
}
