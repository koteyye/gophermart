package server

import (
	"context"
	"errors"
	"net/http"

	"log/slog"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/models"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/service"
)

// Auth проверяет наличие токена авторизации в запросе и прокидывает в контекст
// уникальный идентификатор пользователя по ключу models.CtxUserID; если токен
// авторизации не валиден, то возвращает http.StatusUnauthorized.
func Auth(s *service.Auth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := parseToken(r.Header.Get("Authorization"))
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				slog.Error(err.Error())
				return
			}

			ctx := r.Context()

			id, err := s.Verify(ctx, token)
			if err != nil {
				if errors.Is(err, models.ErrNotFound) {
					w.WriteHeader(http.StatusUnauthorized)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
				slog.Error(err.Error())
				return
			}

			*r = *r.WithContext(context.WithValue(ctx, models.KeyUserID, id))

			next.ServeHTTP(w, r)
		})
	}
}
