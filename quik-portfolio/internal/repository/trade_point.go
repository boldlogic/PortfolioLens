package repository

import (
	"context"
	"database/sql"
	"errors"

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

	getTradePointByID = `
		SELECT point_id, code, name
		FROM quik.trade_points
		WHERE point_id = @p1`
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

func (r *Repository) GetTradePointByID(ctx context.Context, id uint8) (models.TradePoint, error) {
	var row models.TradePoint
	err := r.db.QueryRowContext(ctx, getTradePointByID, id).Scan(&row.Id, &row.Code, &row.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.TradePoint{}, apperrors.ErrNotFound
		}
		if IsExceeded(err) {
			return models.TradePoint{}, err
		}
		r.logger.Error("ошибка получения торговой площадки", zap.Uint8("id", id), zap.Error(err))
		return models.TradePoint{}, apperrors.ErrRetrievingData
	}
	return row, nil
}
