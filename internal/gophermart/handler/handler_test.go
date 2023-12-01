package handler_test

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
	mock_domain "github.com/sergeizaitcev/gophermart/internal/gophermart/domain/mocks"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/handler"
)

type signerStub struct {
	userID domain.UserID
}

func (*signerStub) Sign(payload string) (token string, err error) {
	return "token", nil
}

func (s *signerStub) Parse(token string) (payload string, err error) {
	return s.userID.String(), nil
}

type HandlerSuite struct {
	suite.Suite

	ctrl *gomock.Controller

	auth       *mock_domain.MockAuthService
	operations *mock_domain.MockOperationService
	orders     *mock_domain.MockOrderService
	users      *mock_domain.MockUserService

	handler http.Handler
	userID  domain.UserID
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

func (suite *HandlerSuite) SetupSuite() {
	suite.ctrl = gomock.NewController(suite.T())

	suite.auth = mock_domain.NewMockAuthService(suite.ctrl)
	suite.operations = mock_domain.NewMockOperationService(suite.ctrl)
	suite.orders = mock_domain.NewMockOrderService(suite.ctrl)
	suite.users = mock_domain.NewMockUserService(suite.ctrl)

	suite.userID = uuid.New()

	suite.handler = handler.New(handler.HandlerOptions{
		Auth:       suite.auth,
		Operations: suite.operations,
		Orders:     suite.orders,
		Users:      suite.users,
		Signer:     &signerStub{userID: suite.userID},
	})
}
