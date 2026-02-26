package repository

import (
	"context"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	"go.uber.org/zap"
)

const (
	getTradePoints = `
		SELECT
			point_id,
			code,
			name
		FROM quik.trade_points`
)

func (r *Repository) GetTradePoints(ctx context.Context) ([]models.TradePoint, error) {
	var result []models.TradePoint

	rows, err := r.db.QueryContext(ctx, getTradePoints)
	if err != nil {
		if IsExceeded(err) {
			return nil, err
		}

		r.logger.Error("не удалось получить торговые площадки", zap.Error(err))

		return nil, apperrors.ErrRetrievingData
	}
	defer rows.Close()

	for rows.Next() {
		row := models.TradePoint{}
		err = rows.Scan(&row.Id, &row.Code, &row.Name)
		if err != nil {
			if IsExceeded(err) {
				return nil, err
			}
			r.logger.Error("не удалось получить торговые площадки", zap.Error(err))

			return nil, apperrors.ErrRetrievingData
		}
		result = append(result, row)
	}
	if rows.Err() != nil {
		return nil, apperrors.ErrRetrievingData
	}
	r.logger.Debug("количество найденных торговых площадок", zap.Int("count", len(result)))

	if len(result) == 0 {
		r.logger.Warn("торговые площадки не найдены")
		return nil, apperrors.ErrNotFound

	}
	return result, nil
}
