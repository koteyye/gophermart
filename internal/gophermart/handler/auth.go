package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"log/slog"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
)

// Тип токена авторизации.
const tokenType = "Bearer"

// token определяет токен авторизации.
type token struct {
	value string
}

func (t token) String() string {
	return tokenType + " " + t.value
}

// parseToken парсит строку и возвращает токен авторизации.
func parseToken(s string) (token, error) {
	split := strings.SplitN(s, " ", 2)
	if len(split) != 2 || split[0] != tokenType {
		return token{}, fmt.Errorf("unsupported token type: %q", split[0])
	}
	return token{value: split[1]}, nil
}

// toToken упаковывает уникальный идентификатор пользователя в токен
// авторизации и возвращает его.
func (h *handler) toToken(id domain.UserID) (token, error) {
	tokenValue, err := h.signer.Sign(id.String())
	if err != nil {
		return token{}, err
	}
	return token{value: tokenValue}, nil
}

// parseToken парсит токен авторизации и возвращает уникальный идентификатор
// пользователя.
func (h *handler) parseToken(token string) (domain.UserID, error) {
	tk, err := parseToken(token)
	if err != nil {
		return domain.EmptyUserID, err
	}
	payload, err := h.signer.Parse(tk.value)
	if err != nil {
		return domain.EmptyUserID, err
	}
	userID, err := domain.NewUserID(payload)
	if err != nil {
		return domain.EmptyUserID, err
	}
	return userID, nil
}

// keyUserID определяет ключ для передачи domain.UserID через контекст.
var keyUserID struct{}

// userFromContext возвращает уникальный идентификатор пользователя
// из контекста.
func userFromContext(ctx context.Context) domain.UserID {
	userID, ok := ctx.Value(keyUserID).(domain.UserID)
	if !ok {
		return domain.EmptyUserID
	}
	return userID
}

// authorization проверяет наличие токена авторизации в запросе и прокидывает
// в контекст уникальный идентификатор пользователя по ключу keyUserID; если
// токен авторизации не действителен, то возвращает http.StatusUnauthorized.
func (h *handler) authorization(next http.Handler) http.Handler {
	auth := func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		userID, err := h.parseToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			slog.Error(err.Error())
			return
		}

		ctx := r.Context()

		err = h.auth.Identify(ctx, userID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
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

// register выполняет регистрацию пользователя и возвращает в заголовке
// ответа токен авторизации.
func (h *handler) register(w http.ResponseWriter, r *http.Request) {
	var auth domain.Authentication

	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	err = auth.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	ctx := r.Context()

	userID, err := h.auth.SignUp(ctx, auth)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicate) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		slog.Error(err.Error())
		return
	}

	token, err := h.toToken(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}

	w.Header().Set("Authorization", token.String())
}

// login выполняет аутентификацию пользователя и возвращает в заголовке
// ответа токен авторизации.
func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	var auth domain.Authentication

	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	err = auth.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	ctx := r.Context()

	userID, err := h.auth.SignIn(ctx, auth)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		slog.Error(err.Error())
		return
	}

	token, err := h.toToken(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}

	w.Header().Set("Authorization", token.String())
}
