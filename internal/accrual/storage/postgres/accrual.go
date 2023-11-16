package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sergeizaitcev/gophermart/internal/accrual/storage"
)

var _ storage.Accrual = (*Accrual)(nil)

type Accrual struct {
	db *pgxpool.Pool
}

func NewAccrual(db *pgxpool.Pool) *Accrual {
	return &Accrual{db: db}
}

// CreateOrderWithGoods создание записи о заказе в таблице orders и связанную таблицу goods
func (a *Accrual) CreateOrderWithGoods(ctx context.Context, order string, goods []*storage.Goods) (uuid.UUID, error) {
	var orderID uuid.UUID

	// Транзакция для создания записи в  Order и Goods
	tx, err := a.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return uuid.Nil, fmt.Errorf("transaction err: %w", mapStorageErr(err))
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, "insert into orders (order_number) values ($1) returning id", order).Scan(&orderID)
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

// CreateInvalidOrder создает запись о заказе без Goods со статусом Invalid. Для кейсов, когда не найдено соответствующего товара
func (a *Accrual) CreateInvalidOrder(ctx context.Context, order string) error {
	_, err := a.db.Exec(ctx, "insert into orders (order_number, status) values ($1, $2)", order, storage.OrderStatus(1))
	if err != nil {
		return fmt.Errorf("insert invalid orders err: %w", mapStorageErr(err))
	}
	return nil
}

// UpdateOrder обновляет статус и сумму вознаграждения по заказу
func (a *Accrual) UpdateOrder(ctx context.Context, order *storage.Order) error {
	_, err := a.db.Exec(ctx, "update orders set status = $1, accrual = $2, updated_at = now() where id = $3", order.Status, order.Accrual, order.OrderID)
	if err != nil {
		return fmt.Errorf("update orders err: %w", mapStorageErr(err))
	}
	return nil
}

// UpdateGoodAccrual обновляет сумму вознаграждения за конкретный товар
func (a *Accrual) UpdateGoodAccrual(ctx context.Context, matchID uuid.UUID, accrual float64) error {
	_, err := a.db.Exec(ctx, "update goods set accrual = $1, updated_at = now() where id = $2", accrual, matchID)
	if err != nil {
		return fmt.Errorf("update goods err: %w", mapStorageErr(err))
	}
	return nil
}

// CreateMatch создает новую механику вознаграждения для товара
func (a *Accrual) CreateMatch(ctx context.Context, match *storage.Match) (uuid.UUID, error) {
	var matchID uuid.UUID
	err := a.db.QueryRow(ctx, "insert into matches (match_name, reward, reward_type) values ($1, $2, $3) returning id", match.MatchName, match.Reward, match.Type).Scan(&matchID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create match err: %w", mapStorageErr(err))
	}
	return matchID, nil
}

// GetMatchByName возвращает механику вознаграждения для товара
func (a *Accrual) GetMatchByName(ctx context.Context, matchName string) (*storage.MatchOut, error) {
	var matches storage.MatchOut

	err := a.db.QueryRow(ctx, "select id, match_name, reward, reward_type from matches where match_name in ($1)", matchName).
		Scan(&matches.MatchID, &matches.MatchName, &matches.Reward, &matches.Type)
	if err != nil {
		return nil, fmt.Errorf("select match err: %w", mapStorageErr(err))
	}
	return &matches, nil
}

func (a *Accrual) GetOrderWithGoodsByNumber(ctx context.Context, orderNumber string) (*storage.OrderOut, error) {
	var order storage.OrderOut

	err := a.db.QueryRow(ctx, "select order_number, status, accrual from orders where order_number = $1", orderNumber).Scan(&order.OrderNumber, &order.Status, &order.Accrual)
	if err != nil {
		return &storage.OrderOut{}, fmt.Errorf("select order err: %w", mapStorageErr(err))
	}

	return &order, nil
}

func (a *Accrual) GetGoodsByOrderID(ctx context.Context, orderID uuid.UUID) ([]*storage.Goods, error) {
	rows, err := a.db.Query(ctx, "select id, match_id, price, accrual from goods where order_id = $1", orderID)
	if err != nil {
		return nil, fmt.Errorf("select good by orderID err: %w", mapStorageErr(err))
	}

	goods, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[storage.Goods])
	if err != nil {
		return nil, fmt.Errorf("scan goods err: %w", mapStorageErr(err))
	}

	return goods, nil
}
