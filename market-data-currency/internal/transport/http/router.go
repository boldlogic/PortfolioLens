package currencyserver

import (
	v1 "github.com/boldlogic/PortfolioLens/market-data-currency/internal/transport/http/v1"
	"github.com/boldlogic/PortfolioLens/pkg/metrics"
	"github.com/boldlogic/PortfolioLens/pkg/transport/httpserver/router"
	"go.uber.org/zap"
)

type Router struct {
	*router.Router
}

func NewRouter(handler *v1.Handler, log *zap.Logger, reg metrics.Registry) *Router {
	base := router.NewRouter(log, reg)
	base.Mount("/api/v1", v1.NewRouter(handler, log))
	return &Router{Router: base}
}
