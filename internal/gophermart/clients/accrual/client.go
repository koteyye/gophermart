package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"path"
	"time"

	"log/slog"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/domain"
	"github.com/sergeizaitcev/gophermart/pkg/httputil"
	"github.com/sergeizaitcev/gophermart/pkg/monetary"
)

var defaultOption = &ClientOption{
	Logger:    slog.Default(),
	Transport: http.DefaultTransport.(*http.Transport).Clone(),
	Retry:     3,
	Backoff:   time.Second,
}

// ClientOption определяет не обязательные параметры для Client.
type ClientOption struct {
	// Логирование ошибок.
	Logger *slog.Logger

	// Время ожидания ответа от сервера.
	Timeout time.Duration

	// Пользовательский транспорт.
	Transport http.RoundTripper

	// Максимальное количество попыток выполнить запрос на сервер.
	//
	// По умолчанию 3.
	Retry int

	// Время ожидания между попытками выполнить запрос.
	//
	// По умолчанию 1s.
	Backoff time.Duration

	// Индикатор использования https соединения.
	//
	// По умолчанию false.
	Secure bool
}

func (o *ClientOption) clone() *ClientOption {
	o2 := *o
	return &o2
}

var _ domain.AccrualClient = (*Client)(nil)

// Client определяет HTTP-клиент для запросов в accrual.
type Client struct {
	client *http.Client
	addr   string
	opts   *ClientOption
}

// NewClient возвращает новый экземпляр Client.
func NewClient(addr string, opts *ClientOption) *Client {
	if opts == nil {
		opts = defaultOption
	}
	opts = opts.clone()
	if opts.Logger == nil {
		opts.Logger = defaultOption.Logger
	}
	if opts.Transport == nil {
		opts.Transport = defaultOption.Transport
	}
	if opts.Retry <= 0 {
		opts.Retry = defaultOption.Retry
	}
	if opts.Backoff <= time.Second {
		opts.Backoff = defaultOption.Backoff
	}
	c := &Client{
		client: &http.Client{
			Timeout:   opts.Timeout,
			Transport: opts.Transport,
		},
		addr: addr,
		opts: opts,
	}
	return c
}

// GetAccrualInfo реализует интерфейс domain.AccrualClient.
func (c *Client) GetAccrualInfo(
	ctx context.Context,
	number domain.OrderNumber,
) (domain.AccrualInfo, error) {
	u := c.addr + "/" + path.Join("api", "orders", string(number))

	res, err := c.get(ctx, u)
	if err != nil {
		return domain.AccrualInfo{}, fmt.Errorf("executing a get request: %w", err)
	}
	defer httputil.GracefulClose(res)

	if res.StatusCode != http.StatusOK {
		return domain.AccrualInfo{}, prepareError(res)
	}

	var data accrualData

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return domain.AccrualInfo{}, fmt.Errorf("reading a response body: %w", err)
	}

	return toAccrualInfo(data)
}

type accrualData struct {
	Order   string        `json:"order"`
	Status  string        `json:"status"`
	Accrual monetary.Unit `json:"accrual"`
}

func toAccrualInfo(data accrualData) (domain.AccrualInfo, error) {
	var info domain.AccrualInfo
	orderNumber, err := domain.NewOrderNumber(data.Order)
	if err != nil {
		return info, err
	}
	status, err := domain.NewAccrualStatus(data.Status)
	if err != nil {
		return info, err
	}
	info.OrderNumber = orderNumber
	info.Status = status
	info.Accrual = data.Accrual
	return info, nil
}

func (c *Client) get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creatign a new request: %w", err)
	}

	res, err := c.sendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("sending a request: %w", err)
	}

	return res, nil
}

func (c *Client) sendRequest(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	n := c.opts.Retry

	for n > 0 {
		res, err := c.client.Do(req)
		if err == nil {
			return res, nil
		}

		ne, ok := err.(net.Error)
		if errors.Is(err, io.EOF) || (ok && ne.Timeout()) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.opts.Backoff):
				n--
				continue
			}
		}

		return nil, err
	}

	return nil, errors.New("TODO: add error message")
}
