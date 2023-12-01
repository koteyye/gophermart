package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/service"
	"github.com/sergeizaitcev/gophermart/pkg/randutil"
)

type AuthSuite struct {
	CommonSuite

	auth *service.Auth

	userID         domain.UserID
	authentication domain.Authentication
}

func TestAuth(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}

func (suite *AuthSuite) SetupSuite() {
	suite.CommonSuite.SetupSuite()
	suite.auth = service.NewAuth(suite.CommonSuite.db)
	suite.authentication = domain.Authentication{Login: "login", Password: "password"}
}

func (suite *AuthSuite) TestA_SignUp() {
	ctx := context.Background()

	suite.Run("success", func() {
		var err error
		suite.userID, err = suite.auth.SignUp(ctx, suite.authentication)
		if suite.NoError(err) {
			suite.NotEmpty(suite.userID)
		}
	})

	suite.Run("duplicate", func() {
		_, err := suite.auth.SignUp(ctx, suite.authentication)
		suite.Error(err)
	})

	suite.Run("error", func() {
		invalidPass := randutil.String(100)
		auth := domain.Authentication{
			Login:    suite.authentication.Login,
			Password: invalidPass,
		}
		_, err := suite.auth.SignUp(ctx, auth)
		suite.Error(err)
	})
}

func (suite *AuthSuite) TestB_SignIn() {
	ctx := context.Background()

	suite.Run("success", func() {
		uid, err := suite.auth.SignIn(ctx, suite.authentication)
		if suite.NoError(err) {
			suite.NotEmpty(uid)
		}
	})

	suite.Run("not found by login", func() {
		_, err := suite.auth.SignIn(ctx, domain.Authentication{})
		suite.Error(err)
	})

	suite.Run("not equals", func() {
		_, err := suite.auth.SignIn(ctx, domain.Authentication{Login: suite.authentication.Login})
		suite.Error(err)
	})
}

func (suite *AuthSuite) TestC_Identify() {
	ctx := context.Background()

	suite.Run("success", func() {
		err := suite.auth.Identify(ctx, suite.userID)
		suite.NoError(err)
	})

	suite.Run("not found", func() {
		err := suite.auth.Identify(ctx, domain.EmptyUserID)
		suite.Error(err)
	})
}
