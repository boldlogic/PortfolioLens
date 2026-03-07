package scheduler

import (
	"time"

	UUID "github.com/google/uuid"
)

const (
	DateFormat = "2006-01-02"
)

type TaskStatusID uint8

const (
	TaskStatusScheduled  TaskStatusID = 0
	TaskStatusInProgress TaskStatusID = 1
	TaskStatusCompleted  TaskStatusID = 2
	TaskStatusError      TaskStatusID = 3
)

func (s TaskStatusID) Valid() bool {
	switch s {
	case TaskStatusScheduled, TaskStatusInProgress, TaskStatusCompleted, TaskStatusError:
		return true
	default:
		return false
	}
}

func (s TaskStatusID) String() string {
	switch s {
	case TaskStatusScheduled:
		return "scheduled"
	case TaskStatusInProgress:
		return "in_progress"
	case TaskStatusCompleted:
		return "completed"
	case TaskStatusError:
		return "error"
	default:
		return "unknown"
	}
}

type Task struct {
	Id          int64
	Uuid        UUID.UUID
	ActionId    uint8
	Action      *Action
	CreatedAt   time.Time
	StartedAt   *time.Time
	StatusId    TaskStatusID
	ScheduledAt time.Time
	Error       *string
	CompletedAt *time.Time

	UpdatedAt time.Time
}

type Action struct {
	Id   uint8
	Code string //Пример: currency.cb.fetch.currency_list
	Name string //Пример: Получение справочника валют из ЦБ по www.cbr.ru/scripts/XML_valFull.asp
}

type TaskStatus struct {
	Id   TaskStatusID
	Name string
}
