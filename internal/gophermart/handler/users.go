package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"log/slog"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
)

// getBalance возвращает баланс авторизованного пользователя.
func (h *handler) getBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := userFromContext(ctx)
	if userID == domain.EmptyUserID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	balance, err := h.users.GetBalance(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		slog.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(balance)
	if err != nil {
		slog.Error(err.Error())
	}
}
