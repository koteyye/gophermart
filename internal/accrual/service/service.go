package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/sergeizaitcev/gophermart/internal/accrual/models"
	"github.com/sergeizaitcev/gophermart/internal/accrual/storage"
	"github.com/sergeizaitcev/gophermart/pkg/monetary"
)

// Service определяет бизнес-логику accrual
type Service struct {
	storage storage.Storage
	termCh  chan struct{}
	wg      *sync.WaitGroup
}

// NewService возвращает экземпляр Service
func NewService(s storage.Storage) *Service {
	return &Service{
		termCh:  make(chan struct{}),
		storage: s,
		wg: &sync.WaitGroup{},
	}
}

func (s *Service) Close() {
	close(s.termCh)
	s.wg.Wait()
}

func (s *Service) withCancel() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-ctx.Done():
		case <-s.termCh:
			cancel()
		}
	}()

	return ctx, cancel
}

func (s *Service) GetOrder(ctx context.Context, orderNumber string) (*models.OrderOut, error) {
	order, err := s.storage.GetOrderByNumber(ctx, orderNumber)
	if err != nil {
		slog.Error(fmt.Errorf("get order by number %s err: %w", orderNumber, err).Error())
		return &models.OrderOut{}, err
	}
	return &models.OrderOut{
		Number:  order.OrderNumber,
		Status:  order.Status,
		Accrual: order.Accrual}, nil
}

// CheckMatches проверяет наличие зарегистрированного match в БД
func (s *Service) CheckMatch(ctx context.Context, matchName string) error {
	_, err := s.storage.GetMatchByName(ctx, matchName)
	if err != nil {
		return err
	}

	return nil
}

// CreateMatch создает match в БД
func (s *Service) CreateMatch(ctx context.Context, match *models.Match) error {
	_, err := s.storage.CreateMatch(ctx, &storage.Match{MatchName: match.MatchName, Reward: match.Reward, Type: storage.RewardType(match.RewardType.Uint())})
	if err != nil {
		slog.Error(fmt.Errorf("create match %s err: %w", match.MatchName, err).Error())
		return err
	}
	return nil
}

// CheckOrder проверяет наличие заказа в БД
func (s *Service) CheckOrder(ctx context.Context, orderNumber string) error {
	order, err := s.storage.GetOrderByNumber(ctx, orderNumber)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil
		}
		return err
	}
	if order != nil {
		return ErrOrderRegistered
	}
	return nil
}

// CreateOrder создает заказ с его товарами и отправляет в очередь
func (s *Service) CreateOrder(order *models.Order) {
	ctx, cancel := s.withCancel()
	defer cancel()

	matchNames := make([]string, len(order.Goods))

	for i, good := range order.Goods {
		matchNames[i] = good.Match
	}
	// Проверяем наличие в БД указанного в заказе match и получаем его ID
	matches, err := s.storage.GetMatchesByNames(ctx, matchNames)
	if err != nil {
		slog.Error(err.Error())
	}

	if len(matches) == 0 {
		err := s.storage.CreateInvalidOrder(ctx, order.Number)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		return
	}

	// Заполняем структуры для воркера
	var goods []*storage.Goods
	var workGoods []*workerGoods
	for _, good := range order.Goods {
		if matches[good.Match] == nil {
			continue
		}

		goods = append(goods, &storage.Goods{MatchID: matches[good.Match].MatchID, Price: good.Price})
		workGoods = append(workGoods, &workerGoods{
			matchID:    matches[good.Match].MatchID,
			price:      good.Price,
			reward:     matches[good.Match].Reward.Float64(),
			rewardType: matches[good.Match].Type})
	}

	orderID, err := s.storage.CreateOrderWithGoods(ctx, order.Number, goods)
	if err != nil {
		slog.Error(err.Error())
	}

	workOrder := workerOrder{orderID: orderID, goods: workGoods}
	s.processing(ctx, &workOrder)
}

// processing выполняет процесс по обработке заказа
func (s *Service) processing(ctx context.Context, workOrder *workerOrder) {
	//в бд устанавливаем статус заказа на processing
	s.updateOrderProcessing(ctx, workOrder.orderID)

	//рассчитывается каждый goods в заказе
	workOrderCh := s.calculateOrder(ctx, workOrder)

	for order := range workOrderCh {
		batchGoods := make([]*storage.Goods, len(order.goods))
		for i, good := range order.goods {
			batchGoods[i] = &storage.Goods{MatchID: good.matchID, Price: good.price, Accrual: good.accrual}
		}

		//в бд обновляются рассчитыванные goods в заказе
		err := s.storage.BatchUpdateGoods(ctx, order.orderID, batchGoods)
		if err != nil {
			slog.Error(fmt.Errorf("batch updated goods err: %w", err).Error())
		}

		//в бд обновляем общий accrual по заказу и обновляем статус на processed
		s.updateOrderProcessed(ctx, order)
		close(workOrderCh)
	}
}

// updateOrderProcessing обновляет статус заказа на "processing"
func (s *Service) updateOrderProcessing(ctx context.Context, orderID uuid.UUID) {
	err := s.storage.UpdateOrder(ctx, &storage.Order{OrderID: orderID, Status: 2, Accrual: 0})
	if err != nil {
		slog.Error(fmt.Errorf("update order status processing err: %w", err).Error())
	}
}

func (s *Service) updateOrderProcessed(ctx context.Context, order *workerOrder) {
	err := s.storage.UpdateOrder(ctx, &storage.Order{
		OrderID: order.orderID,
		Status:  3,
		Accrual: order.accrual,
	})
	if err != nil {
		slog.Error(fmt.Errorf("update order status processed err: %w", err).Error())
	}
}

// calculateOrder выполняет расчет каждого goods в заказе
func (s *Service) calculateOrder(ctx context.Context, workOrder *workerOrder) chan *workerOrder {
	workOrderCh := make(chan *workerOrder)

	go func() {
		for _, good := range workOrder.goods {
			good.accrual = calculateAccrual(good)
			workOrder.accrual += good.accrual
		}
		workOrderCh <- workOrder
	}()

	return workOrderCh
}

// calculateAccrual расчет accrual каждого товара в заказе
func calculateAccrual(good *workerGoods) monetary.Unit {
	switch good.rewardType {
	case percent:
		return monetary.Format(good.reward / 100 * good.price.Float64())
	case natural:
		return monetary.Format(good.reward)
	}

	return 0
}
