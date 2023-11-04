package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)


type GophermartDBPostgres struct {
	db *pgx.Conn
}

func NewGophermartPostgres(db *pgx.Conn) *GophermartDBPostgres {
	return &GophermartDBPostgres{db: db}
}

	func (g *GophermartDBPostgres) CreateOrder(ctx context.Context, orderNumber int64) (uuid.UUID, error) {
		return uuid.Nil, errors.New("not implemented")
	}

	func (g *GophermartDBPostgres) UpdateOrder(ctx context.Context, order *UpdateOrder) error {
		return errors.New("not implemented")
	}

	func (g *GophermartDBPostgres) UpdateOrderStatus(ctx context.Context, orderNumber int64, orderStatus Status) error {
		return errors.New("not implemented")
	}

	func (g *GophermartDBPostgres) DeleteOrderByNumber(ctx context.Context, orderNumber int64) error {
		return errors.New("not implemented")
	}

	func (g *GophermartDBPostgres) GetOrderByNumber(ctx context.Context, orderNumber int64) (*OrderItem, error) {
		return nil, errors.New("not implemented")
	}

	func (g *GophermartDBPostgres) CreateBalanceOperation(ctx context.Context, operation *BalanceOperationItem) (uuid.UUID, error) {
		return uuid.Nil, errors.New("not implemented")
	}

	func (g *GophermartDBPostgres) UpdateBalanceOperation(ctx context.Context, operation *BalanceOperationItem) error {
		return errors.New("not implemented")
	}

	func (g *GophermartDBPostgres) DeleteBalanceOperation(ctx context.Context, operationID uuid.UUID) error {
		return errors.New("not implemented")
	}

	func (g *GophermartDBPostgres) GetBalanceOperationByOrderID(ctx context.Context, orderID uuid.UUID) (*BalanceOperationItem, error) {
		return nil, errors.New("not implemented")
	}

	func (g *GophermartDBPostgres) GetBalanceOperationByBalanceID(ctx context.Context, balanceID uuid.UUID) ([]*BalanceOperationItem, error) {
		return nil, errors.New("not implemented")
	}

	func (g *GophermartDBPostgres) GetBalanceByUserID(ctx context.Context, userID uuid.UUID) (*BalanceItem, error) {
		return nil, errors.New("not implemented")
	}

	func (g *GophermartDBPostgres) UpdateBalance(ctx context.Context, userID uuid.UUID, currentSum int64) (error) {
		return errors.New("not implemented")
	}

	func (g *GophermartDBPostgres) IncrementBalance(ctx context.Context, userID uuid.UUID, currentSum int64) (error) {
		return errors.New("not implemented")
	}


	func (g *GophermartDBPostgres) DecrementBalance(ctx context.Context, userID uuid.UUID, currentSum int64) (error) {
		return errors.New("not implemented")
	}