package taskserver

import (
	v1 "github.com/boldlogic/PortfolioLens/task-manager/internal/transport/http/v1"
	"go.uber.org/zap"
)

type Handler = v1.Handler

type CommonHandler = v1.CommonHandler

type Service = v1.TaskService

func NewHandler(commonHandler CommonHandler, svc Service, logger *zap.Logger) *v1.Handler {
	return v1.NewHandler(commonHandler, svc, logger)
}
