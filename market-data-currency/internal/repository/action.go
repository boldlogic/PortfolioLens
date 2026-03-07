package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/boldlogic/PortfolioLens/market-data-currency/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/pkg/models/scheduler"
	"github.com/boldlogic/PortfolioLens/pkg/shutdown"
	"go.uber.org/zap"
)

const (
	selectActionByCode = `
		SELECT id
			,code
			,name
		FROM dbo.actions
		where id=@p1`
)

func (r *Repository) SelectAction(ctx context.Context, id uint8) (scheduler.Action, error) {
	var a scheduler.Action

	row := r.Db.QueryRowContext(ctx, selectActionByCode, id)
	err := row.Scan(&a.Id, &a.Code, &a.Name)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return scheduler.Action{}, err
		}
		if errors.Is(err, sql.ErrNoRows) {
			r.Logger.Debug("не найдено действие")
			return scheduler.Action{}, apperrors.ErrNotFound
		}
		r.Logger.Error("ошибка при получении действия", zap.Error(err))

		return scheduler.Action{}, apperrors.ErrRetrievingData
	}
	return a, nil

}
