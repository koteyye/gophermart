package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
)

var _ domain.AuthService = (*Auth)(nil)

// Auth определяет сервис регистрации и аутентификации пользователя.
type Auth struct {
	db *sql.DB
}

// NewAuth возвращает новый экземпляр Auth.
func NewAuth(db *sql.DB) *Auth {
	return &Auth{db: db}
}

// Identify реализует интерфейс domain.AuthService.
func (a *Auth) Identify(ctx context.Context, id domain.UserID) error {
	u, err := getUserByID(ctx, a.db, id)
	if err != nil {
		return fmt.Errorf("user search: %w", err)
	}

	if u.ID != id {
		return fmt.Errorf("user ID not identified: %q", id)
	}

	return nil
}

// SignIn реализует интерфейс domain.AuthService.
func (a *Auth) SignIn(ctx context.Context, auth domain.Authentication) (domain.UserID, error) {
	user, err := getUser(ctx, a.db, auth)
	if err != nil {
		return domain.EmptyUserID, fmt.Errorf("user search: %w", err)
	}

	if user.Login != auth.Login || !user.ComparePassword(auth.Password) {
		return domain.EmptyUserID, domain.ErrNotFound
	}

	return user.ID, nil
}

// SignUp реализует интерфейс domain.AuthService.
func (a *Auth) SignUp(ctx context.Context, auth domain.Authentication) (domain.UserID, error) {
	user, err := domain.NewUser(auth)
	if err != nil {
		return domain.EmptyUserID, fmt.Errorf("converting to user: %w", err)
	}

	uid, err := createUser(ctx, a.db, user)
	if err != nil {
		return domain.EmptyUserID, fmt.Errorf("creating a new user: %w", err)
	}

	return uid, nil
}

func createUser(
	ctx context.Context,
	db *sql.DB,
	user domain.User,
) (domain.UserID, error) {
	var userID domain.UserID

	query := "INSERT INTO users (login, hashed_password) VALUES ($1, $2) RETURNING id;"

	err := db.QueryRowContext(ctx, query, user.Login, user.HashedPassword).Scan(&userID)
	if err != nil {
		return domain.EmptyUserID, fmt.Errorf("creating a new user: %w", errorHandling(err))
	}

	return userID, nil
}

func getUser(
	ctx context.Context,
	db *sql.DB,
	auth domain.Authentication,
) (domain.User, error) {
	var user domain.User

	query := `SELECT
		id, login, hashed_password, current_balance, withdrawn_balance
	FROM users
	WHERE login = $1;`

	err := db.QueryRowContext(ctx, query, auth.Login).Scan(
		&user.ID,
		&user.Login,
		&user.HashedPassword,
		&user.Balance.Current,
		&user.Balance.Withdrawn,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("user search: %w", errorHandling(err))
	}

	return user, nil
}

func getUserByID(
	ctx context.Context,
	db *sql.DB,
	id domain.UserID,
) (domain.User, error) {
	var user domain.User

	query := `SELECT
		id, login, hashed_password, current_balance, withdrawn_balance
	FROM users
	WHERE id = $1;`

	err := db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Login,
		&user.HashedPassword,
		&user.Balance.Current,
		&user.Balance.Withdrawn,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("user search: %w", errorHandling(err))
	}

	return user, nil
}
