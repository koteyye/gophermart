package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"log/slog"

	"github.com/go-chi/chi/v5"

	"github.com/sergeizaitcev/gophermart/internal/accrual/service"
	"github.com/sergeizaitcev/gophermart/internal/accrual/storage"
)

// handler определяет HTTP-обработчик для accrual
// реализует интерфейс http.Handler
type handler struct {
	logger  *slog.Logger
	mux     *chi.Mux
	service *service.Service
}

// NewHandler возвращает новый экземпляр handler
func NewHandler(logger *slog.Logger, s *service.Service) http.Handler {
	r := &handler{
		logger:  logger,
		mux:     chi.NewRouter(),
		service: s,
	}
	r.init()
	return r
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *handler) init() {
	h.mux.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/orders", h.registerOrder)
			r.Post("/goods", h.createMatch)
			r.Get("/orders/{number}", h.getOrder)
		})
	})
}

func (h *handler) registerOrder(w http.ResponseWriter, r *http.Request) {
	o, err := parseOrder(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	err = h.service.CheckOrder(ctx, o.Number)
	if err != nil {
		mapErrorToResponse(w, err)
		return
	}

	go h.service.CreateOrder(&o)
	w.WriteHeader(http.StatusAccepted)
}

func (h *handler) createMatch(w http.ResponseWriter, r *http.Request) {
	m, err := parseMatch(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	err = h.service.CheckMatch(ctx, m.MatchName)
	if err != nil {
		if !errors.Is(err, storage.ErrNotFound) {
			mapErrorToResponse(w, err)
			return
		}
	}

	err = h.service.CreateMatch(ctx, &m)
	if err != nil {
		mapErrorToResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) getOrder(w http.ResponseWriter, r *http.Request) {
	orderNumber := chi.URLParam(r, "number")
	if orderNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	order, err := h.service.GetOrder(ctx, orderNumber)
	if err != nil {
		mapErrorToResponse(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(order)
}
