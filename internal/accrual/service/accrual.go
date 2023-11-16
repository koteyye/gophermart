package service

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"sync"
	"time"

	"github.com/sergeizaitcev/gophermart/internal/accrual/models"
	"github.com/sergeizaitcev/gophermart/internal/accrual/storage"
	"github.com/sergeizaitcev/gophermart/pkg/monetary"
)

// Accrual определяет сервис расчета вознаграждений
type Accrual struct {
	storage storage.Accrual
}

// NewAccrual возвращает новый экземпляр Accrual
func NewAccrual(accrual storage.Accrual) *Accrual {
	return &Accrual{
		storage: accrual,
	}
}

func (a *Accrual) CheckOrder(ctx context.Context, orderNumber string) bool {
	_, err := a.storage.GetOrderWithGoodsByNumber(ctx, orderNumber)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return true
		}
		slog.Error(err.Error())
		return false
	}
	return false
}

// CreateOrder создает заказ с его товарами и отправляет в очередь
func (a *Accrual) CreateOrder(ctx context.Context, order *models.Order) {
	goods := make([]*storage.Goods, len(order.Goods))
	workGoods := make([]*workerGoods, len(order.Goods))
	for i, good := range order.Goods {
		var price monetary.NullUnit

		// Проверяем наличие в БД указанного в заказе match и получаем его ID
		match, err := a.storage.GetMatchByName(ctx, good.Match)
		if err != nil {
			// Если ErrNotFound, то это штатное выполнение сценария
			if errors.Is(err, models.ErrNotFound) {
				slog.Info(err.Error())
				a.storage.CreateInvalidOrder(ctx, order.Number)
				return
			}
			slog.Error(err.Error())
			return
		}

		// Конвертим Price из запроса в значение для БД
		err = price.Scan(good.Price)
		if err != nil {
			slog.Error(err.Error())
		}
		goods[i] = &storage.Goods{MatchID: match.MatchID, Price: float64(price.Unit)}
		workGoods[i] = &workerGoods{matchID: match.MatchID, price: good.Price, reward: match.Reward, rewardType: match.Type}
	}

	orderID, err := a.storage.CreateOrderWithGoods(ctx, order.Number, goods)
	if err != nil {
		slog.Error(err.Error())
	}

	workOrder := workerOrder{orderID: orderID, goods: workGoods}
	a.processing(&workOrder)
}

// processing выполняет процесс по обработке заказа
func (a *Accrual) processing(workOrder *workerOrder) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	doneCh := make(chan struct{})
	workOrderCh := make(chan *workerOrder)

	go func() {
		defer close(workOrderCh)

		select {
		case <-doneCh:
			return
		case workOrderCh <- workOrder:
			err := a.storage.UpdateOrder(ctx, &storage.Order{OrderID: workOrder.orderID, Status: 2, Accrual: 0})
			if err != nil {
				slog.Error(err.Error())
			}
		}
	}()

	calculatedOrderCh := a.calculateAccrual(workOrderCh)
	for order := range calculatedOrderCh {
		for _, good := range order.goods {
			err := a.storage.UpdateGoodAccrual(ctx, good.matchID, float64(good.accrual))
			if err != nil {
				slog.Error(err.Error())
			}
		}
		err := a.storage.UpdateOrder(ctx, &storage.Order{OrderID: order.orderID, Status: 3, Accrual: float64(order.accrual)})
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

// calculateAccrual расчет accrual каждого товара в заказе
func (a *Accrual) calculateAccrual(tasks chan *workerOrder) chan *workerOrder {
	var wg sync.WaitGroup

	out := make(chan *workerOrder)

	wg.Add(len(tasks))
	for task := range tasks {
		for _, good := range task.goods {
			switch good.rewardType {
			case percent:
				accrualResult := good.reward * int(good.price) / 100
				good.accrual = monetary.Unit(accrualResult)
			case natural:
				good.accrual = monetary.Unit(good.reward)
			}
			task.accrual = +good.accrual
		}
		out <- task
	}
	wg.Done()

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
