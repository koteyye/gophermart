package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/service"
)

type OrderSuite struct {
	CommonSuite

	userID uuid.UUID
	orders []string
}

func TestOrder(t *testing.T) {
	suite.Run(t, new(OrderSuite))
}

func (suite *OrderSuite) SetupSuite() {
	suite.CommonSuite.SetupSuite()

	ctx := context.Background()

	u := service.User{
		Login:    "login",
		Password: "password",
	}

	var err error
	suite.userID, err = suite.storage.CreateUser(ctx, u)
	suite.Require().NoError(err)

	suite.orders = []string{"order_1", "order_2"}
}

func (suite *OrderSuite) TestA_CreateOrder() {
	ctx := context.Background()

	suite.Run("success", func() {
		for _, order := range suite.orders {
			err := suite.storage.CreateOrder(ctx, suite.userID, order)
			suite.NoError(err)
		}
	})

	suite.Run("order_already_exists", func() {
		for _, order := range suite.orders {
			err := suite.storage.CreateOrder(ctx, suite.userID, order)
			suite.Error(err)
		}
	})

	suite.Run("created_another_user", func() {
		userID := uuid.New()
		for _, order := range suite.orders {
			err := suite.storage.CreateOrder(ctx, userID, order)
			suite.Error(err)
		}
	})
}

func (suite *OrderSuite) TestB_Orders() {
	ctx := context.Background()

	suite.Run("found", func() {
		want := []service.Order{
			{
				Number: suite.orders[0],
				Status: service.OrderStatusNew,
			},
			{
				Number: suite.orders[1],
				Status: service.OrderStatusNew,
			},
		}

		got, err := suite.storage.Orders(ctx, suite.userID)

		if suite.NoError(err) && suite.Len(got, len(want)) {
			for i := 0; i < len(want); i++ {
				if suite.NotEmpty(got[i].UploadedAt) {
					got[i].UploadedAt = time.Time{}
				}
			}
			suite.Equal(want, got)
		}
	})

	suite.Run("not_found", func() {
		_, err := suite.storage.Orders(ctx, uuid.Nil)
		suite.Error(err)
	})
}

func (suite *OrderSuite) TestC_ProcessOrder() {
	ctx := context.Background()

	suite.Run("success", func() {
		err := suite.storage.ProcessOrder(ctx, suite.orders[0], 10)
		suite.NoError(err)
	})
}

func (suite *OrderSuite) TestC_UpdateOrderStatus() {
	ctx := context.Background()

	suite.Run("success", func() {
		err := suite.storage.UpdateOrderStatus(ctx, suite.orders[1], service.OrderStatusInvalid)
		suite.NoError(err)
	})
}
