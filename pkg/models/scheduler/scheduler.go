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
	UpdatedAt   time.Time
}

const (
	ActionCodeCbrCurrencyList    = "currency.cb.fetch.currency_list"
	ActionCodeCbrRatesToday      = "currency.cb.fetch.rates_today"
	ActionCodeCbrHistoricalRates = "currency.cb.fetch.historical_rates"
)

type Action struct {
	Id   uint8
	Code string
	Name string
}

type TaskStatus struct {
	Id   TaskStatusID
	Name string
}

type TaskParam struct {
	TaskId  int64
	ParamId int
	Code    string
	Value   string
}
