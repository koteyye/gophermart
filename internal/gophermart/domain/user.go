package domain

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/sergeizaitcev/gophermart/pkg/monetary"
	"github.com/sergeizaitcev/gophermart/pkg/passwords"
)

// User определяет пользователя.
type User struct {
	ID             UserID      // Уникальный идентификатор пользователя.
	Login          string      // Логин пользователя.
	HashedPassword string      // Хеш-сумма пароля пользователя.
	Balance        UserBalance // Баланс пользователя.
}

// ComparePassword возвращает true, если пароль содержит хеш-сумму пароля.
func (u User) ComparePassword(password string) bool {
	return passwords.Compare(u.HashedPassword, password)
}

// UserBalance определяет баланс пользователя.
type UserBalance struct {
	Current   monetary.Unit `json:"current"`   // Начисленные баллы.
	Withdrawn monetary.Unit `json:"withdrawn"` // Списанные баллы.
}

// NewUser конвертирует данные аутентификации в пользователя и возвращает его.
func NewUser(auth Authentication) (User, error) {
	var user User
	hashedPassword, err := passwords.Hash(auth.Password)
	if err != nil {
		return user, fmt.Errorf("hashing password: %w", err)
	}
	user.Login = auth.Login
	user.HashedPassword = hashedPassword
	return user, nil
}

// UserID определяет уникальный идентификатор пользователя.
type UserID = uuid.UUID

var EmptyUserID = uuid.Nil

// NewUserID конвертирует строку в уникальный идентификатор пользователя
// и возвращает его.
func NewUserID(s string) (UserID, error) {
	uid, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parsing user ID: %w", err)
	}
	return uid, nil
}

// Authentication определяет данные для аутентификации пользователя.
type Authentication struct {
	Login    string `json:"login"`    // Логин пользователя.
	Password string `json:"password"` // Пароль пользователя.
}

func (a Authentication) Validate() error {
	const maxPassLen = 72
	if a.Login == "" {
		return errors.New("login must be not empty")
	}
	if a.Password == "" {
		return errors.New("password must be not empty")
	}
	if len(a.Password) > maxPassLen {
		return fmt.Errorf("length of pass must be is less than or equal to %d", maxPassLen)
	}
	return nil
}
