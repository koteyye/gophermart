package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)


type gophermartDBPostgres struct {
	db *pgx.Conn
}

func NewGophermartPostgres(db *pgx.Conn) *gophermartDBPostgres {
	return &gophermartDBPostgres{db: db}
}

	func (g *gophermartDBPostgres) CreateOrder(ctx context.Context, orderNumber int64) (uuid.UUID, error) {
		return uuid.Nil, errors.New("not implemented")
	}

	func (g *gophermartDBPostgres) UpdateOrder(ctx context.Context, order *UpdateOrder) error {
		return errors.New("not implemented")
	}

	func (g *gophermartDBPostgres) UpdateOrderStatus(ctx context.Context, orderNumber int64, orderStatus string) error {
		return errors.New("not implemented")
	}

	func (g *gophermartDBPostgres) DeleteOrderByNumber(ctx context.Context, orderNumber int64) error {
		return errors.New("not implemented")
	}

	func (g *gophermartDBPostgres) GetOrderByNumber(ctx context.Context, orderNumber int64) (*OrderItem, error) {
		return nil, errors.New("not implemented")
	}

	func (g *gophermartDBPostgres) CreateBalanceOperation(ctx context.Context, operation *BalanceOperationItem) (uuid.UUID, error) {
		return uuid.Nil, errors.New("not implemented")
	}

	func (g *gophermartDBPostgres) UpdateBalanceOperation(ctx context.Context, operation *BalanceOperationItem) error {
		return errors.New("not implemented")
	}

	func (g *gophermartDBPostgres) DeleteBalanceOperation(ctx context.Context, operationID uuid.UUID) error {
		return errors.New("not implemented")
	}

	func (g *gophermartDBPostgres) GetBalanceOperationByOrderID(ctx context.Context, orderID uuid.UUID) (*BalanceOperationItem, error) {
		return nil, errors.New("not implemented")
	}

	func (g *gophermartDBPostgres) GetBalanceOperationByBalanceID(ctx context.Context, balanceID uuid.UUID) ([]*BalanceOperationItem, error) {
		return nil, errors.New("not implemented")
	}

	func (g *gophermartDBPostgres) GetBalanceByUserID(ctx context.Context, userID uuid.UUID) (*BalanceItem, error) {
		return nil, errors.New("not implemented")
	}

	func (g *gophermartDBPostgres) UpdateBalance(ctx context.Context, userID uuid.UUID, currentSum int64) (error) {
		return errors.New("not implemented")
	}
