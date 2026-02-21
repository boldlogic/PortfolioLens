package httpserver

import (
	"net/http"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/config"
	v1 "github.com/boldlogic/PortfolioLens/quik-portfolio/internal/transport/http/v1"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type Router struct {
	Mux     *chi.Mux
	Handler *Handler
	V1      *v1.Router
	logger  *zap.Logger
	config  *config.Config
}

func NewRouter(handler *Handler, log *zap.Logger, cfg *config.Config) *Router {
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
