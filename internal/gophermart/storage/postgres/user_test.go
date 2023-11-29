package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/service"
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

	u := service.User{
		Login:    suite.login,
		Password: suite.password,
	}

	suite.Run("success", func() {
		id, err := suite.storage.CreateUser(ctx, u)
		if suite.NoError(err) {
			suite.NotNil(id)
			suite.userID = id
		}
	})

	suite.Run("user_already_exists", func() {
		_, err := suite.storage.CreateUser(ctx, u)
		suite.Error(err)
	})

	suite.Run("hashing_error", func() {
		u.Password = randutil.String(100)
		_, err := suite.storage.CreateUser(ctx, u)
		suite.Error(err)
	})
}

func (suite *UserSuite) TestB_UserID() {
	ctx := context.Background()

	u := service.User{
		Login:    suite.login,
		Password: suite.password,
	}

	suite.Run("found", func() {
		id, err := suite.storage.UserID(ctx, u)
		if suite.NoError(err) {
			suite.Equal(suite.userID, id)
		}
	})

	suite.Run("not_found", func() {
		_, err := suite.storage.UserID(ctx, service.User{})
		suite.Error(err)
	})
}

func (suite *UserSuite) TestB_UserExists() {
	ctx := context.Background()

	suite.Run("found", func() {
		got, err := suite.storage.UserExists(ctx, suite.userID)
		if suite.NoError(err) {
			suite.True(got)
		}
	})

	suite.Run("not_found", func() {
		got, err := suite.storage.UserExists(ctx, uuid.Nil)
		if suite.NoError(err) {
			suite.False(got)
		}
	})
}

func (suite *UserSuite) TestC_GetBalance() {
	ctx := context.Background()

	suite.Run("found", func() {
		want := &service.Balance{}

		got, err := suite.storage.Balance(ctx, suite.userID)
		if suite.NoError(err) {
			suite.Equal(want, got)
		}
	})

	suite.Run("not_found", func() {
		_, err := suite.storage.Balance(ctx, uuid.Nil)
		suite.Error(err)
	})
}
