package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/storage"
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

	var err error
	suite.userID, err = suite.storage.CreateUser(ctx, "login", "password")
	suite.Require().NoError(err)

	suite.order = "order"
	err = suite.storage.CreateOrder(ctx, suite.userID, suite.order)
	suite.Require().NoError(err)
}

func (suite *OperationSuite) TestA_CreateOperation() {
	ctx := context.Background()

	suite.Run("success", func() {
		err := suite.storage.CreateOperation(ctx, suite.userID, suite.order, 10)
		suite.NoError(err)
	})
}

func (suite *OperationSuite) TestB_GetOperation() {
	ctx := context.Background()

	suite.Run("found", func() {
		want := &storage.Operation{
			UserID:      suite.userID,
			OrderNumber: suite.order,
			Amount:      10,
			Status:      storage.OperationStatusRun,
		}

		got, err := suite.storage.GetOperation(ctx, suite.order)
		if suite.NoError(err) && suite.NotEmpty(got.UpdatedAt) {
			got.UpdatedAt = time.Time{}
			suite.Equal(want, got)
		}
	})

	suite.Run("not_found", func() {
		_, err := suite.storage.GetOperation(ctx, "invalid")
		suite.Error(err)
	})
}

func (suite *OperationSuite) TestB_GetOperations() {
	ctx := context.Background()

	suite.Run("found", func() {
		want := []storage.Operation{{
			UserID:      suite.userID,
			OrderNumber: suite.order,
			Amount:      10,
			Status:      storage.OperationStatusRun,
		}}

		got, err := suite.storage.GetOperations(ctx, suite.userID)
		if suite.NoError(err) && suite.Len(got, len(want)) {
			for i := 0; i < len(want); i++ {
				got[i].UpdatedAt = time.Time{}
			}
			suite.Equal(want, got)
		}
	})

	suite.Run("not_found", func() {
		_, err := suite.storage.GetOperations(ctx, uuid.Nil)
		suite.Error(err)
	})
}

func (suite *OperationSuite) TestC_UpdateOperationStatus() {
	ctx := context.Background()

	suite.Run("success", func() {
		err := suite.storage.UpdateOperationStatus(
			ctx,
			suite.order,
			storage.OperationStatusDone,
		)
		suite.NoError(err)
	})

	suite.Run("get_updated", func() {
		want := &storage.Operation{
			UserID:      suite.userID,
			OrderNumber: suite.order,
			Amount:      10,
			Status:      storage.OperationStatusDone,
		}

		got, err := suite.storage.GetOperation(ctx, suite.order)
		if suite.NoError(err) && suite.NotEmpty(got.UpdatedAt) {
			got.UpdatedAt = time.Time{}
			suite.Equal(want, got)
		}
	})
}

func (suite *OperationSuite) TestD_DeleteOperation() {
	ctx := context.Background()

	suite.Run("success", func() {
		err := suite.storage.DeleteOperation(ctx, suite.order)
		suite.NoError(err)
	})

	suite.Run("get_deleted", func() {
		_, err := suite.storage.GetOperation(ctx, suite.order)
		suite.Error(err)
	})
}
