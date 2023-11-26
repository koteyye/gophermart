package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// User определяет связку логин-пароль пользователя.
type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Validate возвращает ошибку, если связка логин-пароль не действительна.
func (u User) Validate() error {
	const maxPassLen = 72
	if u.Login == "" {
		return errors.New("login must be not empty")
	}
	if u.Password == "" {
		return errors.New("pass must be not empty")
	}
	if len(u.Password) > maxPassLen {
		return fmt.Errorf("length of pass must be is less than or equal to %d", maxPassLen)
	}
	return nil
}

// SignUp выполняет регистрацию нового пользователя и возвращает токен аутентификации.
func (s *Service) SignUp(ctx context.Context, u User) (token string, err error) {
	payload, err := s.register(ctx, u)
	if err != nil {
		return "", fmt.Errorf("registering a new user: %w", err)
	}

	token, err = s.signer.Sign(payload)
	if err != nil {
		return "", fmt.Errorf("creating a token: %w", err)
	}

	return token, nil
}

// SignIn выполняет вход пользователя и возвращает токен аутентификации.
func (s *Service) SignIn(ctx context.Context, u User) (token string, err error) {
	payload, err := s.login(ctx, u)
	if err != nil {
		return "", fmt.Errorf("user authorization: %w", err)
	}

	token, err = s.signer.Sign(payload)
	if err != nil {
		return "", fmt.Errorf("creating a token: %w", err)
	}

	return token, nil
}

// Verify проверяет токен аутентификации и возвращает уникальный ID пользователя.
func (s *Service) Verify(ctx context.Context, token string) (uuid.UUID, error) {
	payload, err := s.signer.Parse(token)
	if err != nil {
		return uuid.Nil, fmt.Errorf("token verification: %w", err)
	}

	userID, err := uuid.Parse(payload)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parsing UUID: %w", err)
	}

	exists, err := s.storage.UserExists(ctx, userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("user search: %w", err)
	}
	if !exists {
		return uuid.Nil, errors.New("user not found")
	}

	return userID, nil
}

// register регистрирует нового пользователя и возвращает уникальный ID.
func (s *Service) register(ctx context.Context, u User) (id string, err error) {
	err = u.Validate()
	if err != nil {
		return "", fmt.Errorf("validation: %w", err)
	}

	userID, err := s.storage.CreateUser(ctx, u)
	if err != nil {
		return "", fmt.Errorf("creating a new user: %w", err)
	}

	return userID.String(), nil
}

// login проводит аутентификацию пользователя и возвращает уникальный ID.
func (s *Service) login(ctx context.Context, u User) (id string, err error) {
	err = u.Validate()
	if err != nil {
		return "", fmt.Errorf("validation: %w", err)
	}

	userID, err := s.storage.UserID(ctx, u)
	if err != nil {
		return "", fmt.Errorf("user search: %w", err)
	}

	return userID.String(), nil
}
