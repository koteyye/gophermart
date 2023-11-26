package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"log/slog"

	"github.com/google/uuid"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/service"
)

// keyUserID определяет ключ для передачи UserID через контекст.
var keyUserID struct{}

func isAuthorized(ctx context.Context, w http.ResponseWriter) uuid.UUID {
	userID, ok := ctx.Value(keyUserID).(uuid.UUID)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return uuid.Nil
	}
	return userID
}

// auth проверяет наличие токена авторизации в запросе и прокидывает в контекст
// уникальный идентификатор пользователя по ключу keyUserID; если токен
// авторизации не валиден, то возвращает http.StatusUnauthorized.
func (h *handler) auth(next http.Handler) http.Handler {
	auth := func(w http.ResponseWriter, r *http.Request) {
		token, err := parseToken(r.Header.Get("Authorization"))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			h.logger.Error(err.Error())
			return
		}

		ctx := r.Context()

		userID, err := h.service.Verify(ctx, token)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			slog.Error(err.Error())
			return
		}

		*r = *r.WithContext(context.WithValue(ctx, keyUserID, userID))

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(auth)
}

func parseToken(s string) (string, error) {
	split := strings.SplitN(s, " ", 2)
	if len(split) != 2 || split[0] != "Bearer" {
		return "", errors.New("unsupported token")
	}
	return split[1], nil
}
