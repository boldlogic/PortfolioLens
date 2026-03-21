package taskserver

import (
	"github.com/boldlogic/PortfolioLens/pkg/transport/httpserver/handler"
	v1 "github.com/boldlogic/PortfolioLens/task-manager/internal/transport/http/v1"
	"go.uber.org/zap"
)

type Handler = v1.Handler

type Service = v1.TaskService

func NewHandler(commonHandler handler.Adapter, svc Service, logger *zap.Logger) *v1.Handler {
	return v1.NewHandler(commonHandler, svc, logger)
}
