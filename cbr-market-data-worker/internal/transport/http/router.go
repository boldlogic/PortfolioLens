package httpserver

import (
	"net/http"

	"github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/config"
	v1 "github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/transport/http/v1"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

type Router struct {
	Mux     *chi.Mux
	Handler *Handler
	V1      *v1.Router
	logger  logrus.FieldLogger
	config  *config.Config
}

func NewRouter(handler *Handler, log logrus.FieldLogger, cfg *config.Config) *Router {
	r := chi.NewRouter()
	r.Get("/healthcheck", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	v1Router := v1.NewRouter(handler, log)
	r.Mount("/api/v1", v1Router.Mux)

	return &Router{
		Mux:    r,
		V1:     v1Router,
		logger: log,
		config: cfg,
	}
}
