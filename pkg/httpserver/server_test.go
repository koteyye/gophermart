package httpserver_test

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sergeizaitcev/gophermart/pkg/httpserver"
	"github.com/sergeizaitcev/gophermart/pkg/tcputil"
)

func TestListenAndServe(t *testing.T) {
	ping := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)

	srv := &httpserver.Server{
		Handler: mux,
	}

	port, err := tcputil.FreePort()
	require.NoError(t, err)
	addr := net.JoinHostPort("localhost", port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error)
	go func() { errCh <- srv.ListenAndServe(ctx, addr) }()

	u := &url.URL{
		Scheme: "http",
		Host:   addr,
		Path:   "/ping",
	}

	t.Logf("server addr: %s", u.String())

	res, err := http.Get(u.String())
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, res.StatusCode)
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}

	cancel()
	assert.NoError(t, <-errCh)
}
