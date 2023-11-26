package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/storage"
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

	var err error
	suite.userID, err = suite.storage.CreateUser(ctx, "login", "password")
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

func (suite *OrderSuite) TestB_GetOrder() {
	ctx := context.Background()

	suite.Run("found", func() {
		want := &storage.Order{
			UserID: suite.userID,
			Number: suite.orders[0],
			Status: storage.OrderStatusNew,
		}

		got, err := suite.storage.GetOrder(ctx, suite.orders[0])

		if suite.NoError(err) && suite.NotEmpty(got.UpdatedAt) {
			got.UpdatedAt = time.Time{}
			suite.Equal(want, got)
		}
	})

	suite.Run("not_found", func() {
		_, err := suite.storage.GetOrder(ctx, "invalid")
		suite.Error(err)
	})
}

func (suite *OrderSuite) TestB_GetOrders() {
	ctx := context.Background()

	suite.Run("found", func() {
		want := []storage.Order{
			{
				UserID: suite.userID,
				Number: suite.orders[0],
				Status: storage.OrderStatusNew,
			},
			{
				UserID: suite.userID,
				Number: suite.orders[1],
				Status: storage.OrderStatusNew,
			},
		}

		got, err := suite.storage.GetOrders(ctx, suite.userID)

		if suite.NoError(err) && suite.Len(got, len(want)) {
			for i := 0; i < len(want); i++ {
				if suite.NotEmpty(got[i].UpdatedAt) {
					got[i].UpdatedAt = time.Time{}
				}
			}
			suite.Equal(want, got)
		}
	})

	suite.Run("not_found", func() {
		_, err := suite.storage.GetOrders(ctx, uuid.Nil)
		suite.Error(err)
	})
}

func (suite *OrderSuite) TestC_UpdateOrder() {
	ctx := context.Background()

	suite.Run("success", func() {
		err := suite.storage.UpdateOrder(ctx, suite.orders[0], storage.OrderStatusProcessed, 10)
		suite.NoError(err)
	})

	suite.Run("get_updated", func() {
		want := &storage.Order{
			UserID:  suite.userID,
			Number:  suite.orders[0],
			Status:  storage.OrderStatusProcessed,
			Accrual: 10,
		}

		got, err := suite.storage.GetOrder(ctx, suite.orders[0])

		if suite.NoError(err) && suite.NotEmpty(got.UpdatedAt) {
			got.UpdatedAt = time.Time{}
			suite.Equal(want, got)
		}
	})
}

func (suite *OrderSuite) TestC_UpdateOrderStatus() {
	ctx := context.Background()

	suite.Run("success", func() {
		err := suite.storage.UpdateOrderStatus(ctx, suite.orders[1], storage.OrderStatusInvalid)
		suite.NoError(err)
	})

	suite.Run("get_updated", func() {
		want := &storage.Order{
			UserID: suite.userID,
			Number: suite.orders[1],
			Status: storage.OrderStatusInvalid,
		}

		got, err := suite.storage.GetOrder(ctx, suite.orders[1])

		if suite.NoError(err) && suite.NotEmpty(got.UpdatedAt) {
			got.UpdatedAt = time.Time{}
			suite.Equal(want, got)
		}
	})
}

func (suite *OrderSuite) TestD_DeleteOrder() {
	ctx := context.Background()

	suite.Run("success", func() {
		err := suite.storage.DeleteOrder(ctx, suite.orders[1])
		suite.NoError(err)
	})

	suite.Run("get_deleted", func() {
		_, err := suite.storage.GetOrder(ctx, suite.orders[1])
		suite.Error(err)
	})
}
