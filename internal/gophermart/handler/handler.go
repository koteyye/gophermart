package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
	"github.com/sergeizaitcev/gophermart/pkg/sign"
)

// HandlerOptions определяет опции для HTTP-обработчика.
type HandlerOptions struct {
	Auth       domain.AuthService
	Operations domain.OperationService
	Orders     domain.OrderService
	Users      domain.UserService
	Signer     sign.Signer
}

// handler определяет HTTP-обработчик для gophermart.
type handler struct {
	mux    *chi.Mux
	signer sign.Signer

	auth       domain.AuthService
	operations domain.OperationService
	orders     domain.OrderService
	users      domain.UserService
}

// New возвращает новый HTTP-обработчик.
func New(opt HandlerOptions) http.Handler {
	r := &handler{
		mux:        chi.NewRouter(),
		signer:     opt.Signer,
		auth:       opt.Auth,
		users:      opt.Users,
		orders:     opt.Orders,
		operations: opt.Operations,
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
			r.Use(h.authorization)

			r.Post("/orders", h.orderProcess)
			r.Get("/orders", h.getOrders)

			r.Get("/balance", h.getBalance)

			r.Post("/balance/withdraw", h.operationPerform)
			r.Get("/withdrawals", h.getOperations)
		})
	})
}
