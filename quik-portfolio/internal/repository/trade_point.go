package repository

import (
	"context"
	"database/sql"
	"errors"

	md "github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/shutdown"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
	"go.uber.org/zap"
)

const (
	getTradePoints = `
		SELECT
			point_id,
			code,
			name
		FROM quik.trade_points`

	getTradePointByID = `
		SELECT point_id, code, name
		FROM quik.trade_points
		WHERE point_id = @p1`
)

func (r *Repository) GetTradePoints(ctx context.Context) ([]md.TradePoint, error) {
	var result []md.TradePoint

	rows, err := r.Db.QueryContext(ctx, getTradePoints)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return nil, err
		}

		r.Logger.Error("не удалось получить торговые площадки", zap.Error(err))

		return nil, apperrors.ErrRetrievingData
	}
	defer rows.Close()

	for rows.Next() {
		row := md.TradePoint{}
		err = rows.Scan(&row.Id, &row.Code, &row.Name)
		if err != nil {
			if shutdown.IsExceeded(err) {
				return nil, err
			}

			r.Logger.Error("ошибка чтения торговой площадки", zap.Error(err))
			return nil, apperrors.ErrRetrievingData
		}
		result = append(result, row)
	}
	if rows.Err() != nil {
		return nil, apperrors.ErrRetrievingData
	}

	r.Logger.Debug("количество найденных торговых площадок", zap.Int("count", len(result)))

	if len(result) == 0 {
		r.Logger.Warn("торговые площадки не найдены")
		return nil, apperrors.ErrNotFound

	}
	return result, nil
}

func (r *Repository) GetTradePointByID(ctx context.Context, id uint8) (md.TradePoint, error) {
	var row md.TradePoint
	err := r.Db.QueryRowContext(ctx, getTradePointByID, id).Scan(&row.Id, &row.Code, &row.Name)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return md.TradePoint{}, err
		}

		if errors.Is(err, sql.ErrNoRows) {
			return md.TradePoint{}, apperrors.ErrNotFound
		}

		r.Logger.Error("ошибка получения торговой площадки", zap.Uint8("id", id), zap.Error(err))
		return md.TradePoint{}, apperrors.ErrRetrievingData
	}
	r.Logger.Debug("торговая площадка получена", zap.Uint8("id", id))

	return row, nil
}
