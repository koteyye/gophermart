package storage

import (
	"context"
	"io"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sergeizaitcev/gophermart/pkg/monetary"
)

type Storage interface {
	io.Closer

	Accrual
}

//go:generate mockgen -source=storage.go -destination=mocks/mock.go
// Accrual методы для CRUD в БД
type Accrual interface {
	// CreateOrderWithGoods создание записи о заказе в таблице orders и связанную таблицу goods
	CreateOrderWithGoods(ctx context.Context, order string, goods []*Goods) (uuid.UUID, error)

	// CreateInvalidOrder создает запись о заказе без Goods со статусом Invalid. Для кейсов, когда не найдено соответствующего товара
	CreateInvalidOrder(ctx context.Context, order string) error

	// UpdateOrder обновляет статус и сумму вознаграждения по заказу
	UpdateOrder(ctx context.Context, order *Order) error

	// UpdateGoodAccrual обновляет сумму вознаграждения за конкретный товар
	UpdateGoodAccrual(ctx context.Context, orderID uuid.UUID, matchID uuid.UUID, accrual int) error

	// BatchUpdateGoods обновляет записи в таблице goods по комбинации orderID + mathcID
	BatchUpdateGoods(ctx context.Context, orderID uuid.UUID, goods[]*Goods) error

	// CreateMatch создает новую механику вознаграждения для товара
	CreateMatch(ctx context.Context, match *Match) (uuid.UUID, error)

	// GetMatchByName возвращает механику вознаграждения для товара
	GetMatchByName(ctx context.Context, matchName string) (*MatchOut, error)

	// GetMathesByNames возвращает список механик вознагрждений и их ID по списку имен
	GetMathesByNames(ctx context.Context, matchNames[]string) (map[string]*MatchOut, error)

	// GetOrderByNumber возвращает статус заказа и вознаграждение
	GetOrderByNumber(ctx context.Context, orderNumber string) (*OrderOut, error)
}

// OrderStatus тип статуса заказа
type OrderStatus uint8

const (
	Registered = iota
	Invalid
	Processing
	Processed
)

func (o OrderStatus) String() string {
	switch o {
	case Registered:
		return "registered"
	case Invalid:
		return "invalid"
	case Processing:
		return "processing"
	case Processed:
		return "processed"
	}
	return "unknow"
}

// Goods структура для создания записи в таблице goods
type Goods struct {
	GoodID  uuid.UUID `db:"id"`
	MatchID uuid.UUID `db:"match_id"`
	Price   monetary.Unit   `db:"price"`
	Accrual monetary.Unit   `db:"accrual"`
}

// Order структура для обновления записи в таблице orders
type Order struct {
	OrderID uuid.UUID
	Status  OrderStatus
	Accrual monetary.Unit
}

// OrderOut структура для выгрузки записи из таблицы orders
type OrderOut struct {
	OrderNumber string
	Status      string
	Accrual     monetary.Unit
}

// Match структура для создания записи в таблице matches
type Match struct {
	MatchName string
	Reward    monetary.Unit
	Type      RewardType
}

type RewardType uint8

const (
	Percent = iota
	Natural
)

func (r RewardType) String() string {
	switch r {
	case Percent:
		return "percent"
	case Natural:
		return "natural"
	}
	return "unknow"
}

// MatchOut структура для получения записи из таблицы matches
type MatchOut struct {
	MatchID   uuid.UUID
	MatchName string
	Reward    monetary.Unit
	Type      string
}
