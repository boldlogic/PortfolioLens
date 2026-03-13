package service

import (
	"context"

	"github.com/boldlogic/PortfolioLens/pkg/models/scheduler"
	"go.uber.org/zap"
)

type SchedulerRepository interface {
	CreateTask(ctx context.Context, actionCode string, taskUUID string) (scheduler.Task, error)
	InsertTaskParams(ctx context.Context, taskId int64, params map[string]string) error
	UpdateTaskStatus(ctx context.Context, id int64, newStatus scheduler.TaskStatusID, errMsg string) error
}

type Service struct {
	logger        *zap.Logger
	schedulerRepo SchedulerRepository
}

func NewService(ctx context.Context,
	schedulerRepo SchedulerRepository,
	logger *zap.Logger) *Service {

	return &Service{

		logger:        logger,
		schedulerRepo: schedulerRepo,
	}
}
