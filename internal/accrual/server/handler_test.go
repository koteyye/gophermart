package server_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/sergeizaitcev/gophermart/internal/accrual/models"
	"github.com/sergeizaitcev/gophermart/internal/accrual/storage"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/sergeizaitcev/gophermart/internal/accrual/server"
	"github.com/sergeizaitcev/gophermart/internal/accrual/service"
	mockStorage "github.com/sergeizaitcev/gophermart/internal/accrual/storage/mocks"
)

const (
	baseURL      = "http://localhost:8080"
	register     = "/api/orders"
	createMatch  = "/api/goods"
	getOrder     = "/api/orders/"
	testOrderNum = "1234567812345670"
)

func testInitHandle(t *testing.T) (http.Handler, *mockStorage.MockStorage) {
	c := gomock.NewController(t)
	defer c.Finish()
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	testHandler := slog.NewTextHandler(os.Stdout, opts)
	testLogger := slog.New(testHandler)

	s := mockStorage.NewMockStorage(c)
	srv := service.NewService(s)
	handler := server.NewHandler(testLogger, srv)

	return handler, s
}

func TestGetOrder(t *testing.T) {
	type mockBehavior func(r *mockStorage.MockStorage, order string)
	tests := []struct {
		name                 string
		order                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:  "success",
			order: "1234567812345670",
			mockBehavior: func(r *mockStorage.MockStorage, order string) {
				r.EXPECT().GetOrderByNumber(gomock.Any(), order).Return(&storage.OrderOut{
					OrderNumber: testOrderNum,
					Status:      "processed",
					Accrual:     10000,
				}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `{
				"order": "1234567812345670",
				"status": "processed",
				"accrual": 100
				}`,
		},
		{
			name:  "no order",
			order: testOrderNum,
			mockBehavior: func(r *mockStorage.MockStorage, order string) {
				r.EXPECT().GetOrderByNumber(gomock.Any(), order).Return(&storage.OrderOut{}, storage.ErrNotFound)
			},
			expectedStatusCode: 404,
		},
		{
			name:  "internal err",
			order: testOrderNum,
			mockBehavior: func(r *mockStorage.MockStorage, order string) {
				r.EXPECT().GetOrderByNumber(gomock.Any(), order).Return(&storage.OrderOut{}, storage.ErrOther)
			},
			expectedStatusCode: 500,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h, s := testInitHandle(t)

			trgURL, err := url.JoinPath(baseURL, getOrder, test.order)
			assert.NoError(t, err)

			r := httptest.NewRequest(http.MethodGet, trgURL, nil)
			w := httptest.NewRecorder()

			test.mockBehavior(s, test.order)

			h.ServeHTTP(w, r)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			if test.expectedResponseBody != "" {
				assert.JSONEq(t, w.Body.String(), test.expectedResponseBody)
			}
		})
	}
}

func TestCreateMatch(t *testing.T) {
	testMatch := storage.Match{
		MatchName: "testMatch",
		Reward:    2000,
		Type:      0,
	}

	testRequest := strings.NewReader(`{
		"match": "testMatch",
		"reward": 20,
		"reward_type": "%"
		}`)

	trgURL, err := url.JoinPath(baseURL, createMatch)
	assert.NoError(t, err)

	t.Run("createMatch", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			t.Parallel()

			h, s := testInitHandle(t)

			r := httptest.NewRequest(http.MethodPost, trgURL, testRequest)
			w := httptest.NewRecorder()

			s.EXPECT().GetMatchByName(gomock.Any(), gomock.Any()).Return(&storage.MatchOut{}, storage.ErrNotFound)
			s.EXPECT().CreateMatch(gomock.Any(), &testMatch).Return(uuid.New(), (error)(nil))

			h.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("badRequest", func(t *testing.T) {
			t.Parallel()

			h, _ := testInitHandle(t)

			testBadRequest := strings.NewReader(`{
				"match": "testMatch",
				"reward": 20,
				"reward_type": "процентики"
				}`)

			r := httptest.NewRequest(http.MethodPost, trgURL, testBadRequest)
			w := httptest.NewRecorder()

			h.ServeHTTP(w, r)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	})
}

func TestRegisterOrder(t *testing.T) {
	trgURL, err := url.JoinPath(baseURL, register)
	assert.NoError(t, err)

	testRequest := models.Order{
		Number: testOrderNum,
		Goods: []models.Goods{
			{
				Match: "item1",
				Price: 100,
			},
			{
				Match: "item2",
				Price: 20099,
			},
		},
	}

	testRequestBody, err := json.Marshal(testRequest)
	assert.NoError(t, err)

	t.Run("createOrder", func(t *testing.T) {
		t.Run("accept", func(t *testing.T) {
			t.Parallel()
			h, s := testInitHandle(t)

			r := httptest.NewRequest(http.MethodPost, trgURL, bytes.NewReader(testRequestBody))
			w := httptest.NewRecorder()

			s.EXPECT().GetOrderByNumber(gomock.Any(), gomock.Eq(testRequest.Number)).Return(&storage.OrderOut{}, storage.ErrNotFound)

			h.ServeHTTP(w, r)

			assert.Equal(t, http.StatusAccepted, w.Code)
		})

		t.Run("conflict", func(t *testing.T) {
			h, s := testInitHandle(t)

			r := httptest.NewRequest(http.MethodPost, trgURL, bytes.NewReader(testRequestBody))
			w := httptest.NewRecorder()

			s.EXPECT().GetOrderByNumber(gomock.Any(), gomock.Eq(testRequest.Number)).Return(&storage.OrderOut{}, storage.ErrDuplicate)

			h.ServeHTTP(w, r)

			assert.Equal(t, http.StatusConflict, w.Code)
		})

		t.Run("internal", func(t *testing.T) {
			h, s := testInitHandle(t)

			r := httptest.NewRequest(http.MethodPost, trgURL, bytes.NewReader(testRequestBody))
			w := httptest.NewRecorder()

			s.EXPECT().GetOrderByNumber(gomock.Any(), gomock.Eq(testRequest.Number)).Return(&storage.OrderOut{}, storage.ErrOther)

			h.ServeHTTP(w, r)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("badRequest", func(t *testing.T) {
			h, _ := testInitHandle(t)

			r := httptest.NewRequest(http.MethodPost, trgURL, strings.NewReader(`{"order": "1234567812345670"}`))
			w := httptest.NewRecorder()

			h.ServeHTTP(w, r)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	})
}
