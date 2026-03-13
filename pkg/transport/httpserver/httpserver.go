package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	addr       string
}

func NewServer(handler http.Handler, cfg ServerConfig) *Server {
	listenHost := cfg.ListenHost
	if listenHost == "" {
		listenHost = "127.0.0.1"
	}
	addr := fmt.Sprintf("%s:%d", listenHost, cfg.Port)
	timeout := time.Duration(cfg.Timeout) * time.Second

	return &Server{
		addr: addr,
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadTimeout:       timeout,
			ReadHeaderTimeout: timeout,
			WriteTimeout:      timeout,
			IdleTimeout:       timeout,
		},
	}
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
