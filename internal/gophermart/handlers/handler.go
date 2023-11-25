package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"log/slog"

	"github.com/go-chi/chi/v5"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/handlers/internal/user"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/models"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/service"
)

// handler определяет HTTP-обработчик для gophermart;
// реализует интерфейс http.Handler.
type handler struct {
	mux     *chi.Mux
	service *service.Service
}

// NewHandler возвращает новый экземпляр handler.
func NewHandler(s *service.Service) http.Handler {
	r := &handler{
		mux:     chi.NewRouter(),
		service: s,
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
			r.Get("/hello", h.hello)
		})
	})
}

func (h *handler) register(w http.ResponseWriter, r *http.Request) {
	u, err := decodeUser(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	ctx := r.Context()

	value, err := h.service.Auth.SignUp(ctx, u.Login, u.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicate) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		slog.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "applicaton/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(newToken(value))
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	u, err := decodeUser(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	ctx := r.Context()

	value, err := h.service.Auth.SignIn(ctx, u.Login, u.Password)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		slog.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "applicaton/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(newToken(value))
}

func (h *handler) hello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, ok := user.FromContext(ctx)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_ = id

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, %s!\n", b)
}
