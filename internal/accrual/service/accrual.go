package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/sergeizaitcev/gophermart/internal/accrual/models"
	"github.com/sergeizaitcev/gophermart/internal/accrual/storage"
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

func (a *Accrual) CheckOrder(ctx context.Context, orderNumber string) (bool, error) {
	_, err := a.storage.GetOrderByNumber(ctx, orderNumber)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return true, nil
		}
		slog.Error(err.Error())
		return false, err
	}
	return false, nil
}

// CreateOrder создает заказ с его товарами и отправляет в очередь
func (a *Accrual) CreateOrder(ctx context.Context, order *models.Order) {
	goods := make([]*storage.Goods, len(order.Goods))
	workGoods := make([]*workerGoods, len(order.Goods))
	for i, good := range order.Goods {
		// Проверяем наличие в БД указанного в заказе match и получаем его ID
		match, err := a.storage.GetMatchByName(ctx, good.Match)
		if err != nil {
			// Если ErrNotFound, то это штатное выполнение сценария
			if errors.Is(err, models.ErrNotFound) {
				slog.Info(fmt.Errorf("%w: %s", err, good.Match).Error())
				a.storage.CreateInvalidOrder(ctx, order.Number)
				return
			}
			slog.Error(err.Error())
			return
		}

		goods[i] = &storage.Goods{MatchID: match.MatchID, Price: good.Price.Float64()}
		workGoods[i] = &workerGoods{matchID: match.MatchID, price: good.Price.Float64(), reward: match.Reward, rewardType: match.Type}
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
	doneCh := make(chan struct{})
	workOrderCh := make(chan *workerOrder)

	go func() {
		defer close(workOrderCh)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

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
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			err := a.storage.UpdateGoodAccrual(ctx, order.orderID, good.matchID, good.accrual)
			if err != nil {
				slog.Error(err.Error())
			}
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

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

	for task := range tasks {
		wg.Add(1)
		go func(task *workerOrder) {
			defer wg.Done()
			for _, good := range task.goods {
				switch good.rewardType {
				case percent:
					accrualResult := good.reward * good.price / 100
					good.accrual = accrualResult
				case natural:
					good.accrual = good.reward
				}
				task.accrual += good.accrual
			}
			out <- task
		}(task)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
