package tcputil

import (
	"net"
	"strconv"
)

// FreePort запрашивает у ядра свободный открытый порт, готовый
// к использованию.
func FreePort() (port string, err error) {
	a, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "", err
	}

	l, err := net.ListenTCP("tcp", a)
	if err != nil {
		return "", err
	}

	p := l.Addr().(*net.TCPAddr).Port
	_ = l.Close()

	return strconv.FormatInt(int64(p), 10), nil
}
