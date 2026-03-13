package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/models/scheduler"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	fetchOneNewTask = `
		WITH next_task AS (
			SELECT TOP (1) t.id
			FROM dbo.tasks t WITH (UPDLOCK, ROWLOCK)
			JOIN dbo.task_statuses st ON st.id = t.status_id AND st.name = N'scheduled'
			WHERE t.scheduled_at<=SYSDATETIMEOFFSET()
			ORDER BY t.scheduled_at ASC
		)
		UPDATE t
		SET
			t.started_at  = SYSDATETIMEOFFSET(),
			t.status_id   = (SELECT id FROM dbo.task_statuses WHERE name = N'in_progress'),
			t.updated_at  = SYSDATETIMEOFFSET()
		OUTPUT inserted.*
		FROM dbo.tasks t
		JOIN next_task n ON n.id = t.id;	
		`
	updateTaskStatus = `
		UPDATE t
		SET
			t.updated_at   = SYSDATETIMEOFFSET(),
			t.completed_at = CASE
				WHEN @p1 = (SELECT id FROM dbo.task_statuses WHERE name = N'completed')
				THEN SYSDATETIMEOFFSET()
				ELSE NULL
			END,
			t.status_id = @p1,
			t.error     = @p2
		FROM dbo.tasks t
		WHERE t.id = @p3`
)

type rawTask struct {
	Id          int64
	Uuid        uuid.UUID
	ActionId    uint8
	StatusId    uint8
	CreatedAt   time.Time
	StartedAt   sql.NullTime
	ScheduledAt time.Time
	Error       sql.NullString
	CompletedAt sql.NullTime
	UpdatedAt   time.Time
}

func (r *Repository) UpdateTaskStatus(ctx context.Context, id int64, newStatus scheduler.TaskStatusID, errMsg string) error {

	_, err := r.Db.ExecContext(ctx, updateTaskStatus, newStatus, errMsg, id)
	if err != nil {
		if r.isShutdown(err) {
			return err
		}
		r.Logger.Error("ошибка при обновлении статуса задачи", zap.Error(err))

		return models.ErrSavingData
	}

	return nil
}

func (r *Repository) FetchOneNewTask(ctx context.Context) (scheduler.Task, error) {
	var raw rawTask
	row := r.Db.QueryRowContext(ctx, fetchOneNewTask)
	err := row.Scan(&raw.Id,
		&raw.Uuid,
		&raw.ActionId,
		&raw.StatusId,
		&raw.CreatedAt,
		&raw.StartedAt,
		&raw.ScheduledAt,
		&raw.CompletedAt,
		&raw.UpdatedAt,
		&raw.Error)
	if err != nil {
		if r.isShutdown(err) {
			return scheduler.Task{}, err
		}
		if errors.Is(err, sql.ErrNoRows) {
			r.Logger.Debug("новых задач не найдено")
			return scheduler.Task{}, models.ErrNotFound
		}
		r.Logger.Error("ошибка при получении задачи", zap.Error(err))

		return scheduler.Task{}, models.ErrRetrievingData
	}

	return rawToTask(raw), nil
}

func rawToTask(raw rawTask) scheduler.Task {
	var out scheduler.Task
	out.Id = raw.Id
	out.Uuid = raw.Uuid
	out.ActionId = raw.ActionId
	out.CreatedAt = raw.CreatedAt
	out.ScheduledAt = raw.ScheduledAt
	out.UpdatedAt = raw.UpdatedAt
	if raw.StartedAt.Valid {
		out.StartedAt = &raw.StartedAt.Time
	}
	if raw.CompletedAt.Valid {
		out.CompletedAt = &raw.CompletedAt.Time
	}
	if raw.Error.Valid {
		out.Error = &raw.Error.String
	}
	out.StatusId = scheduler.TaskStatusID(raw.StatusId)
	return out
}
