package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	CreateOrder(ctx context.Context, orderNumber int64) (uuid.UUID, error)
	UpdateOrder(ctx context.Context, order *UpdateOrder) error
	UpdateOrderStatus(ctx context.Context, orderNumber int64, orderStatus Status) error
	DeleteOrderByNumber(ctx context.Context, orderNumber int64) error
	GetOrderByNumber(ctx context.Context, orderNumber int64) (*OrderItem, error)
}

// Balance - CRUD с балансом
type Balance interface {
	// CRUD BalanceOperation
	CreateBalanceOperation(ctx context.Context, operation *BalanceOperationItem) (uuid.UUID, error)
	UpdateBalanceOperation(ctx context.Context, operation *BalanceOperationItem) error
	DeleteBalanceOperation(ctx context.Context, operationID uuid.UUID) error
	GetBalanceOperationByOrderID(
		ctx context.Context,
		orderID uuid.UUID,
	) (*BalanceOperationItem, error)
	GetBalanceOperationByBalanceID(
		ctx context.Context,
		balanceID uuid.UUID,
	) ([]*BalanceOperationItem, error)
	// CRUD Balance
	GetBalanceByUserID(ctx context.Context, userID uuid.UUID) (*BalanceItem, error)
	UpdateBalance(ctx context.Context, userID uuid.UUID, currentSum int64) error
	IncrementBalance(ctx context.Context, userID uuid.UUID, currentSum int64) error
	DecrementBalance(ctx context.Context, userID uuid.UUID, currentSum int64) error
}

type Storage struct {
	conn *pgx.Conn
}

// FIXME: в это функции создаются два разных подключения, чтобы накатить
// миграции и создать одно соединение; необходимо использовать пакет
// `database/sql` для создания единого способа подключения к БД.
func NewStorage(ctx context.Context, c *config.Config) (*Storage, error) {
	err := migrationUp(ctx, c.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("migration: %w", err)
	}

	conn, err := newConn(ctx, c.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("connect to the database: %w", err)
	}

	s := &Storage{
		conn: conn,
	}

	return s, nil
}

func (s *Storage) Close() error {
	err := s.conn.Close(context.Background())
	if err != nil {
		return fmt.Errorf("close connection: %w", err)
	}
	return nil
}

func (s *Storage) Auth() Auth {
	return NewAuthPostgres(s.conn)
}

func newConn(ctx context.Context, dsn string) (conn *pgx.Conn, err error) {
	conn, err = pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create a new connection: %w", err)
	}
	defer func() {
		if err != nil {
			_ = conn.Close(context.Background())
		}
	}()

	pingCtx, pingCancel := context.WithTimeout(ctx, 3*time.Second)
	defer pingCancel()

	err = conn.Ping(pingCtx)
	if err != nil {
		return nil, fmt.Errorf("database ping: %w", err)
	}

	return conn, nil
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
