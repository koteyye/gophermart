package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/storage"
	"github.com/sergeizaitcev/gophermart/pkg/monetary"
	"github.com/sergeizaitcev/gophermart/pkg/randutil"
)

type UserSuite struct {
	CommonSuite

	userID   uuid.UUID
	login    string
	password string
}

func TestUser(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (suite *UserSuite) SetupSuite() {
	suite.CommonSuite.SetupSuite()
	suite.login = "login"
	suite.password = "password"
}

func (suite *UserSuite) TestA_CreateUser() {
	ctx := context.Background()

	suite.Run("success", func() {
		id, err := suite.storage.CreateUser(ctx, suite.login, suite.password)
		if suite.NoError(err) {
			suite.NotNil(id)
			suite.userID = id
		}
	})

	suite.Run("user_already_exists", func() {
		_, err := suite.storage.CreateUser(ctx, suite.login, suite.password)
		suite.Error(err)
	})

	suite.Run("hashing_error", func() {
		invalid := randutil.String(100)
		_, err := suite.storage.CreateUser(ctx, suite.login, invalid)
		suite.Error(err)
	})
}

func (suite *UserSuite) TestB_GetUser() {
	ctx := context.Background()

	suite.Run("found", func() {
		id, err := suite.storage.GetUser(ctx, suite.login, suite.password)
		if suite.NoError(err) {
			suite.Equal(suite.userID, id)
		}
	})

	suite.Run("not_found", func() {
		_, err := suite.storage.GetUser(ctx, "", "")
		suite.Error(err)
	})
}

func (suite *UserSuite) TestB_GetLogin() {
	ctx := context.Background()

	suite.Run("found", func() {
		login, err := suite.storage.GetLogin(ctx, suite.userID)
		if suite.NoError(err) {
			suite.Equal(suite.login, login)
		}
	})

	suite.Run("not_found", func() {
		_, err := suite.storage.GetLogin(ctx, uuid.Nil)
		suite.Error(err)
	})
}

func (suite *UserSuite) TestC_GetBalance() {
	ctx := context.Background()

	suite.Run("found", func() {
		want := &storage.UserBalance{
			UserID: suite.userID,
		}

		got, err := suite.storage.GetBalance(ctx, suite.userID)
		if suite.NoError(err) {
			suite.Equal(want, got)
		}
	})

	suite.Run("not_found", func() {
		_, err := suite.storage.GetBalance(ctx, uuid.Nil)
		suite.Error(err)
	})
}

func (suite *UserSuite) TestD_UpdateBalance() {
	ctx := context.Background()

	suite.Run("update", func() {
		units := []monetary.Unit{0, 100, -50}
		for _, unit := range units {
			suite.NoError(suite.storage.UpdateBalance(ctx, suite.userID, unit))
		}
	})

	suite.Run("get_updated", func() {
		want := &storage.UserBalance{
			UserID:    suite.userID,
			Amount:    50,
			Withdrawn: 50,
		}

		got, err := suite.storage.GetBalance(ctx, suite.userID)
		if suite.NoError(err) {
			suite.Equal(want, got)
		}
	})

	suite.Run("below_zero", func() {
		suite.Error(suite.storage.UpdateBalance(ctx, suite.userID, -51))
	})
}
