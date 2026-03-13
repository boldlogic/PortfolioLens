package v1

import (
	"context"
	"net/http"

	"github.com/boldlogic/PortfolioLens/pkg/models/scheduler"
	"github.com/boldlogic/PortfolioLens/pkg/transport/httpserver/handler"
	"go.uber.org/zap"
)

type TaskService interface {
	CreateTask(ctx context.Context, actionCode string, taskUUID string, params map[string]string) (scheduler.Task, error)
}

type CommonHandler interface {
	Adapt(fn handler.HandlerFunc) http.HandlerFunc
}

type Handler struct {
	commonHandler CommonHandler
	taskSvc       TaskService
	logger        *zap.Logger
}

func NewHandler(commonHandler CommonHandler, taskSvc TaskService, log *zap.Logger) *Handler {
	return &Handler{
		logger:        log,
		taskSvc:       taskSvc,
		commonHandler: commonHandler,
	}
}

func (h *Handler) Adapt(fn handler.HandlerFunc) http.HandlerFunc {
	return h.commonHandler.Adapt(fn)
}
