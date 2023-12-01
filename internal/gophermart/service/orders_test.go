package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
	mock_domain "github.com/sergeizaitcev/gophermart/internal/gophermart/domain/mocks"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/service"
	"github.com/sergeizaitcev/gophermart/pkg/monetary"
)

type OrderSuite struct {
	CommonSuite

	ctrl    *gomock.Controller
	accrual *mock_domain.MockAccrualClient

	orders *service.Orders
	userID domain.UserID
}

func TestOrders(t *testing.T) {
	suite.Run(t, new(OrderSuite))
}

func (suite *OrderSuite) SetupSuite() {
	suite.CommonSuite.SetupSuite()

	suite.ctrl = gomock.NewController(suite.T())
	suite.accrual = mock_domain.NewMockAccrualClient(suite.ctrl)
	suite.orders = service.NewOrders(suite.CommonSuite.db, suite.accrual)

	auth := service.NewAuth(suite.CommonSuite.db)

	var err error
	suite.userID, err = auth.SignUp(
		context.Background(),
		domain.Authentication{Login: "login", Password: "password"},
	)

	suite.Require().NoError(err)
	suite.Require().NotEmpty(suite.userID)
}

func (suite *OrderSuite) TearDownSuite() {
	suite.CommonSuite.TearDownSuite()
	suite.orders.Close()
}

func (suite *OrderSuite) TestA_Process() {
	ctx := context.Background()

	suite.Run("processed", func() {
		order := domain.Order{
			UserID: suite.userID,
			Number: domain.OrderNumber("1"),
		}

		suite.accrual.EXPECT().GetAccrualInfo(gomock.Any(), order.Number).Return(
			domain.AccrualInfo{
				OrderNumber: order.Number,
				Status:      domain.AccrualStatusProcessed,
				Accrual:     monetary.Format(1000),
			}, nil,
		).Times(1)

		err := suite.orders.Process(ctx, order)
		if suite.NoError(err) {
			time.Sleep(time.Second)
			suite.ctrl.Finish()
		}
	})

	suite.Run("processing", func() {
		order := domain.Order{
			UserID: suite.userID,
			Number: domain.OrderNumber("2"),
		}

		suite.accrual.EXPECT().GetAccrualInfo(gomock.Any(), order.Number).Return(
			domain.AccrualInfo{
				OrderNumber: order.Number,
				Status:      domain.AccrualStatusProcessing,
				Accrual:     monetary.Format(1000),
			}, nil,
		).Times(1)

		suite.accrual.EXPECT().GetAccrualInfo(gomock.Any(), order.Number).Return(
			domain.AccrualInfo{
				OrderNumber: order.Number,
				Status:      domain.AccrualStatusProcessed,
				Accrual:     monetary.Format(1000),
			}, nil,
		).Times(1)

		err := suite.orders.Process(ctx, order)
		if suite.NoError(err) {
			time.Sleep(time.Second)
			suite.ctrl.Finish()
		}
	})

	suite.Run("invalid", func() {
		order := domain.Order{
			UserID: suite.userID,
			Number: domain.OrderNumber("3"),
		}

		suite.accrual.EXPECT().GetAccrualInfo(gomock.Any(), order.Number).Return(
			domain.AccrualInfo{
				OrderNumber: order.Number,
				Status:      domain.AccrualStatusInvalid,
			}, nil,
		).Times(1)

		err := suite.orders.Process(ctx, order)
		if suite.NoError(err) {
			time.Sleep(time.Second)
			suite.ctrl.Finish()
		}
	})

	suite.Run("duplicate", func() {
		order := domain.Order{
			UserID: suite.userID,
			Number: domain.OrderNumber("1"),
		}

		err := suite.orders.Process(ctx, order)
		suite.ErrorIs(err, domain.ErrDuplicate)
	})

	suite.Run("duplicate other user", func() {
		order := domain.Order{
			UserID: uuid.New(),
			Number: domain.OrderNumber("2"),
		}

		err := suite.orders.Process(ctx, order)
		suite.ErrorIs(err, domain.ErrDuplicateOtherUser)
	})
}

func (suite *OrderSuite) TestB_GetOrders() {
	ctx := context.Background()

	suite.Run("success", func() {
		values, err := suite.orders.GetOrders(ctx, suite.userID)
		if suite.NoError(err) {
			suite.Len(values, 3)
		}
	})

	suite.Run("not found", func() {
		values, err := suite.orders.GetOrders(ctx, domain.EmptyUserID)
		if suite.Error(err) {
			suite.Len(values, 0)
		}
	})
}
