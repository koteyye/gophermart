package accrual_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/service/accrual"
)

type transportMock struct {
	mock.Mock
	header http.Header
}

func newTransportMock() *transportMock {
	return &transportMock{
		header: make(http.Header),
	}
}

func (m *transportMock) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req.Method, req.URL.EscapedPath())

	err := args.Error(2)
	if err != nil {
		return nil, err
	}

	statusCode := args.Int(0)
	data := args.Get(1)

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	res := &http.Response{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Header:     m.header.Clone(),
		Body:       io.NopCloser(bytes.NewReader(b)),
		Request:    req,
	}

	return res, nil
}

func TestClient(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		mock := newTransportMock()
		client := accrual.NewClient("localhost", &accrual.ClientOption{
			Transport: mock,
		})

		want := accrual.OrderInfo{
			Order:   "1",
			Status:  accrual.StatusRegistered,
			Accrual: 1000,
		}

		mock.On("RoundTrip", "GET", "/api/orders/1").
			Return(http.StatusOK, want, nil)

		got, err := client.OrderInfo(context.Background(), "1")
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("not_registered", func(t *testing.T) {
		mock := newTransportMock()
		client := accrual.NewClient("localhost", &accrual.ClientOption{
			Transport: mock,
		})

		mock.On("RoundTrip", "GET", "/api/orders/2").
			Return(http.StatusNoContent, struct{}{}, nil)

		_, err := client.OrderInfo(context.Background(), "2")
		require.ErrorIs(t, err, accrual.ErrOrderNotRegistered)
	})

	t.Run("too_many_requests", func(t *testing.T) {
		mock := newTransportMock()
		client := accrual.NewClient("localhost", &accrual.ClientOption{
			Transport: mock,
		})

		mock.header.Add("Retry-After", "60")
		mock.On("RoundTrip", "GET", "/api/orders/3").
			Return(http.StatusTooManyRequests, struct{}{}, nil)

		_, err := client.OrderInfo(context.Background(), "3")

		var exhausted *accrual.ResourceExhaustedError
		require.ErrorAs(t, err, &exhausted)
		require.Equal(t, 60*time.Second, exhausted.RetryAfter())
	})

	t.Run("internal_server_error", func(t *testing.T) {
		mock := newTransportMock()
		client := accrual.NewClient("localhost", &accrual.ClientOption{
			Transport: mock,
		})

		mock.On("RoundTrip", "GET", "/api/orders/4").
			Return(http.StatusInternalServerError, struct{}{}, nil)

		_, err := client.OrderInfo(context.Background(), "4")
		require.ErrorIs(t, err, accrual.ErrInternalServerError)
	})
}
