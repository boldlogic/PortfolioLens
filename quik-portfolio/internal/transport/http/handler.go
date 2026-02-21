package httpserver

import (
	v1 "github.com/boldlogic/PortfolioLens/quik-portfolio/internal/transport/http/v1"
	"go.uber.org/zap"
)

type Handler = v1.Handler

// type CommonHandler = v1.CommonHandler
type Service = v1.Service

func NewHandler(svc Service, logger *zap.Logger) *v1.Handler {
	return v1.NewHandler(svc, logger)
}
