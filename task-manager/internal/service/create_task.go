package service

import (
	"context"
	"fmt"

	"github.com/boldlogic/PortfolioLens/pkg/models/scheduler"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) CreateTask(ctx context.Context, actionCode string, taskUUID string, params map[string]string) (scheduler.Task, error) {

	if taskUUID == "" {
		u, err := uuid.NewV7()
		if err != nil {
			return scheduler.Task{}, err
		}
		taskUUID = u.String()
	}

	task, err := s.schedulerRepo.CreateTask(ctx, actionCode, taskUUID)
	if err != nil {
		s.logger.Warn("ошибка создания задачи",
			zap.String("action", actionCode),
			zap.Error(err))
		return scheduler.Task{}, err
	}

	if len(params) > 0 {
		if err = s.schedulerRepo.InsertTaskParams(ctx, task.Id, params); err != nil {
			s.logger.Error("ошибка сохранения параметров задачи",
				zap.Int64("task_id", task.Id),
				zap.Error(err))
			_ = s.schedulerRepo.UpdateTaskStatus(ctx, task.Id, scheduler.TaskStatusError,
				fmt.Sprintf("ошибка сохранения параметров: %v", err))
			return scheduler.Task{}, err
		}
	}

	s.logger.Info("задача создана",
		zap.Int64("task_id", task.Id),
		zap.String("uuid", task.Uuid.String()),
		zap.String("action", actionCode))

	return task, nil
}
