package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/config"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

func New(handler http.Handler, cfg config.HttpConfig, logger *zap.Logger) *Server {
	timeout := time.Duration(cfg.Timeout) * time.Second
	addr := fmt.Sprintf(":%d", cfg.Port)
	return &Server{
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadTimeout:       timeout,
			ReadHeaderTimeout: timeout,
			WriteTimeout:      timeout,
			IdleTimeout:       timeout,
		},
		logger: logger,
	}
}

func (s *Server) ListenAndServe() error {

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {

	return s.httpServer.Shutdown(ctx)
}
