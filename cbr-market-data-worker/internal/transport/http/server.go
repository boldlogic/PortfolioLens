package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/config"
	"github.com/sirupsen/logrus"
)

type Server struct {
	httpServer *http.Server
	addr       string
	log        logrus.FieldLogger
}

func NewServer(handler http.Handler, cfg config.ServerConfig, log logrus.FieldLogger) *Server {
	listenHost := cfg.ListenHost
	if listenHost == "" {
		listenHost = "127.0.0.1"
	}
	addr := fmt.Sprintf("%s:%d", listenHost, cfg.Port)
	timeout := time.Duration(cfg.Timeout) * time.Second

	return &Server{
		addr: addr,
		log:  log,
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
	s.log.WithFields(logrus.Fields{
		"addr": s.addr,
	}).Info("HTTP-сервер запускается")

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.WithFields(logrus.Fields{
		"addr": s.addr,
	}).Info("HTTP-сервер останавливается")

	return s.httpServer.Shutdown(ctx)
}
