package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
)

var _ domain.UserService = (*Users)(nil)

// Users определяет сервис для работы с пользователем.
type Users struct {
	db *sql.DB
}

// NewUsers возвращает новый экземпляр User.
func NewUsers(db *sql.DB) *Users {
	return &Users{db: db}
}

// GetBalance реализует интерфейс domain.UserService.
func (u *Users) GetBalance(ctx context.Context, id domain.UserID) (domain.UserBalance, error) {
	user, err := getUserByID(ctx, u.db, id)
	if err != nil {
		return domain.UserBalance{}, fmt.Errorf("user search: %w", err)
	}
	return user.Balance, nil
}
