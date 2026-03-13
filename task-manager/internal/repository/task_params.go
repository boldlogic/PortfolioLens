package repository

import (
	"context"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/shutdown"

	"go.uber.org/zap"
)

const insertTaskParam = `
	INSERT INTO dbo.task_params (task_id, param_id, value)
	SELECT @p1, p.id, @p3
	FROM dbo.params p
	WHERE p.code = @p2`

func (r *Repository) InsertTaskParam(ctx context.Context, taskId int64, paramCode string, value string) error {
	res, err := r.Db.ExecContext(ctx, insertTaskParam, taskId, paramCode, value)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		r.Logger.Error("ошибка при сохранении параметра задачи",
			zap.Int64("task_id", taskId),
			zap.String("param_code", paramCode),
			zap.Error(err))
		return models.ErrSavingData
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		r.Logger.Warn("параметр задачи не сохранён: param code не найден в dbo.params",
			zap.Int64("task_id", taskId),
			zap.String("param_code", paramCode))
		return models.ErrBusinessValidation
	}
	return nil
}

func (r *Repository) InsertTaskParams(ctx context.Context, taskId int64, params map[string]string) error {
	for code, value := range params {
		if err := r.InsertTaskParam(ctx, taskId, code, value); err != nil {
			return err
		}
	}
	return nil
}
