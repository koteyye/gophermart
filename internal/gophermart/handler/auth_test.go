package handler_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
	"github.com/sergeizaitcev/gophermart/pkg/randutil"
)

func (suite *HandlerSuite) TestRegister() {
	suite.Run("success", func() {
		body := `{"login":"login","password":"password"}`

		suite.auth.EXPECT().SignUp(
			gomock.Any(),
			domain.Authentication{Login: "login", Password: "password"},
		).Return(uuid.New(), nil).Times(1)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(body))

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusOK, rec.Code) {
			suite.NotEmpty(rec.Header().Get("Authorization"))
			suite.ctrl.Finish()
		}
	})

	suite.Run("duplicate", func() {
		body := `{"login":"login","password":"password"}`

		suite.auth.EXPECT().SignUp(
			gomock.Any(),
			domain.Authentication{Login: "login", Password: "password"},
		).Return(uuid.Nil, domain.ErrDuplicate).Times(1)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(body))

		suite.handler.ServeHTTP(rec, req)

		suite.Equal(http.StatusConflict, rec.Code)
		suite.ctrl.Finish()
	})

	suite.Run("internal server error", func() {
		body := `{"login":"login","password":"password"}`

		suite.auth.EXPECT().SignUp(
			gomock.Any(),
			domain.Authentication{Login: "login", Password: "password"},
		).Return(uuid.Nil, errors.New("error")).Times(1)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(body))

		suite.handler.ServeHTTP(rec, req)

		suite.Equal(http.StatusInternalServerError, rec.Code)
		suite.ctrl.Finish()
	})

	suite.Run("bad request", func() {
		body := `[]`

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(body))

		suite.handler.ServeHTTP(rec, req)

		suite.Equal(http.StatusBadRequest, rec.Code)
	})

	suite.Run("auth invalid", func() {
		body := fmt.Sprintf(`{"login":"login","password":%q}`, randutil.String(100))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(body))

		suite.handler.ServeHTTP(rec, req)

		suite.Equal(http.StatusBadRequest, rec.Code)
	})
}

func (suite *HandlerSuite) TestLogin() {
	suite.Run("success", func() {
		body := `{"login":"login","password":"password"}`

		suite.auth.EXPECT().SignIn(
			gomock.Any(),
			domain.Authentication{Login: "login", Password: "password"},
		).Return(uuid.New(), nil).Times(1)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(body))

		suite.handler.ServeHTTP(rec, req)

		if suite.Equal(http.StatusOK, rec.Code) {
			suite.NotEmpty(rec.Header().Get("Authorization"))
			suite.ctrl.Finish()
		}
	})

	suite.Run("not found", func() {
		body := `{"login":"login","password":"password"}`

		suite.auth.EXPECT().SignIn(
			gomock.Any(),
			domain.Authentication{Login: "login", Password: "password"},
		).Return(uuid.Nil, domain.ErrNotFound).Times(1)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(body))

		suite.handler.ServeHTTP(rec, req)

		suite.Equal(http.StatusUnauthorized, rec.Code)
		suite.ctrl.Finish()
	})

	suite.Run("internal server error", func() {
		body := `{"login":"login","password":"password"}`

		suite.auth.EXPECT().SignIn(
			gomock.Any(),
			domain.Authentication{Login: "login", Password: "password"},
		).Return(uuid.Nil, errors.New("error")).Times(1)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(body))

		suite.handler.ServeHTTP(rec, req)

		suite.Equal(http.StatusInternalServerError, rec.Code)
		suite.ctrl.Finish()
	})

	suite.Run("bad request", func() {
		body := `[]`

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(body))

		suite.handler.ServeHTTP(rec, req)

		suite.Equal(http.StatusBadRequest, rec.Code)
	})

	suite.Run("auth invalid", func() {
		body := fmt.Sprintf(`{"login":"login","password":%q}`, randutil.String(100))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(body))

		suite.handler.ServeHTTP(rec, req)

		suite.Equal(http.StatusBadRequest, rec.Code)
	})
}
