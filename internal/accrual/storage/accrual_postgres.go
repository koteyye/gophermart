package storage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccrualPostgres struct {
	db *pgxpool.Pool
}

func NewAccrualPostgres(db *pgxpool.Pool) *AccrualPostgres {
	return &AccrualPostgres{db: db}
}

// CreateOrder создание записи о заказе в таблице orders и связанную таблицу goods
func (a *AccrualPostgres) CreateOrderWithGoods(ctx context.Context, order string, goods []Goods) (uuid.UUID, error) {
	var orderID uuid.UUID

	// Транзакция для создания записи в  Order и Goods
	tx, err := a.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return uuid.Nil, fmt.Errorf("transaction err: %w", mapStorageErr(err))
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, "insert into orders (order_number) values ($1)", order).Scan(&orderID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create order err: %w", mapStorageErr(err))
	}

	for _, good := range goods {
		_, err = tx.Exec(ctx, "insert into goods (order_id, match_id, price) values ($1, $2, $3)", orderID, good.MatchID, good.Price)
		if err != nil {
			return uuid.Nil, fmt.Errorf("create goods err: %w", mapStorageErr(err))
		}
	}
	tx.Commit(ctx)

	return orderID, nil
}

// UpdateOrder обновляет статус и сумму вознаграждения по заказу
func (a *AccrualPostgres) UpdateOrder(ctx context.Context, order Order) error {
	_, err := a.db.Exec(ctx, "update orders set status = $1, accrual = $2, updated_at = now() where id = $3", order.Status, order.Accrual, order.OrderID)
	if err != nil {
		return fmt.Errorf("update orders err: %w", mapStorageErr(err))
	}
	return nil
}

