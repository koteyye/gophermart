package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/service"
	"github.com/sergeizaitcev/gophermart/pkg/luhn"
	"github.com/sergeizaitcev/gophermart/pkg/strutil"
)

func (h *handler) register(w http.ResponseWriter, r *http.Request) {
	var u service.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err.Error())
		return
	}

	ctx := r.Context()

	token, err := h.service.SignUp(ctx, u)
	if err != nil {
		if errors.Is(err, service.ErrDuplicate) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		h.logger.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, token)
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	var u service.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err.Error())
		return
	}

	ctx := r.Context()

	token, err := h.service.SignIn(ctx, u)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		h.logger.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, token)
}

func (h *handler) balance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := isAuthorized(ctx, w)
	if userID == uuid.Nil {
		return
	}

	balance, err := h.service.Balance(ctx, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		h.logger.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(balance)
	if err != nil {
		h.logger.Error(err.Error())
	}
}

func (h *handler) balanceWithdraw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := isAuthorized(ctx, w)
	if userID == uuid.Nil {
		return
	}

	var operation service.Operation

	err := json.NewDecoder(r.Body).Decode(&operation)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !strutil.OnlyDigits(operation.Order) || !luhn.Check(operation.Order) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if operation.Sum < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.service.Withdraw(ctx, userID, operation.Order, operation.Sum)
	if err != nil {
		if errors.Is(err, service.ErrBalanceBelowZero) {
			w.WriteHeader(http.StatusPaymentRequired)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		h.logger.Error(err.Error())
	}
}

func (h *handler) withdrawals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := isAuthorized(ctx, w)
	if userID == uuid.Nil {
		return
	}

	operations, err := h.service.Withdrawals(ctx, userID)
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

	err = json.NewEncoder(w).Encode(operations)
	if err != nil {
		h.logger.Error(err.Error())
	}
}
