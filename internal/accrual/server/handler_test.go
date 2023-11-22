package server_test

import (
	"errors"
	"github.com/google/uuid"
	"github.com/sergeizaitcev/gophermart/internal/accrual/storage"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

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

	s := mockStorage.NewMockStorage(c)
	srv := service.NewService(s)
	handler := server.NewHandler(srv)

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
				r.EXPECT().GetOrderByNumber(gomock.Any(), order).Return(&storage.OrderOut{}, errors.New("other err"))
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
	type mockBehavior func(r *mockStorage.MockStorage, match *storage.Match)
	type mockBehavior2 func(r *mockStorage.MockStorage, order string)

	tests := []struct {
		name               string
		requestBody        io.Reader
		mockBehavior       mockBehavior
		mockBehavior2      mockBehavior2
		expectedStatusCode int
	}{
		{
			name: "success",
			requestBody: strings.NewReader(`{
				"match": "testMatch",
				"reward": 20,
				"reward_type": "%"
				}`),
			mockBehavior: func(r *mockStorage.MockStorage, match *storage.Match) {
				r.EXPECT().CreateMatch(gomock.Any(), match).Return(uuid.New(), nil)
			},
			mockBehavior2: func(r *mockStorage.MockStorage, matchName string) {
				r.EXPECT().GetMatchByName(gomock.Any(), matchName).Return(&storage.MatchOut{}, storage.ErrNotFound)
			},
			expectedStatusCode: 200,
		},
		{
			name: "bad request",
			requestBody: strings.NewReader(`{
				"match": "testMatch",
				"reward": 20,
				"reward_type": "процентики"
				}`),
			mockBehavior: func(r *mockStorage.MockStorage, match *storage.Match) {
				r.EXPECT().CreateMatch(gomock.Any(), match).Return(uuid.New(), nil)
			},
			mockBehavior2: func(r *mockStorage.MockStorage, matchName string) {
				r.EXPECT().GetMatchByName(gomock.Any(), matchName).Return(&storage.MatchOut{}, storage.ErrNotFound)
			},
			expectedStatusCode: 400,
		},
		{
			name: "duplicate",
			requestBody: strings.NewReader(`{
				"match": "testMatch",
				"reward": 20,
				"reward_type": "%"
				}`),
			mockBehavior: func(r *mockStorage.MockStorage, match *storage.Match) {
				r.EXPECT().CreateMatch(gomock.Any(), match).Return(uuid.Nil, storage.ErrDuplicate)
			},
			mockBehavior2: func(r *mockStorage.MockStorage, matchName string) {
				r.EXPECT().GetMatchByName(gomock.Any(), matchName).Return(&storage.MatchOut{
					MatchID:   uuid.New(),
					MatchName: "testMatch",
					Reward:    10000,
					Type:      "percent",
				}, nil)
			},
			expectedStatusCode: 409,
		},
		{
			name: "internal err",
			requestBody: strings.NewReader(`{
				"match": "testMatch",
				"reward": 20,
				"reward_type": "%"
				}`),
			mockBehavior: func(r *mockStorage.MockStorage, match *storage.Match) {
				r.EXPECT().CreateMatch(gomock.Any(), match).Return(uuid.Nil, errors.New("other err"))
			},
			mockBehavior2: func(r *mockStorage.MockStorage, matchName string) {
				r.EXPECT().GetMatchByName(gomock.Any(), matchName).Return(&storage.MatchOut{
					MatchID:   uuid.New(),
					MatchName: "testMatch",
					Reward:    10000,
					Type:      "percent",
				}, nil)
			},
			expectedStatusCode: 500,
		},
	}

	h, s := testInitHandle(t)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			trgURL, err := url.JoinPath(baseURL, createMatch)
			assert.NoError(t, err)

			r := httptest.NewRequest(http.MethodPost, trgURL, test.requestBody)
			w := httptest.NewRecorder()

			testMatch := storage.Match{}

			test.mockBehavior(s, &testMatch)

			h.ServeHTTP(w, r)

			assert.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}

func TestRegisterOrder(t *testing.T) {
	type mockBehavior func(r *mockStorage.MockStorage, order string)

	requestBody := strings.NewReader(`{
		"order": "1234567812345670",
		"goods": [
			{
				"description": "item1",
				"price": 100
			},
			{
				"description": "item2",
				"price": 200.99
			}
		]
		}`)

	tests := []struct {
		name               string
		requestBody        io.Reader
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name:        "accept",
			requestBody: requestBody,
			mockBehavior: func(r *mockStorage.MockStorage, order string) {
				r.EXPECT().GetOrderByNumber(gomock.Any(), order).Return(&storage.OrderOut{}, storage.ErrNotFound)
			},
			expectedStatusCode: 202,
		},
		{
			name:        "bad request",
			requestBody: strings.NewReader(`{"order": "1234567812345670"}`),
			mockBehavior: func(r *mockStorage.MockStorage, order string) {
				r.EXPECT().GetOrderByNumber(gomock.Any(), order).Return(&storage.OrderOut{OrderNumber: testOrderNum, Status: "processed", Accrual: 10000}, nil)
			},
			expectedStatusCode: 400,
		},
		{
			name:        "duplicate",
			requestBody: requestBody,
			mockBehavior: func(r *mockStorage.MockStorage, order string) {
				r.EXPECT().GetOrderByNumber(gomock.Any(), order).Return(&storage.OrderOut{
					OrderNumber: testOrderNum,
					Status:      "processed",
					Accrual:     10000,
				}, nil)
			},
			expectedStatusCode: 409,
		},
		{
			name:        "internal err",
			requestBody: requestBody,
			mockBehavior: func(r *mockStorage.MockStorage, order string) {
				r.EXPECT().GetOrderByNumber(gomock.Any(), order).Return(&storage.OrderOut{}, errors.New("other err"))
			},
			expectedStatusCode: 500,
		},
	}

	h, s := testInitHandle(t)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			trgURL, err := url.JoinPath(baseURL, register)
			assert.NoError(t, err)

			r := httptest.NewRequest(http.MethodPost, trgURL, test.requestBody)
			w := httptest.NewRecorder()

			test.mockBehavior(s, testOrderNum)

			testMap := make(map[string]*storage.MatchOut)
			testMap["item1"] = &storage.MatchOut{
				MatchID:   uuid.New(),
				MatchName: "item1",
				Reward:    10,
				Type:      "percent",
			}
			testMap["item2"] = &storage.MatchOut{
				MatchID:   uuid.New(),
				MatchName: "item2",
				Reward:    10,
				Type:      "percent",
			}

			s.EXPECT().GetMatchesByNames(gomock.Any(), []string{"item1", "item2"}).Return(testMap, nil)

			h.ServeHTTP(w, r)

			assert.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}
