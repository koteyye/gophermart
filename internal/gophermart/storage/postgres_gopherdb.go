package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/models"
)

type GophermartDBPostgres struct {
	db *pgx.Conn
}

func NewGophermartPostgres(db *pgx.Conn) *GophermartDBPostgres {
	return &GophermartDBPostgres{db: db}
}

// CreateOrder - создает в таблице orders запись с номером заказа на статусе new
func (g *GophermartDBPostgres) CreateOrder(ctx context.Context, orderNumber int64) (uuid.UUID, error) {
	var orderID uuid.UUID
	userID := ctx.Value(models.CtxUserID)

	// creatorUserID ID пользователя, ранее загрузившего заказ
	var creatorUserID uuid.UUID

	// Запускаем транзацкию чтобы сначала проверить наличие в БД добавляемого номера заказа и кто его добавил
	// затем добавляем запись если ее нет
	tx, err := g.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
		tx.Commit(ctx)
	}()
	if err != nil {
		return uuid.Nil, fmt.Errorf("trasnaction err: %w", mapStorageErr(err))
	}

	err = tx.QueryRow(ctx, "select id, user_created from orders where order_number = $1", orderNumber).Scan(&orderID, &creatorUserID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("select order err: %w", mapStorageErr(err))
		}
	}

	// Возвращаем ошибку: текущий пользователь или другой пользователь ранее загрузил номер заказа
	if creatorUserID == userID {
		return orderID, ErrDuplicate
	} else if creatorUserID != uuid.Nil && creatorUserID != userID{
		return orderID, ErrDuplicateOtherUser
	}
	
	// Если нет ошибок, то создаем новую запись
	
	err = tx.QueryRow(ctx, "insert into orders (order_number, user_created) values ($1, $2) returning id", orderNumber, userID).Scan(&orderID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("insert orders err: %w", mapStorageErr(err))
	}
	return orderID, nil
}

// UpdatedOrder обновляет статус и баллы по номеру заказа
func (g *GophermartDBPostgres) UpdateOrder(ctx context.Context, order *UpdateOrder) error {
	_, err := g.db.Exec(ctx, "update orders set status = $1, accrual = $2 where order_number = $3", order.Status, order.Accrual, order.Number)
	if err != nil {
		return fmt.Errorf("update order err: %w", mapStorageErr(err))
	}

	return nil
}

// UpdateOrderStatus - обновляет статус заказа
func (g *GophermartDBPostgres) UpdateOrderStatus(ctx context.Context, orderNumber int64, orderStatus Status) error {
	_, err := g.db.Exec(ctx, "update orders set status = $1 where order_number = $2", orderStatus, orderNumber)
	if err != nil {
		return fmt.Errorf("update order status: %w", mapStorageErr(err))
	}
	return nil
}

// GetOrderByNumber - возвращает информацию о заказе по номеру
func (g *GophermartDBPostgres) GetOrderByNumber(ctx context.Context, orderNumber int64) (*OrderItem, error) {
	var order OrderItem
	row := g.db.QueryRow(ctx, "select id, order_number, status, accrual, user_created from orders where order_number = $1 and deleted_at is null", orderNumber)
	err := row.Scan(&order.ID, &order.OrderNumber, &order.Status, &order.Accrual, &order.UserID)
	if err != nil {
		return nil, fmt.Errorf("select order err: %w", mapStorageErr(err))
	}
	return &order, nil
}

// DeleteOrderByNumber - удаляет заказ по номеру
func (g *GophermartDBPostgres) DeleteOrderByNumber(ctx context.Context, orderNumber int64) error {
	_, err := g.db.Exec(ctx, "update orders set deleted_at = now() where order_number = $1", orderNumber)
	if err != nil {
		return fmt.Errorf("delete order err: %w", mapStorageErr(err))
	}
	return nil
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

func (g *GophermartDBPostgres) UpdateBalance(ctx context.Context, userID uuid.UUID, currentSum int64) error {
	return errors.New("not implemented")
}

func (g *GophermartDBPostgres) IncrementBalance(ctx context.Context, userID uuid.UUID, currentSum int64) error {
	return errors.New("not implemented")
}

func (g *GophermartDBPostgres) DecrementBalance(ctx context.Context, userID uuid.UUID, currentSum int64) error {
	return errors.New("not implemented")
}
