package handlers

import (
	"net/http"

	"log/slog"

	"github.com/go-chi/chi/v5"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/service"
)

// handler определяет HTTP-обработчик для gophermart;
// реализует интерфейс http.Handler.
type handler struct {
	logger  *slog.Logger
	mux     *chi.Mux
	service *service.Service
}

// NewHandler возвращает новый экземпляр handler.
func NewHandler(logger *slog.Logger, service *service.Service) http.Handler {
	r := &handler{
		logger:  logger,
		mux:     chi.NewRouter(),
		service: service,
	}
	r.init()
	return r
}

// ServeHTTP реализует интерфейс http.Handler.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *handler) init() {
	h.mux.Route("/api/user", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/register", h.register)
			r.Post("/login", h.login)
		})

		r.Group(func(r chi.Router) {
			r.Use(h.auth)

			r.Post("/orders", h.addOrder)
			r.Get("/orders", h.orders)

			r.Get("/balance", h.balance)
			r.Post("/balance/withdraw", h.balanceWithdraw)
			r.Get("/withdrawals", h.withdrawals)
		})
	})
}
