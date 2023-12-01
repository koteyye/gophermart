package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/golang/mock/gomock"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
)

func (suite *HandlerSuite) TestGetBalance() {
	suite.Run("success", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(nil).Times(1)
		suite.users.EXPECT().GetBalance(gomock.Any(), suite.userID).Return(
			domain.UserBalance{}, nil,
		)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/user/balance", http.NoBody)
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusOK, rec.Code) {
			suite.NotEmpty(rec.Body.String())
			suite.ctrl.Finish()
		}
	})

	suite.Run("unauthorized", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(domain.ErrNotFound).Times(1)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/user/balance", http.NoBody)
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusUnauthorized, rec.Code) {
			suite.ctrl.Finish()
		}
	})

	suite.Run("no token", func() {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/user/balance", http.NoBody)

		suite.handler.ServeHTTP(rec, req)

		suite.Equal(http.StatusUnauthorized, rec.Code)
	})

	suite.Run("not found", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(nil).Times(1)
		suite.users.EXPECT().GetBalance(gomock.Any(), suite.userID).Return(
			domain.UserBalance{}, domain.ErrNotFound,
		)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/user/balance", http.NoBody)
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusNotFound, rec.Code) {
			suite.ctrl.Finish()
		}
	})

	suite.Run("internal server error", func() {
		suite.auth.EXPECT().Identify(gomock.Any(), suite.userID).Return(nil).Times(1)
		suite.users.EXPECT().GetBalance(gomock.Any(), suite.userID).Return(
			domain.UserBalance{}, errors.New("error"),
		)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/user/balance", http.NoBody)
		req.Header.Set("Authorization", "Bearer token")

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusInternalServerError, rec.Code) {
			suite.ctrl.Finish()
		}
	})
}
