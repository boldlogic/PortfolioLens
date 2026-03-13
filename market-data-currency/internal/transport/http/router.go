package currencyserver

import (
	v1 "github.com/boldlogic/PortfolioLens/market-data-currency/internal/transport/http/v1"
	"github.com/boldlogic/PortfolioLens/pkg/transport/httpserver"
	"github.com/boldlogic/PortfolioLens/pkg/transport/httpserver/router"
	"go.uber.org/zap"
)

type Router struct {
	CommonRouter *router.Router
	V1           *v1.Router
	logger       *zap.Logger
	config       *httpserver.ServerConfig
}

func NewRouter(handler *v1.Handler, log *zap.Logger, cfg *httpserver.ServerConfig) *Router {
	commonRouter := router.NewRouter(handler)

	v1Router := v1.NewRouter(handler, log)
	commonRouter.Mux.Mount("/api/v1", v1Router.Mux)

	return &Router{
		CommonRouter: commonRouter,
		V1:           v1Router,
		logger:       log,
		config:       cfg,
	}
}
