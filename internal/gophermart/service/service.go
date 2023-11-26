package service

import (
	"context"
	"errors"
	"sync"

	"log/slog"

	"github.com/sergeizaitcev/gophermart/pkg/sign"
)

// ServiceOptions определеяет зависимости gophermart.
type ServiceOptions struct {
	Logger  *slog.Logger
	Accrual AccrualClient
	Storage Storage
	Signer  sign.Signer
}

// Service определяет бизнес-логику gophermart.
type Service struct {
	logger *slog.Logger

	accrual AccrualClient
	storage Storage
	signer  sign.Signer

	orders     OrderQueue
	operations OperationQueue

	wg     *sync.WaitGroup
	termCh chan struct{}
}

// NewService возвращает новый экземпляр Service.
func NewService(opt ServiceOptions) *Service {
	s := &Service{
		logger:  opt.Logger,
		accrual: opt.Accrual,
		storage: opt.Storage,
		signer:  opt.Signer,
		wg:      new(sync.WaitGroup),
		termCh:  make(chan struct{}),
	}

	for i := 0; i < 2; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.orderProcessing()
		}()
	}

	for i := 0; i < 2; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.operationProcessing()
		}()
	}

	return s
}

// Close сигнализирует всем процессам о завершении и блокируется до тех пор,
// пока они не будут завершены.
func (s *Service) Close() {
	if s.closed() {
		return
	}
	close(s.termCh)
	s.wg.Wait()
}

func (s *Service) closed() bool {
	select {
	case <-s.termCh:
		return true
	default:
		return false
	}
}

func (s *Service) withCancel() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	s.wg.Add(1)
	go func() {
		select {
		case <-ctx.Done():
		case <-s.termCh:
			cancel()
		}
		s.wg.Done()
	}()
	return ctx, cancel
}

func (s *Service) retry(fn func(ctx context.Context) error) {
	ctx, cancel := s.withCancel()
	defer cancel()

	for i := 0; i < 2; i++ {
		err := fn(ctx)
		if err == nil || errors.Is(err, context.Canceled) {
			return
		}
	}

	err := fn(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error(err.Error())
	}
}
