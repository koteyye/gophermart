package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/storage"
	"github.com/sergeizaitcev/gophermart/pkg/sign"
)

// userService определяет сервис регистрации и аутентификации.
type userService struct {
	storage storage.Auth
	signer  sign.Signer
}

// newUserService возвращает новый экземпляр userService.
func newUserService(auth storage.Auth, signer sign.Signer) *userService {
	return &userService{
		storage: auth,
		signer:  signer,
	}
}

// SignUp выполняет регистрацию нового пользователя и возвращает токен аутентификации.
func (s *userService) SignUp(ctx context.Context, login, pass string) (token string, err error) {
	id, err := s.register(ctx, login, pass)
	if err != nil {
		return "", fmt.Errorf("registering a new user: %w", err)
	}

	token, err = s.signer.Sign(id)
	if err != nil {
		return "", fmt.Errorf("creating a token: %w", err)
	}

	return token, nil
}

// SignIn выполняет вход пользователя и возвращает токен аутентификации.
func (s *userService) SignIn(ctx context.Context, login, pass string) (token string, err error) {
	id, err := s.login(ctx, login, pass)
	if err != nil {
		return "", fmt.Errorf("user authorization: %w", err)
	}

	token, err = s.signer.Sign(id)
	if err != nil {
		return "", fmt.Errorf("creating a token: %w", err)
	}

	return token, nil
}

// Verify проверяет токен аутентификации и возвращает уникальный ID пользователя.
func (s *userService) Verify(ctx context.Context, token string) (id string, err error) {
	id, err = s.signer.Parse(token)
	if err != nil {
		return "", fmt.Errorf("token verification: %w", err)
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		return "", fmt.Errorf("parsing UUID: %w", err)
	}

	_, err = s.storage.GetLogin(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("checking a user: %w", err)
	}

	return id, nil
}

func (s *userService) Balance(ctx context.Context, id string)

// register регистрирует нового пользователя и возвращает уникальный ID.
func (s *userService) register(ctx context.Context, login, pass string) (id string, err error) {
	err = validate(login, pass)
	if err != nil {
		return "", fmt.Errorf("validation: %w", err)
	}

	uid, err := s.storage.CreateUser(ctx, login, pass)
	if err != nil {
		return "", fmt.Errorf("creating a new user: %w", err)
	}

	return uid.String(), nil
}

// login проводит аутентификацию пользователя и возвращает уникальный ID.
func (s *userService) login(ctx context.Context, login, pass string) (id string, err error) {
	err = validate(login, pass)
	if err != nil {
		return "", fmt.Errorf("validation: %w", err)
	}

	uid, err := s.storage.GetUser(ctx, login, pass)
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
