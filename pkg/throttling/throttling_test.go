package throttling_test

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/sergeizaitcev/gophermart/pkg/throttling"
)

var _ throttling.Limiter = (*limiterMock)(nil)

type limiterMock chan struct{}

func newLimiterMock() limiterMock {
	return make(limiterMock, 1)
}

func (m limiterMock) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return m.wait(ctx)
}

func (m limiterMock) wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case m <- struct{}{}:
	}
	return nil
}

func (m limiterMock) Done() {
	<-m
}

var _ http.RoundTripper = (*roundTripperMock)(nil)

type roundTripperMock struct {
	mock.Mock
	limiter limiterMock
}

func newRoundTripperMock(limiter limiterMock) *roundTripperMock {
	return &roundTripperMock{limiter: limiter}
}

func (m *roundTripperMock) RoundTrip(*http.Request) (*http.Response, error) {
	defer m.limiter.Done()
	args := m.Called()
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestTransport(t *testing.T) {
	limiter := newLimiterMock()
	roundTripper := newRoundTripperMock(limiter)
	roundTripper.On("RoundTrip").Return(new(http.Response), nil)

	client := http.Client{Transport: throttling.NewTransport(roundTripper, limiter)}

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	for i := 0; i < 2; i++ {
		res, err := client.Do(req)
		require.NoError(t, err)
		if res != nil {
			io.Copy(io.Discard, res.Body)
			res.Body.Close()
		}
	}

	res, err := client.Do(req.WithContext(ctx))
	require.Error(t, err)
	if res != nil {
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}

	roundTripper.AssertNumberOfCalls(t, "RoundTrip", 2)
}
