package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/sergeizaitcev/gophermart/internal/accrual/storage"
)

// CreateOrderWithGoods создание записи о заказе в таблице orders и связанную таблицу goods
func (a *Storage) CreateOrderWithGoods(ctx context.Context, order string, goods []*storage.Goods) (uuid.UUID, error) {
	var orderID uuid.UUID

	query1 := "insert into orders (order_number) values ($1) returning id"
	query2 := "insert into goods (order_id, match_id, price) values ($1, $2, $3)"

	// Транзакция для создания записи в  Order и Goods
	err := a.transaction(ctx, func(tx *sql.Tx) error {
		err := tx.QueryRowContext(ctx, query1, order).Scan(&orderID)
		if err != nil {
			return fmt.Errorf("create order err: %w", errorHandle(err))
		}

		for _, good := range goods {
			_, err := tx.ExecContext(ctx, query2, orderID, good.MatchID, good.Price)
			if err != nil {
				return fmt.Errorf("create goods err: %w", errorHandle(err))
			}
		}

		return nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	return orderID, nil
}

// CreateInvalidOrder создает запись о заказе без Goods со статусом Invalid. Для кейсов, когда не найдено соответствующего товара
func (a *Storage) CreateInvalidOrder(ctx context.Context, order string) error {
	_, err := a.db.ExecContext(ctx, "insert into orders (order_number, status) values ($1, $2)", order, storage.OrderStatus(1))
	if err != nil {
		return fmt.Errorf("insert invalid orders err: %w", errorHandle(err))
	}
	return nil
}

// UpdateOrder обновляет статус и сумму вознаграждения по заказу
func (a *Storage) UpdateOrder(ctx context.Context, order *storage.Order) error {
	_, err := a.db.ExecContext(ctx, "update orders set status = $1, accrual = $2, updated_at = now() where id = $3", order.Status, order.Accrual, order.OrderID)
	if err != nil {
		return fmt.Errorf("update orders err: %w", errorHandle(err))
	}
	return nil
}

// UpdateGoodAccrual обновляет сумму вознаграждения за конкретный товар
func (a *Storage) UpdateGoodAccrual(ctx context.Context, orderID uuid.UUID, matchID uuid.UUID, accrual float64) error {
	query1 := "update goods set accrual = $1, updated_at = now() where match_id = $2 and order_id = $3"

	_, err := a.db.ExecContext(ctx, query1, accrual, matchID, orderID)
	if err != nil {
		return fmt.Errorf("update goods err: %w", errorHandle(err))
	}

	return nil
}

// CreateMatch создает новую механику вознаграждения для товара
func (a *Storage) CreateMatch(ctx context.Context, match *storage.Match) (uuid.UUID, error) {
	var matchID uuid.UUID
	err := a.db.QueryRowContext(ctx, "insert into matches (match_name, reward, reward_type) values ($1, $2, $3) returning id", match.MatchName, match.Reward, match.Type).Scan(&matchID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create match err: %w", errorHandle(err))
	}
	return matchID, nil
}

// GetMatchByName возвращает механику вознаграждения для товара
func (a *Storage) GetMatchByName(ctx context.Context, matchName string) (*storage.MatchOut, error) {
	var matches storage.MatchOut

	err := a.db.QueryRowContext(ctx, "select id, match_name, reward, reward_type from matches where match_name in ($1)", matchName).
		Scan(&matches.MatchID, &matches.MatchName, &matches.Reward, &matches.Type)
	if err != nil {
		return nil, fmt.Errorf("select match err: %w", errorHandle(err))
	}
	return &matches, nil
}

// GetOrderByNumber возвращает статус заказа и вознаграждение
func (a *Storage) GetOrderByNumber(ctx context.Context, orderNumber string) (*storage.OrderOut, error) {
	var order storage.OrderOut
	query := "select order_number, status, accrual from orders where order_number = $1 and deleted_at is null"

	err := a.db.QueryRowContext(ctx, query, orderNumber).Scan(&order.OrderNumber, &order.Status, &order.Accrual)
	if err != nil {
		return &storage.OrderOut{}, fmt.Errorf("select order err: %w", errorHandle(err))
	}

	return &order, nil
}

func (a *Storage) GetGoodsByOrderID(ctx context.Context, orderID uuid.UUID) ([]*storage.Goods, error) {
	var goods []*storage.Goods

	rows, err := a.db.QueryContext(ctx, "select id, match_id, price, accrual from goods where order_id = $1", orderID)
	if err != nil {
		return nil, fmt.Errorf("select good by orderID err: %w", errorHandle(err))
	}
	defer rows.Close()

	for rows.Next() {
		var good storage.Goods

		err := rows.Scan(
			&good.GoodID,
			&good.MatchID,
			&good.Price,
			&good.Accrual,
		)
		if err != nil {
			return nil, fmt.Errorf("scan goods err: %w", errorHandle(err))
		}

		goods = append(goods, &good)
	}

	return goods, nil
}
