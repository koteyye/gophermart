package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/storage"
	"github.com/sergeizaitcev/gophermart/pkg/sign"
)

// Auth определяет сервис регистрации и аутентификации.
type Auth struct {
	storage storage.Auth
	signer  sign.Signer
}

// NewAuth возвращает новый экземпляр Auth.
func NewAuth(auth storage.Auth, signer sign.Signer) *Auth {
	return &Auth{
		storage: auth,
		signer:  signer,
	}
}

// SignUp выполняет регистрацию нового пользователя и возвращает токен аутентификации.
func (a *Auth) SignUp(ctx context.Context, login, pass string) (token string, err error) {
	id, err := a.register(ctx, login, pass)
	if err != nil {
		return "", fmt.Errorf("registering a new user: %w", err)
	}

	token, err = a.signer.Sign(id)
	if err != nil {
		return "", fmt.Errorf("creating a token: %w", err)
	}

	return token, nil
}

// SignIn выполняет вход пользователя и возвращает токен аутентификации.
func (a *Auth) SignIn(ctx context.Context, login, pass string) (token string, err error) {
	id, err := a.login(ctx, login, pass)
	if err != nil {
		return "", fmt.Errorf("user authorization: %w", err)
	}

	token, err = a.signer.Sign(id)
	if err != nil {
		return "", fmt.Errorf("creating a token: %w", err)
	}

	return token, nil
}

// Verify проверяет токен аутентификации и возвращает уникальный ID пользователя.
func (a *Auth) Verify(ctx context.Context, token string) (id string, err error) {
	id, err = a.signer.Parse(token)
	if err != nil {
		return "", fmt.Errorf("token verification: %w", err)
	}

	return id, nil
}

// register регистрирует нового пользователя и возвращает уникальный ID.
func (a *Auth) register(ctx context.Context, login, pass string) (id string, err error) {
	err = validate(login, pass)
	if err != nil {
		return "", fmt.Errorf("validation: %w", err)
	}

	uid, err := a.storage.CreateUser(ctx, login, pass)
	if err != nil {
		return "", fmt.Errorf("creating a new user: %w", err)
	}

	return uid.String(), nil
}

// login проводит аутентификацию пользователя и возвращает уникальный ID.
func (a *Auth) login(ctx context.Context, login, pass string) (id string, err error) {
	err = validate(login, pass)
	if err != nil {
		return "", fmt.Errorf("validation: %w", err)
	}

	uid, err := a.storage.GetUser(ctx, login, pass)
	if err != nil {
		return "", fmt.Errorf("getting a user: %w", err)
	}

	return uid.String(), nil
}

func validate(login, pass string) error {
	const maxPassLen = 72
	if login == "" {
		return errors.New("login must be not empty")
	}
	if pass == "" {
		return errors.New("pass must be not empty")
	}
	if len(pass) > maxPassLen {
		return fmt.Errorf("length of pass must be is less than or equal to %d", maxPassLen)
	}
	return nil
}
