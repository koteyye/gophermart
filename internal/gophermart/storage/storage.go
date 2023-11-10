package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/sergeizaitcev/gophermart/deployments/gophermart/migrations"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/config"
)

//go:generate mockgen -source=storage.go -destination=mocks/mock.go

// Auth методы для CRUD с users, используется в signIn и signUp
type Auth interface {
	CreateUser(ctx context.Context, login string, password string) (uuid.UUID, error)
	GetUser(ctx context.Context, login string, passwrod string) (uuid.UUID, error)
}

// GophermartDB CRUD операции с БД
type GophermartDB interface {
	Orders
	Balance
}

// Orders - CRUD с заказами
type Orders interface {
	CreateOrder(ctx context.Context, order string) (uuid.UUID, error)
	UpdateOrder(ctx context.Context, order *UpdateOrder) error
	UpdateOrderStatus(ctx context.Context, order string, orderStatus Status) error
	DeleteOrderByNumber(ctx context.Context, order string) error
	GetOrderByNumber(ctx context.Context, order string) (*OrderItem, error)
	GetOrdersByUser(ctx context.Context) ([]*OrderItem, error)
}

// Balance - CRUD с балансом
type Balance interface {
	// CRUD BalanceOperation
	CreateBalanceOperation(ctx context.Context, operation *BalanceOperationItem) (uuid.UUID, error)
	UpdateBalanceOperation(ctx context.Context, order string, done bool) error
	DeleteBalanceOperation(ctx context.Context, order string) error
	GetBalanceOperationByOrderID(
		ctx context.Context,
		order string,
	) (*BalanceOperationItem, error)
	GetBalanceOperation(
		ctx context.Context,
	) ([]*BalanceOperationItem, error)
	// CRUD Balance
	GetBalanceByUserID(ctx context.Context) (*BalanceItem, error)
	IncrementBalance(ctx context.Context, incrementSum int64) error
	DecrementBalance(ctx context.Context, decrementSum int64) error
}

type Storage struct {
	pool *pgxpool.Pool
}

// FIXME: в это функции создаются два разных подключения, чтобы накатить
// миграции и создать одно соединение; необходимо использовать пакет
// `database/sql` для создания единого способа подключения к БД.
func NewStorage(ctx context.Context, c *config.Config) (*Storage, error) {
	err := migrationUp(ctx, c.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("migration: %w", err)
	}

	pool, err := newPool(ctx, c.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("connect to the database: %w", err)
	}

	s := &Storage{
		pool: pool,
	}

	return s, nil
}

func (s *Storage) Close() error {
	s.pool.Close()
	return nil
}

func (s *Storage) Auth() Auth {
	return NewAuthPostgres(s.pool)
}

func newPool(ctx context.Context, dsn string) (pool *pgxpool.Pool, err error) {

	pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create a new connection: %w", err)
	}
	defer func() {
		if err != nil {
			pool.Close()
		}
	}()

	pingCtx, pingCancel := context.WithTimeout(ctx, 3*time.Second)
	defer pingCancel()

	err = pool.Ping(pingCtx)
	if err != nil {
		return nil, fmt.Errorf("database ping: %w", err)
	}

	return pool, nil
}

func migrationUp(ctx context.Context, dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("create a new connection: %w", err)
	}
	defer db.Close()

	err = migrations.Up(ctx, db)
	if err != nil {
		return fmt.Errorf("migration up: %w", err)
	}

	return nil
}
