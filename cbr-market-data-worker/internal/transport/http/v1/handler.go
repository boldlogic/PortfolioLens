package v1

import (
	"github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/service"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Service service.Service
	log     logrus.FieldLogger
}

func NewHandler(logger logrus.FieldLogger, svc service.Service) *Handler {
	return &Handler{
		log:     logger,
		Service: svc,
	}
}
