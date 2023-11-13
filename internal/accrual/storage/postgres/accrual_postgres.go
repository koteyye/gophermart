package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sergeizaitcev/gophermart/internal/accrual/storage/models"
)

type AccrualPostgres struct {
	db *pgxpool.Pool
}

func NewAccrualPostgres(db *pgxpool.Pool) *AccrualPostgres {
	return &AccrualPostgres{db: db}
}

// CreateOrder создание записи о заказе в таблице orders и связанную таблицу goods
func (a *AccrualPostgres) CreateOrderWithGoods(ctx context.Context, order string, goods []*models.Goods) (uuid.UUID, error) {
	var orderID uuid.UUID

	// Транзакция для создания записи в  Order и Goods
	tx, err := a.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return uuid.Nil, fmt.Errorf("transaction err: %w", models.MapStorageErr(err))
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, "insert into orders (order_number) values ($1)", order).Scan(&orderID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create order err: %w", models.MapStorageErr(err))
	}

	for _, good := range goods {
		_, err = tx.Exec(ctx, "insert into goods (order_id, match_id, price) values ($1, $2, $3)", orderID, good.MatchID, good.Price)
		if err != nil {
			return uuid.Nil, fmt.Errorf("create goods err: %w", models.MapStorageErr(err))
		}
	}
	tx.Commit(ctx)

	return orderID, nil
}

// UpdateOrder обновляет статус и сумму вознаграждения по заказу
func (a *AccrualPostgres) UpdateOrder(ctx context.Context, order *models.Order) error {
	_, err := a.db.Exec(ctx, "update orders set status = $1, accrual = $2, updated_at = now() where id = $3", order.Status, order.Accrual, order.OrderID)
	if err != nil {
		return fmt.Errorf("update orders err: %w", models.MapStorageErr(err))
	}
	return nil
}

// UpdateGoodAccrual обновляет сумму вознаграждения за конкретный товар
func (a *AccrualPostgres) UpdateGoodAccrual(ctx context.Context, goodID uuid.UUID, accrual int) error {
	_, err := a.db.Exec(ctx, "update goods set accrual = $1, updated_at = now() where id = $2", accrual, goodID)
	if err != nil {
		return fmt.Errorf("update goods err: %w", models.MapStorageErr(err))
	}
	return nil
}

func (a *AccrualPostgres) CreateMatch(ctx context.Context, match *models.Matches) (uuid.UUID, error) {
	var matchID uuid.UUID
	err := a.db.QueryRow(ctx, "insert into matches (match_name, reward, reward_type) values ($1 $2, $3) returning id", match.MatchName, match.Reward, match.RewardType).Scan(&matchID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create match err: %w", models.MapStorageErr(err))
	}
	return matchID, nil
}