package v1

import (
	"github.com/boldlogic/cbr-market-data-worker/internal/service"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Service service.Client
	log     logrus.FieldLogger
}

func NewHandler(logger logrus.FieldLogger, svc service.Client) *Handler {
	return &Handler{
		log:     logger,
		Service: svc,
	}
}
