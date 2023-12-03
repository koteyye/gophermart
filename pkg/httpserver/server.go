package httpserver

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"
)

// Server определяет HTTP-сервер.
type Server struct {
	Handler http.Handler

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	once sync.Once
	srv  *http.Server
}

func (s *Server) lazyInit() {
	s.once.Do(func() {
		s.srv = &http.Server{
			Handler:      s.Handler,
			ReadTimeout:  s.ReadTimeout,
			WriteTimeout: s.WriteTimeout,
			IdleTimeout:  s.IdleTimeout,
		}
	})
}

// ListenAndServe запускает сервер и блокируется до тех пор, пока не сработает
// контекст или метод не вернёт ошибку.
func (s *Server) ListenAndServe(ctx context.Context, addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer l.Close()

	errc := make(chan error)
	go func() { errc <- s.Serve(l) }()

	select {
	case <-ctx.Done():
	case err := <-errc:
		return err
	}

	select {
	case err := <-errc:
		return err
	default:
	}

	err = s.Shutdown()
	<-errc

	return err
}

// Serve запускает сервер, прослушивая входящие соединения при помощи
// net.Listener, и блокируется до тех пор, пока метод не вернёт ошибку.
func (s *Server) Serve(l net.Listener) error {
	s.lazyInit()
	return s.srv.Serve(l)
}

// Shutdown изящно останавливает сервер и закрывает активный net.Listener.
func (s *Server) Shutdown() error {
	s.lazyInit()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.srv.Shutdown(ctx)
}

// ListenAndServe запускает сервер и блокируется до тех пор, пока не сработает
// контекст или функция не вернёт ошибку.
func ListenAndServe(ctx context.Context, addr string, h http.Handler) error {
	s := &Server{
		Handler:      h,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return s.ListenAndServe(ctx, addr)
}
