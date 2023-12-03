package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/service"
)

type UserSuite struct {
	CommonSuite

	users  *service.Users
	userID domain.UserID
}

func TestUsers(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (suite *UserSuite) SetupSuite() {
	suite.CommonSuite.SetupSuite()
	suite.users = service.NewUsers(suite.CommonSuite.db)

	auth := service.NewAuth(suite.CommonSuite.db)

	var err error
	suite.userID, err = auth.SignUp(
		context.Background(),
		domain.Authentication{
			Login:    "login",
			Password: "password",
		},
	)

	suite.Require().NoError(err)
	suite.Require().NotEmpty(suite.userID)
}

func (suite *UserSuite) TestGetBalance() {
	ctx := context.Background()

	suite.Run("success", func() {
		balance, err := suite.users.GetBalance(ctx, suite.userID)
		if suite.NoError(err) {
			suite.Empty(balance)
		}
	})

	suite.Run("not found", func() {
		_, err := suite.users.GetBalance(ctx, domain.EmptyUserID)
		suite.Error(err)
	})
}
