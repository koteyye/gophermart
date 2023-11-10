package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/models"
)

type GophermartDBPostgres struct {
	db *pgxpool.Pool
}

func NewGophermartPostgres(db *pgxpool.Pool) *GophermartDBPostgres {
	return &GophermartDBPostgres{db: db}
}

// CreateOrder - создает в таблице orders запись с номером заказа на статусе new
func (g *GophermartDBPostgres) CreateOrder(
	ctx context.Context,
	order string,
) (uuid.UUID, error) {
	var orderID uuid.UUID
	userID := ctx.Value(models.KeyUserID)

	// creatorUserID ID пользователя, ранее загрузившего заказ
	var creatorUserID uuid.UUID

	// Запускаем транзацкию чтобы сначала проверить наличие в БД добавляемого номера заказа и кто его добавил
	// затем добавляем запись если ее нет
	tx, err := g.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return uuid.Nil, fmt.Errorf("trasnaction err: %w", models.MapStorageErr(err))
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, "select id, user_created from orders where order_number = $1", order).
		Scan(&orderID, &creatorUserID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("select order err: %w", models.MapStorageErr(err))
		}
	}

	// Возвращаем ошибку: текущий пользователь или другой пользователь ранее загрузил номер заказа
	if creatorUserID == userID {
		return orderID, models.ErrDuplicate
	} else if creatorUserID != uuid.Nil && creatorUserID != userID {
		return orderID, models.ErrDuplicateOtherUser
	}

	// Если нет ошибок, то создаем новую запись

	err = tx.QueryRow(ctx, "insert into orders (order_number, user_created) values ($1, $2) returning id", order, userID).
		Scan(&orderID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("insert orders err: %w", models.MapStorageErr(err))
	}
	tx.Commit(ctx)
	return orderID, nil
}

// UpdatedOrder обновляет статус и баллы по номеру заказа
func (g *GophermartDBPostgres) UpdateOrder(ctx context.Context, order *UpdateOrder) error {
	_, err := g.db.Exec(
		ctx,
		"update orders set status = $1, accrual = $2, updated_at = now() where order_number = $3",
		order.Status,
		order.Accrual,
		order.Order,
	)
	if err != nil {
		return fmt.Errorf("update order err: %w", models.MapStorageErr(err))
	}

	return nil
}

// UpdateOrderStatus - обновляет статус заказа
func (g *GophermartDBPostgres) UpdateOrderStatus(
	ctx context.Context,
	order string,
	orderStatus Status,
) error {
	_, err := g.db.Exec(
		ctx,
		"update orders set status = $1, updated_at = now() where order_number = $2",
		orderStatus,
		order,
	)
	if err != nil {
		return fmt.Errorf("update order status: %w", models.MapStorageErr(err))
	}
	return nil
}

// GetOrderByNumber - возвращает информацию о заказе по номеру
func (g *GophermartDBPostgres) GetOrderByNumber(
	ctx context.Context,
	order string,
) (*OrderItem, error) {
	var orderInfo OrderItem
	row := g.db.QueryRow(
		ctx,
		"select id, order_number, status, accrual, user_created, updated_at from orders where order_number = $1 and deleted_at is null",
		order,
	)
	err := row.Scan(&orderInfo.ID, &orderInfo.Order, &orderInfo.Status, &orderInfo.Accrual, &orderInfo.UserID, &orderInfo.UploadedAt)
	if err != nil {
		return nil, fmt.Errorf("select order err: %w", models.MapStorageErr(err))
	}
	return &orderInfo, nil
}

//GetOrdersByUser - возвращает все заказы по текущему пользователю
func (g *GophermartDBPostgres) GetOrdersByUser(ctx context.Context) ([]*OrderItem, error) {
	userID := ctx.Value(models.KeyUserID)

	rows, err := g.db.Query(ctx, "select id, order_number, status, accrual, user_created, updated_at from orders where user_created = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("select orders by userID err: %w", models.MapStorageErr(err))
	}

	orders, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[OrderItem])
	if err != nil {
		return nil, fmt.Errorf("scan orders err: %w", models.MapStorageErr(err))
	}

	return orders, nil
}

// DeleteOrderByNumber - удаляет заказ по номеру
func (g *GophermartDBPostgres) DeleteOrderByNumber(ctx context.Context, order string) error {
	_, err := g.db.Exec(
		ctx,
		"update orders set deleted_at = now() where order_number = $1",
		order,
	)
	if err != nil {
		return fmt.Errorf("delete order err: %w", models.MapStorageErr(err))
	}
	return nil
}

// CreateBalanceOperation - создает операцию с балансом (запрос на списание)
func (g *GophermartDBPostgres) CreateBalanceOperation(
	ctx context.Context,
	operation int64,
	order string,
) (uuid.UUID, error) {
	var bOperationID uuid.UUID
	var balanceID uuid.UUID
	var orderID uuid.UUID
	userID := ctx.Value(models.KeyUserID)

	// Транзакция для получения OrderID и BalanceID, чтоб использовать их для создания Balance operation
	tx, err := g.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return uuid.Nil, fmt.Errorf("transaction err: %w", models.MapStorageErr(err))
	}
	defer tx.Rollback(ctx)

	if err = tx.QueryRow(ctx, "select id from balance where user_id = $1", userID).Scan(&balanceID); err != nil {
		return uuid.Nil, fmt.Errorf("select balance err: %w", models.MapStorageErr(err))
	}

	if err = tx.QueryRow(ctx, "select id from orders where order_number = $1", order).Scan(&orderID); err != nil {
		return uuid.Nil, fmt.Errorf("select order err: %w", models.MapStorageErr(err))
	}

	err = g.db.QueryRow(ctx, `insert into balance_operations (order_id, balance_id, sum_operation) values ($1, $2, $3) returning id`, orderID, balanceID, operation).Scan(&bOperationID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("insert balance_operation err: %w", models.MapStorageErr(err))
	}
	tx.Commit(ctx)

	return bOperationID, nil
}

// UpdateBalanceOperation - обновление состояния операции с балансом по номеру заказа
func (g *GophermartDBPostgres) UpdateBalanceOperation(
	ctx context.Context,
	order string,
	operationState BalanceOperationState,
) error {
	var orderID uuid.UUID

	// Транзакция для получения OrderID, чтоб затем использовать его для обновления Balance operation
	tx, err := g.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return fmt.Errorf("transaction err: %w", models.MapStorageErr(err))
	}
	defer tx.Rollback(ctx)

	if err = tx.QueryRow(ctx, "select id from orders where order_number = $1", order).Scan(&orderID); err != nil {
		return fmt.Errorf("select order err: %w", models.MapStorageErr(err))
	}

	_, err = tx.Exec(ctx, "update balance_operations set operation_state = $1, updated_at = now() where order_id = $2", operationState, orderID)
	if err != nil {
		return fmt.Errorf("update balance_operation err: %w", models.MapStorageErr(err))
	}
	tx.Commit(ctx)

	return nil
}

// DeleteBalanceOperation - удаление балансовой операции по номеру заказа
func (g *GophermartDBPostgres) DeleteBalanceOperationByOrderID(
	ctx context.Context,
	order string,
) error {
	var orderID uuid.UUID

	// Транзакция для получения OrderID, чтоб затем использовать его для удаления Balance operation
	tx, err := g.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return fmt.Errorf("transaction err: %w", models.MapStorageErr(err))
	}
	defer tx.Rollback(ctx)

	if err = tx.QueryRow(ctx, "select id from orders where order_number = $1", order).Scan(&orderID); err != nil {
		return fmt.Errorf("select order err: %w", models.MapStorageErr(err))
	}

	_, err = tx.Exec(ctx, "update balance_operations set deleted_at = now() where order_id = $1", orderID)
	if err != nil {
		return fmt.Errorf("update deleted_at balance_operation err: %w", models.MapStorageErr(err))
	}
	tx.Commit(ctx)

	return nil
}

// GetBalanceOperationByOrder - возвращает балансовую операцию по номеру заказа
func (g *GophermartDBPostgres) GetBalanceOperationByOrder(
	ctx context.Context,
	order string,
) (*BalanceOperationItem, error) {
	var orderID uuid.UUID
	var balanceOperation BalanceOperationItem

	// Транзакция для получения OrderID, чтоб затем использовать его для получения Balance operation
	tx, err := g.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return nil, fmt.Errorf("transaction err: %w", models.MapStorageErr(err))
	}
	defer tx.Rollback(ctx)

	if err = tx.QueryRow(ctx, "select id from orders where order_number = $1", order).Scan(&orderID); err != nil {
		return nil, fmt.Errorf("select order err: %w", models.MapStorageErr(err))
	}

	if err = tx.QueryRow(ctx, "select id, order_id, balance_id, sum_operation, updated_at from balance_operations where order_id = $1 and deleted_at is null", orderID).Scan(&balanceOperation.ID, &balanceOperation.OrderID, &balanceOperation.BalanceID, &balanceOperation.SumOperation, &balanceOperation.UpdatedAt); err != nil {
		return nil, fmt.Errorf("select balance operation err: %w", models.MapStorageErr(err))
	}
	tx.Commit(ctx)

	return &balanceOperation, nil
}

// GetBalanceOperation - возвращает все балансовые операции текущему пользователю отсортированных от старых к новым
func (g *GophermartDBPostgres) GetBalanceOperation(
	ctx context.Context,
) ([]*BalanceOperationItem, error) {
	var balanceID uuid.UUID
	userID := ctx.Value(models.KeyUserID)

	// Транзакция для получения BalanceID, чтобы затем использовать его для получения Balance operations
	tx, err := g.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return nil, fmt.Errorf("transaction err: %w", models.MapStorageErr(err))
	}
	defer tx.Rollback(ctx)

	if err = tx.QueryRow(ctx, "select id from balance where user_id = $1", userID).Scan(&balanceID); err != nil {
		return nil, fmt.Errorf("select balance err: %w", models.MapStorageErr(err))
	}

	rows, err := tx.Query(ctx, "select id, order_id, balance_id, sum_operation, updated_at from balance_operations where balance_id = $1 and deleted_at is null order by updated_at", balanceID)
	if err != nil {
		return nil, fmt.Errorf("select balance operation err: %w", models.MapStorageErr(err))
	}

	balanceOperations, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[BalanceOperationItem])
	if err != nil {
		return nil, fmt.Errorf("scan balance operations err: %w", models.MapStorageErr(err))
	}

	tx.Commit(ctx)

	return balanceOperations, nil
}

// GetBalanceByUserID - возвращает баланс текущего пользователя
func (g *GophermartDBPostgres) GetBalanceByUserID(
	ctx context.Context,
) (*BalanceItem, error) {
	var balance BalanceItem
	userID := ctx.Value(models.KeyUserID)

	if err := g.db.QueryRow(ctx, "select id, user_id, current_balance, withdrawn from balance where user_id = $1 and deleted_at is null", userID).Scan(&balance.ID, &balance.UserID, &balance.CurrentBalance, &balance.Withdrawn); err != nil {
		return nil, fmt.Errorf("select balance err: %w", models.MapStorageErr(err))
	}

	return &balance, nil
}

// IncrementBalance - добавляет к текущему балансу указанную сумму
func (g *GophermartDBPostgres) IncrementBalance(
	ctx context.Context,
	incrementSum int64,
) error {
	userID := ctx.Value(models.KeyUserID)

	_, err := g.db.Exec(ctx, "update balance set current_balance = current_balance + $1, updated_at = now() where user_id = $2", incrementSum, userID)
	if err != nil {
		return fmt.Errorf("incerement err: %w", models.MapStorageErr(err))
	}

	return nil
}

// DecrementBalance отнимает от текущего баланса указанную сумму
func (g *GophermartDBPostgres) DecrementBalance(
	ctx context.Context,
	decrementSum int64,
) error {
	userID := ctx.Value(models.KeyUserID)
	var currentSum int64
	var withdrawn int64

	//Транзакция для получения текущего баланса и сравнение с декрментом
	tx, err := g.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return fmt.Errorf("transaction err: %w", models.MapStorageErr(err))
	}
	defer tx.Rollback(ctx)

	if err = tx.QueryRow(ctx, "select current_balance, withdrawn from balance where user_id = $1", userID).Scan(&currentSum, &withdrawn); err != nil {
		return fmt.Errorf("select current sum err: %w", models.MapStorageErr(err))
	}

	if currentSum-decrementSum <= 0 {
		return models.ErrBalanceBelowZero
	}

	_, err = tx.Exec(ctx, "update balance set current_balance = current_balance - $1, withdrawn = withdrawn + $1, updated_at = now() where user_id = $2", decrementSum, userID)
	if err != nil {
		return fmt.Errorf("incerement err: %w", models.MapStorageErr(err))
	}
	tx.Commit(ctx)

	return nil
}
