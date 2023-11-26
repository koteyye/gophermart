package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/service"
	"github.com/sergeizaitcev/gophermart/pkg/luhn"
	"github.com/sergeizaitcev/gophermart/pkg/strutil"
)

func (h *handler) addOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := isAuthorized(ctx, w)
	if userID == uuid.Nil {
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

	order := string(b)
	if !strutil.OnlyDigits(order) || !luhn.Check(order) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = h.service.AddOrder(ctx, userID, order)
	if err != nil {
		if errors.Is(err, service.ErrDuplicate) {
			w.WriteHeader(http.StatusOK)
		} else if errors.Is(err, service.ErrDuplicateOtherUser) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *handler) orders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := isAuthorized(ctx, w)
	if userID == uuid.Nil {
		return
	}

	orders, err := h.service.Orders(ctx, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
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
		h.logger.Error(err.Error())
	}
}
