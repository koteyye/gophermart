package handlers

import (
	"context"
	"errors"
	"net/http"

	"log/slog"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/models"
)

// Auth проверяет наличие токена авторизации в запросе и прокидывает в контекст
// уникальный идентификатор пользователя по ключу models.CtxUserID; если токен
// авторизации не валиден, то возвращает http.StatusUnauthorized.
func (h *handler) auth(next http.Handler) http.Handler {
	auth := func(w http.ResponseWriter, r *http.Request) {
		token, err := parseToken(r.Header.Get("Authorization"))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			slog.Error(err.Error())
			return
		}

		ctx := r.Context()

		id, err := h.service.Auth.Verify(ctx, token)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			slog.Error(err.Error())
			return
		}

		*r = *r.WithContext(context.WithValue(ctx, keyUserID, id))

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(auth)
}
