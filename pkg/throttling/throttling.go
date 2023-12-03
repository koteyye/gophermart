package throttling

import (
	"context"
	"net/http"
)

// Limiter описывает интерфейс ограничителя запросов.
type Limiter interface {
	// Wait блокируется до тех пор, пока не будет снято ограничение или
	// сработает метод Done у контекста.
	Wait(context.Context) error
}

var _ http.RoundTripper = (*Transport)(nil)

// Transport определяет обёртку над http.RoundTripper, которая органичивает вызов
// метода RoundTrip в течение установленного времени.
type Transport struct {
	roundTripper http.RoundTripper
	limiter      Limiter
}

func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	err := t.limiter.Wait(r.Context())
	if err != nil {
		return nil, err
	}
	return t.roundTripper.RoundTrip(r)
}

// NewTransport возвращает новый экземпляр Transport.
func NewTransport(roundTripper http.RoundTripper, limiter Limiter) *Transport {
	return &Transport{
		roundTripper: roundTripper,
		limiter:      limiter,
	}
}
