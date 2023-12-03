package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/golang/mock/gomock"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
)

func (suite *HandlerSuite) TestOrderProcess() {
	orderNumber := "49927398716"

	suite.Run("success", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(nil).Times(1)
		suite.orders.EXPECT().Process(gomock.Any(), domain.Order{
			UserID: suite.userID,
			Number: domain.OrderNumber(orderNumber),
			Status: domain.OrderStatusNew,
		}).Return(nil)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPost,
			"/api/user/orders",
			strings.NewReader(orderNumber),
		)
		req.Header.Set("Content-Type", "text/plain")
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusAccepted, rec.Code) {
			suite.ctrl.Finish()
		}
	})

	suite.Run("unsupported media type", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(nil).Times(1)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/user/orders", http.NoBody)
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusBadRequest, rec.Code) {
			suite.ctrl.Finish()
		}
	})

	suite.Run("invalid order number", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(nil).Times(1)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPost,
			"/api/user/orders",
			strings.NewReader("invalid"),
		)
		req.Header.Set("Content-Type", "text/plain")
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusUnprocessableEntity, rec.Code) {
			suite.ctrl.Finish()
		}
	})

	suite.Run("duplicate", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(nil).Times(1)
		suite.orders.EXPECT().Process(gomock.Any(), domain.Order{
			UserID: suite.userID,
			Number: domain.OrderNumber(orderNumber),
			Status: domain.OrderStatusNew,
		}).Return(domain.ErrDuplicate)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPost,
			"/api/user/orders",
			strings.NewReader(orderNumber),
		)
		req.Header.Set("Content-Type", "text/plain")
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusOK, rec.Code) {
			suite.ctrl.Finish()
		}
	})

	suite.Run("conflict", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(nil).Times(1)
		suite.orders.EXPECT().Process(gomock.Any(), domain.Order{
			UserID: suite.userID,
			Number: domain.OrderNumber(orderNumber),
			Status: domain.OrderStatusNew,
		}).Return(domain.ErrDuplicateOtherUser)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPost,
			"/api/user/orders",
			strings.NewReader(orderNumber),
		)
		req.Header.Set("Content-Type", "text/plain")
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusConflict, rec.Code) {
			suite.ctrl.Finish()
		}
	})

	suite.Run("internal server error", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(nil).Times(1)
		suite.orders.EXPECT().Process(gomock.Any(), domain.Order{
			UserID: suite.userID,
			Number: domain.OrderNumber(orderNumber),
			Status: domain.OrderStatusNew,
		}).Return(errors.New("error"))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPost,
			"/api/user/orders",
			strings.NewReader(orderNumber),
		)
		req.Header.Set("Content-Type", "text/plain")
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusInternalServerError, rec.Code) {
			suite.ctrl.Finish()
		}
	})

	suite.Run("unauthorized", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(domain.ErrNotFound).Times(1)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/user/orders", http.NoBody)
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusUnauthorized, rec.Code) {
			suite.ctrl.Finish()
		}
	})

	suite.Run("no token", func() {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/user/orders", http.NoBody)

		suite.handler.ServeHTTP(rec, req)

		suite.Equal(http.StatusUnauthorized, rec.Code)
	})
}

func (suite *HandlerSuite) TestGerOrders() {
	suite.Run("success", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(nil).Times(1)
		suite.orders.EXPECT().GetOrders(gomock.Any(), suite.userID).Return(
			[]domain.Order{}, nil,
		)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/user/orders", http.NoBody)
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusOK, rec.Code) {
			suite.NotEmpty(rec.Body.String())
			suite.ctrl.Finish()
		}
	})

	suite.Run("not found", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(nil).Times(1)
		suite.orders.EXPECT().GetOrders(gomock.Any(), suite.userID).Return(
			nil, domain.ErrNotFound,
		)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/user/orders", http.NoBody)
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusNoContent, rec.Code) {
			suite.ctrl.Finish()
		}
	})

	suite.Run("internal server error", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(nil).Times(1)
		suite.orders.EXPECT().GetOrders(gomock.Any(), suite.userID).Return(
			nil, errors.New("error"),
		)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/user/orders", http.NoBody)
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusInternalServerError, rec.Code) {
			suite.ctrl.Finish()
		}
	})

	suite.Run("unauthorized", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(domain.ErrNotFound).Times(1)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/user/orders", http.NoBody)
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusUnauthorized, rec.Code) {
			suite.ctrl.Finish()
		}
	})

	suite.Run("no token", func() {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/user/orders", http.NoBody)

		suite.handler.ServeHTTP(rec, req)

		suite.Equal(http.StatusUnauthorized, rec.Code)
	})
}
