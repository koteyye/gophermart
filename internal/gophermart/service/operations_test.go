package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
	mock_domain "github.com/sergeizaitcev/gophermart/internal/gophermart/domain/mocks"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/service"
	"github.com/sergeizaitcev/gophermart/pkg/monetary"
)

type OperationSuite struct {
	CommonSuite

	operations *service.Operations

	userID      domain.UserID
	orderNumber domain.OrderNumber
}

func TestOperations(t *testing.T) {
	suite.Run(t, new(OperationSuite))
}

func (suite *OperationSuite) SetupSuite() {
	suite.CommonSuite.SetupSuite()

	ctrl := gomock.NewController(suite.T())
	accrual := mock_domain.NewMockAccrualClient(ctrl)

	suite.orderNumber = domain.OrderNumber("49927398716")
	accrual.EXPECT().GetAccrualInfo(gomock.Any(), suite.orderNumber).Return(
		domain.AccrualInfo{
			OrderNumber: suite.orderNumber,
			Status:      domain.AccrualStatusProcessed,
			Accrual:     monetary.Format(2000),
		}, nil,
	).Times(1)

	auth := service.NewAuth(suite.CommonSuite.db)
	suite.operations = service.NewOperations(suite.CommonSuite.db)

	orders := service.NewOrders(suite.CommonSuite.db, accrual)
	defer orders.Close()

	var err error
	ctx := context.Background()

	suite.userID, err = auth.SignUp(
		ctx,
		domain.Authentication{Login: "login", Password: "password"},
	)
	suite.Require().NoError(err)
	suite.Require().NotEmpty(suite.userID)

	err = orders.Process(ctx, domain.Order{
		UserID: suite.userID,
		Number: suite.orderNumber,
	})
	suite.Require().NoError(err)

	time.Sleep(time.Second)
	ctrl.Finish()
}

func (suite *OperationSuite) TestA_Perform() {
	ctx := context.Background()

	suite.Run("perform", func() {
		err := suite.operations.Perform(ctx, domain.Operation{
			UserID:      suite.userID,
			OrderNumber: suite.orderNumber,
			Sum:         monetary.Format(1000),
		})
		suite.NoError(err)
	})

	suite.Run("balance below zero", func() {
		err := suite.operations.Perform(ctx, domain.Operation{
			UserID:      suite.userID,
			OrderNumber: suite.orderNumber,
			Sum:         monetary.Format(10_000),
		})
		suite.ErrorIs(err, domain.ErrBalanceBelowZero)
	})
}

func (suite *OperationSuite) TestB_GetOperations() {
	ctx := context.Background()

	suite.Run("success", func() {
		values, err := suite.operations.GetOperations(ctx, suite.userID)
		if suite.NoError(err) {
			suite.Len(values, 1)
		}
	})

	suite.Run("not found", func() {
		values, err := suite.operations.GetOperations(ctx, domain.EmptyUserID)
		if suite.Error(err) {
			suite.Len(values, 0)
		}
	})
}
