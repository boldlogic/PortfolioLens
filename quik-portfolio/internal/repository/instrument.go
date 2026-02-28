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
	selInstrumentId = `
		SELECT instrument_id
		FROM quik.instruments
		where ticker=@p1
	`
	insInstrument = `
		INSERT INTO quik.instruments (
			ticker,
			registration_number,
			full_name,
			short_name,
			isin,
			face_value,
			maturity_date,
			coupon_duration
			)
		output inserted.instrument_id
		VALUES (
			@p1
			,@p2
			,@p3
			,@p4
			,@p5
			,@p6
			,@p7
			,@p8
		)	
	`
)

func (r *Repository) GetInstrumentId(ctx context.Context, ticker string) (int, error) {
	var instrumentId int
	r.logger.Debug("получение id инструмента по тикеру", zap.String("ticker", ticker))
	row := r.db.QueryRowContext(ctx, selInstrumentId, ticker)
	err := row.Scan(&instrumentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("инструмент не найден", zap.String("ticker", ticker))
			return 0, apperrors.ErrNotFound
		}
		r.logger.Error("ошибка получения инструмента", zap.String("ticker", ticker), zap.Error(err))
		return 0, err
	}
	return instrumentId, nil
}

func (r *Repository) InsInstrument(ctx context.Context, i models.Instrument) (int, error) {
	var instrumentId int
	r.logger.Debug("сохранение инструмента", zap.String("Ticker", i.Ticker))
	row := r.db.QueryRowContext(ctx, insInstrument, i.Ticker, i.RegistrationNumber,
		i.FullName, i.ShortName, i.ISIN, i.FaceValue, i.MaturityDate,
		i.CouponDuration)

	err := row.Scan(&instrumentId)
	if err != nil {
		r.logger.Error("ошибка сохранения инструмента", zap.String("Ticker", i.Ticker), zap.Error(err))
		return 0, apperrors.ErrSavingData
	}

	return instrumentId, nil
}
