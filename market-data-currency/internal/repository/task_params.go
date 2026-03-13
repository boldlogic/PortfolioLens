package repository

import (
	"context"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/models/scheduler"

	"go.uber.org/zap"
)

const selectTaskParams = `
	SELECT
		tp.task_id,
		tp.param_id,
		p.code,
		tp.value
	FROM dbo.task_params tp
	JOIN dbo.params p ON p.id = tp.param_id
	WHERE tp.task_id = @p1`

func (r *Repository) SelectTaskParams(ctx context.Context, taskId int64) ([]scheduler.TaskParam, error) {
	rows, err := r.Db.QueryContext(ctx, selectTaskParams, taskId)
	if err != nil {
		if r.isShutdown(err) {
			return nil, err
		}
		r.Logger.Error("ошибка при получении параметров задачи", zap.Int64("task_id", taskId), zap.Error(err))
		return nil, models.ErrRetrievingData
	}
	defer rows.Close()

	var result []scheduler.TaskParam
	for rows.Next() {
		var p scheduler.TaskParam
		if err = rows.Scan(&p.TaskId, &p.ParamId, &p.Code, &p.Value); err != nil {
			if r.isShutdown(err) {
				return nil, err
			}
			r.Logger.Error("ошибка при чтении параметра задачи", zap.Int64("task_id", taskId), zap.Error(err))
			return nil, models.ErrRetrievingData
		}
		result = append(result, p)
	}
	if rows.Err() != nil {
		r.Logger.Error("ошибка при итерации параметров задачи", zap.Error(rows.Err()))
		return nil, models.ErrRetrievingData
	}
	return result, nil
}
