package tcputil_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sergeizaitcev/gophermart/pkg/tcputil"
)

func TestFreePort(t *testing.T) {
	port, err := tcputil.FreePort()
	require.NoError(t, err)

	addr := net.JoinHostPort("", port)

	l, err := net.Listen("tcp", addr)
	require.NoError(t, err)
	require.NoError(t, l.Close())
}
