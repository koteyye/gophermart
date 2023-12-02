package httputil

import (
	"io"
	"net/http"
)

// GracefulClose изящно закрывает тело HTTP-ответа.
func GracefulClose(res *http.Response) {
	io.Copy(io.Discard, res.Body)
	_ = res.Body.Close()
}
