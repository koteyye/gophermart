package handler_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/golang/mock/gomock"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
	"github.com/sergeizaitcev/gophermart/pkg/monetary"
)

func (suite *HandlerSuite) TestPerformOperation() {
	orderNumber := "49927398716"

	suite.Run("success", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(nil).Times(1)
		suite.operations.EXPECT().Perform(gomock.Any(), domain.Operation{
			UserID:      suite.userID,
			OrderNumber: domain.OrderNumber(orderNumber),
			Sum:         monetary.Format(1000),
		}).Return(nil)

		body := fmt.Sprintf(`{"order":%q,"sum":1000}`, orderNumber)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPost,
			"/api/user/balance/withdraw",
			strings.NewReader(body),
		)
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusOK, rec.Code) {
			suite.ctrl.Finish()
		}
	})

	suite.Run("invalid order number", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(nil).Times(1)

		body := fmt.Sprintf(`{"order":%q,"sum":1000}`, "invalid")

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPost,
			"/api/user/balance/withdraw",
			strings.NewReader(body),
		)
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusUnprocessableEntity, rec.Code) {
			suite.ctrl.Finish()
		}
	})

	suite.Run("balance below zero", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(nil).Times(1)
		suite.operations.EXPECT().Perform(gomock.Any(), domain.Operation{
			UserID:      suite.userID,
			OrderNumber: domain.OrderNumber(orderNumber),
			Sum:         monetary.Format(1000),
		}).Return(domain.ErrBalanceBelowZero)

		body := fmt.Sprintf(`{"order":%q,"sum":1000}`, orderNumber)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPost,
			"/api/user/balance/withdraw",
			strings.NewReader(body),
		)
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusPaymentRequired, rec.Code) {
			suite.ctrl.Finish()
		}
	})

	suite.Run("internal server error", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(nil).Times(1)
		suite.operations.EXPECT().Perform(gomock.Any(), domain.Operation{
			UserID:      suite.userID,
			OrderNumber: domain.OrderNumber(orderNumber),
			Sum:         monetary.Format(1000),
		}).Return(errors.New("error"))

		body := fmt.Sprintf(`{"order":%q,"sum":1000}`, orderNumber)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPost,
			"/api/user/balance/withdraw",
			strings.NewReader(body),
		)
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusInternalServerError, rec.Code) {
			suite.ctrl.Finish()
		}
	})

	suite.Run("unauthorized", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(domain.ErrNotFound).Times(1)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", http.NoBody)
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusUnauthorized, rec.Code) {
			suite.ctrl.Finish()
		}
	})

	suite.Run("no token", func() {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", http.NoBody)

		suite.handler.ServeHTTP(rec, req)

		suite.Equal(http.StatusUnauthorized, rec.Code)
	})
}

func (suite *HandlerSuite) TestGerOperations() {
	suite.Run("success", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(nil).Times(1)
		suite.operations.EXPECT().GetOperations(gomock.Any(), suite.userID).Return(
			[]domain.Operation{}, nil,
		)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", http.NoBody)
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusOK, rec.Code) {
			suite.NotEmpty(rec.Body.String())
			suite.ctrl.Finish()
		}
	})

	suite.Run("not found", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(nil).Times(1)
		suite.operations.EXPECT().GetOperations(gomock.Any(), suite.userID).Return(
			nil, domain.ErrNotFound,
		)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", http.NoBody)
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusNoContent, rec.Code) {
			suite.ctrl.Finish()
		}
	})

	suite.Run("internal server error", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(nil).Times(1)
		suite.operations.EXPECT().GetOperations(gomock.Any(), suite.userID).Return(
			nil, errors.New("error"),
		)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", http.NoBody)
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusInternalServerError, rec.Code) {
			suite.ctrl.Finish()
		}
	})

	suite.Run("unauthorized", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(domain.ErrNotFound).Times(1)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", http.NoBody)
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusUnauthorized, rec.Code) {
			suite.ctrl.Finish()
		}
	})

	suite.Run("no token", func() {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", http.NoBody)

		suite.handler.ServeHTTP(rec, req)

		suite.Equal(http.StatusUnauthorized, rec.Code)
	})
}
