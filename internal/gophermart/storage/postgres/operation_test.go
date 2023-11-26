package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/service"
)

type OperationSuite struct {
	CommonSuite

	userID uuid.UUID
	order  string
}

func TestOperation(t *testing.T) {
	suite.Run(t, new(OperationSuite))
}

func (suite *OperationSuite) SetupSuite() {
	suite.CommonSuite.SetupSuite()

	ctx := context.Background()

	u := service.User{
		Login:    "login",
		Password: "password",
	}

	var err error
	suite.userID, err = suite.storage.CreateUser(ctx, u)
	suite.Require().NoError(err)

	suite.order = "order"
	err = suite.storage.CreateOrder(ctx, suite.userID, suite.order)
	suite.Require().NoError(err)
}

func (suite *OperationSuite) TestA_CreateOperation() {
	ctx := context.Background()

	suite.Run("success", func() {
		_, err := suite.storage.CreateOperation(ctx, suite.userID, suite.order, 10)
		suite.NoError(err)
	})
}

func (suite *OperationSuite) TestB_Operations() {
	ctx := context.Background()

	suite.Run("found", func() {
		want := []service.Operation{{
			Order:  suite.order,
			Sum:    10,
			Status: service.OperationStatusRun,
		}}

		got, err := suite.storage.Operations(ctx, suite.userID)
		if suite.NoError(err) && suite.Len(got, len(want)) {
			for i := 0; i < len(want); i++ {
				got[i].ProcessedAt = time.Time{}
			}
			suite.Equal(want, got)
		}
	})

	suite.Run("not_found", func() {
		_, err := suite.storage.Operations(ctx, uuid.Nil)
		suite.Error(err)
	})
}

func (suite *OperationSuite) TestC_UpdateOperationStatus() {
	ctx := context.Background()

	suite.Run("success", func() {
		err := suite.storage.UpdateOperationStatus(
			ctx,
			suite.order,
			service.OperationStatusDone,
		)
		suite.NoError(err)
	})
}

func (suite *OperationSuite) TestC_PerformOperation() {
	// TODO
}
