package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/models/scheduler"
	"github.com/boldlogic/PortfolioLens/pkg/shutdown"
	mssql "github.com/microsoft/go-mssqldb"
	"go.uber.org/zap"
)

// TO_DO в прекрасном будущем: func (r *Repository) InsertTasks(ctx context.Context, tasks []scheduler.Task) error {}

const (
	insertTask = `
		INSERT INTO dbo.tasks (uuid, action_id)
		OUTPUT inserted.id, inserted.uuid, inserted.created_at, inserted.scheduled_at
		SELECT @p1, a.id
		FROM dbo.actions a
		WHERE a.code = @p2`
	updateTaskStatus = `
		UPDATE t
		SET
			t.updated_at  = SYSDATETIMEOFFSET(),
			t.completed_at  = SYSDATETIMEOFFSET(),
			t.status_id   = @p1,  
			t.error=@p2
		FROM dbo.tasks t
		WHERE t.id=@p3`
)

func (r *Repository) CreateTask(ctx context.Context, actionCode string, taskUUID string) (scheduler.Task, error) {

	var out scheduler.Task
	row := r.Db.QueryRowContext(ctx, insertTask, taskUUID, actionCode)
	err := row.Scan(&out.Id, &out.Uuid, &out.CreatedAt, &out.ScheduledAt)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return scheduler.Task{}, err
		}
		if errors.Is(err, sql.ErrNoRows) {
			return scheduler.Task{}, models.ErrBusinessValidation
		}

		r.Logger.Error("ошибка создания задачи", zap.String("action.code", actionCode), zap.String("uuid", taskUUID), zap.Error(err))

		var mssqlErr mssql.Error
		if errors.As(err, &mssqlErr) && (mssqlErr.Number == 2627 || mssqlErr.Number == 2601) {
			return scheduler.Task{}, models.ErrConflict
		}
		return scheduler.Task{}, models.ErrSavingData
	}

	return scheduler.Task{
		Id:          out.Id,
		Uuid:        out.Uuid,
		CreatedAt:   out.CreatedAt,
		ScheduledAt: out.ScheduledAt,
	}, nil
}

func (r *Repository) UpdateTaskStatus(ctx context.Context, id int64, newStatus scheduler.TaskStatusID, errMsg string) error {

	_, err := r.Db.ExecContext(ctx, updateTaskStatus, newStatus, errMsg, id)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		r.Logger.Error("ошибка при обновлении статуса задачи", zap.Error(err))

		return models.ErrSavingData
	}

	return nil
}
