package httpserver

import (
	"github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/service"
	v1 "github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/transport/http/v1"
	"github.com/sirupsen/logrus"
)

type Handler = v1.Handler

func NewHandler(logger logrus.FieldLogger, svc service.Service) *v1.Handler {
	return v1.NewHandler(logger, svc)
}
