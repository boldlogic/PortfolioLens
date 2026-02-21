package service

import (
	"context"
	"errors"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	UUID "github.com/google/uuid"
	"gorm.io/gorm"
)

func (c *Service) CreateTask(ctx context.Context, task models.Task, action string) (models.Task, error) {
	act, err := c.schedulerRepo.GetActionId(action)
	if err != nil && err != gorm.ErrRecordNotFound {
		return models.Task{}, err
	}
	if act.Id == 0 || err == gorm.ErrRecordNotFound {
		c.log.Warnf("Не найден тип события с кодом %s", action)
		return models.Task{}, models.ErrActionNotFound
	}
	task.ActionId = act.Id

	if task.Uuid != "" {
		exists, err := c.schedulerRepo.GetTask(task.Uuid)
		if err != nil && !errors.Is(err, models.ErrTaskNotFound) {
			c.log.Warnf("Не удалось получить задачу по uuid %s", task.Uuid)
			return models.Task{}, models.ErrTaskRetrievingError
		}
		if exists.Id != 0 {
			c.log.Infof("Задача с uuid %s уже существует, id=%d", task.Uuid, exists.Id)
			return models.Task{}, models.ErrTaskAlreadyExists

		}
	} else {
		task.Uuid = UUID.New().String()
	}

	c.log.Infof("Создаем задачу по uuid %s", task.Uuid)
	created, err := c.schedulerRepo.CreateTask(&task)
	if err != nil {
		c.log.Warnf("Не удалось создать задачу по uuid %s", task.Uuid)
		return models.Task{}, models.ErrTaskCreating

	}

	return created, nil
}
