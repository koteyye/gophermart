//go:build integration

package e2e

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
	"github.com/sergeizaitcev/gophermart/pkg/httputil"
)

var validOrderNumber = "49927398716"

type E2ESuite struct {
	suite.Suite
	authorization string
}

func TestE2E(t *testing.T) {
	suite.Run(t, new(E2ESuite))
}

func (suite *E2ESuite) TestA_RegisterMechanic() {
	body := `{
		"match": "Tefal",
		"reward": 30,
		"reward_type": "%"
	}`

	url := "http://localhost:8080/api/goods"

	suite.Run("success", func() {
		res, err := http.Post(url, "application/json", strings.NewReader(body))
		if suite.NoError(err) {
			suite.Equal(http.StatusOK, res.StatusCode)
			httputil.GracefulClose(res)
		}
	})

	suite.Run("conflict", func() {
		res, err := http.Post(url, "application/json", strings.NewReader(body))
		if suite.NoError(err) {
			suite.Equal(http.StatusConflict, res.StatusCode)
			httputil.GracefulClose(res)
		}
	})
}

func (suite *E2ESuite) TestB_RegisterOrder() {
	body := `{
		"order": "` + validOrderNumber + `",
		"goods": [
			{
				"description": "Чайник Tefal",
				"price": 2499.99
			}
		]
	}`

	url := "http://localhost:8080/api/orders"

	suite.Run("success", func() {
		res, err := http.Post(url, "application/json", strings.NewReader(body))
		if suite.NoError(err) {
			suite.Equal(http.StatusAccepted, res.StatusCode)
			httputil.GracefulClose(res)
		}
	})

	suite.Run("conflict", func() {
		res, err := http.Post(url, "application/json", strings.NewReader(body))
		if suite.NoError(err) {
			suite.Equal(http.StatusConflict, res.StatusCode)
			httputil.GracefulClose(res)
		}
	})
}

func (suite *E2ESuite) TestC_RegisterUser() {
	body := `{"login":"login","password":"password"}`
	url := "http://localhost:8081/api/user/register"

	suite.Run("success", func() {
		res, err := http.Post(url, "application/json", strings.NewReader(body))
		if suite.NoError(err) {
			suite.Equal(http.StatusOK, res.StatusCode)
			suite.authorization = res.Header.Get("Authorization")
			suite.NotEmpty(suite.authorization)
			httputil.GracefulClose(res)
		}
	})

	suite.Run("conflict", func() {
		res, err := http.Post(url, "application/json", strings.NewReader(body))
		if suite.NoError(err) {
			suite.Equal(http.StatusConflict, res.StatusCode)
			httputil.GracefulClose(res)
		}
	})
}

func (suite *E2ESuite) TestD_Login() {
	url := "http://localhost:8081/api/user/login"

	suite.Run("success", func() {
		body := `{"login":"login","password":"password"}`
		res, err := http.Post(url, "application/json", strings.NewReader(body))
		if suite.NoError(err) {
			suite.Equal(http.StatusOK, res.StatusCode)
			authorization := res.Header.Get("Authorization")
			suite.Equal(suite.authorization, authorization)
			httputil.GracefulClose(res)
		}
	})

	suite.Run("unauthorized", func() {
		body := `{"login":"login2","password":"password"}`
		res, err := http.Post(url, "application/json", strings.NewReader(body))
		if suite.NoError(err) {
			suite.Equal(http.StatusUnauthorized, res.StatusCode)
			httputil.GracefulClose(res)
		}
	})
}

func (suite *E2ESuite) TestF_GetBalance() {
	url := "http://localhost:8081/api/user/balance"

	suite.Run("success", func() {
		req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
		suite.NoError(err)

		req.Header.Set("Authorization", suite.authorization)

		res, err := http.DefaultClient.Do(req)
		if suite.NoError(err) {
			suite.Equal(http.StatusOK, res.StatusCode)

			var balance domain.UserBalance
			suite.NoError(json.NewDecoder(res.Body).Decode(&balance))
			suite.Empty(balance)

			httputil.GracefulClose(res)
		}
	})

	suite.Run("unauthorized", func() {
		req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
		suite.NoError(err)

		res, err := http.DefaultClient.Do(req)
		if suite.NoError(err) {
			suite.Equal(http.StatusUnauthorized, res.StatusCode)
			httputil.GracefulClose(res)
		}
	})
}

func (suite *E2ESuite) TestG_CreateOrder() {
	body := validOrderNumber
	url := "http://localhost:8081/api/user/orders"

	suite.Run("success", func() {
		req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
		suite.NoError(err)

		req.Header.Set("Authorization", suite.authorization)
		req.Header.Set("Content-Type", "text/plain")

		res, err := http.DefaultClient.Do(req)
		if suite.NoError(err) {
			suite.Equal(http.StatusAccepted, res.StatusCode)
			httputil.GracefulClose(res)
		}
	})

	suite.Run("duplicate", func() {
		req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
		suite.NoError(err)

		req.Header.Set("Authorization", suite.authorization)
		req.Header.Set("Content-Type", "text/plain")

		res, err := http.DefaultClient.Do(req)
		if suite.NoError(err) {
			suite.Equal(http.StatusOK, res.StatusCode)
			httputil.GracefulClose(res)
		}
	})

	suite.Run("unauthorized", func() {
		req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
		suite.NoError(err)

		req.Header.Set("Content-Type", "text/plain")

		res, err := http.DefaultClient.Do(req)
		if suite.NoError(err) {
			suite.Equal(http.StatusUnauthorized, res.StatusCode)
			httputil.GracefulClose(res)
		}
	})
}

func (suite *E2ESuite) TestH_Withdrawal() {
	url := "http://localhost:8081/api/user/balance/withdraw"
	body := `{"order":"` + validOrderNumber + `","sum":500}`

	suite.Run("success", func() {
		req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
		suite.NoError(err)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", suite.authorization)

		var code int

		for i := 0; i < 3; i++ {
			res, err := http.DefaultClient.Do(req)
			defer httputil.GracefulClose(res)

			code = res.StatusCode
			if suite.NoError(err) && http.StatusOK == code {
				return
			}

			time.Sleep(time.Second)
		}

		suite.Equal(http.StatusOK, code)
	})
}

func (suite *E2ESuite) TestI_GetBalance() {
	url := "http://localhost:8081/api/user/balance"

	suite.Run("success", func() {
		req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
		suite.NoError(err)

		req.Header.Set("Authorization", suite.authorization)

		res, err := http.DefaultClient.Do(req)
		if suite.NoError(err) {
			suite.Equal(http.StatusOK, res.StatusCode)

			var balance domain.UserBalance
			suite.NoError(json.NewDecoder(res.Body).Decode(&balance))
			suite.NotEmpty(balance)

			httputil.GracefulClose(res)
		}
	})
}
