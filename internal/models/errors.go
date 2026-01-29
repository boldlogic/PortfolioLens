package models

import "errors"

var (
	ErrTaskRetrievingError = errors.New("не удалось получить task")
	ErrActionNotFound      = errors.New("action не найден")
	ErrTaskNotFound        = errors.New("задача не найдена")
	ErrTaskAlreadyExists   = errors.New("задача уже существует")
	ErrTaskCreating        = errors.New("не удалось создать задачу")
	ErrTaskUpdating        = errors.New("не удалось обновить задачу")
	ErrNoNewTasks          = errors.New("не найдены новые задачи")
)
