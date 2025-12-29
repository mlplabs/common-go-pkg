package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

const (
	readWriteTimeoutSec = 60
	shtdwnTimeoutSrc    = 5

	defaultReadTimeout     = readWriteTimeoutSec * time.Second
	defaultWriteTimeout    = readWriteTimeoutSec * time.Second
	defaultAddr            = ":80"
	defaultShutdownTimeout = shtdwnTimeoutSrc * time.Second
)

// Server - HTTP-сервер.
type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func NewServer(handler *chi.Mux, opts ...Option) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         defaultAddr,
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	s.start()

	return s
}

func (s *Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

// Notify ...
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown ...
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
