package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/boldlogic/PortfolioLens/instruments/internal/models"
	md "github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/shutdown"
	"go.uber.org/zap"
)

const (
	selectInstrumentId = `
		SELECT
			instrument_id
		FROM
			quik.instruments
		WHERE
			ticker = @p1
			AND trade_point_id = @p2
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
			coupon_duration,
			trade_point_id
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
			,@p9
		)	
	`
)

func (r *Repository) GetInstrumentId(ctx context.Context, ticker string, tradePointId uint8) (int, error) {
	var instrumentId int
	row := r.Db.QueryRowContext(ctx, selectInstrumentId, ticker, tradePointId)
	err := row.Scan(&instrumentId)

	if err != nil {
		if shutdown.IsExceeded(err) {
			return 0, err
		}

		if errors.Is(err, sql.ErrNoRows) {
			r.Logger.Debug("инструмент не найден", zap.String("ticker", ticker), zap.Uint8("trade_point_id", tradePointId))
			return 0, md.ErrNotFound
		}

		r.Logger.Error("ошибка получения инструмента", zap.String("ticker", ticker), zap.Uint8("trade_point_id", tradePointId), zap.Error(err))
		return 0, err
	}

	r.Logger.Debug("получение id инструмента по тикеру", zap.String("ticker", ticker), zap.Uint8("trade_point_id", tradePointId), zap.Int("instrumentId", instrumentId))

	return instrumentId, nil
}

func (r *Repository) InsInstrument(ctx context.Context, i models.Instrument) (int, error) {
	var instrumentId int
	r.Logger.Debug("сохранение инструмента", zap.String("Ticker", i.Ticker))
	row := r.Db.QueryRowContext(ctx, insInstrument, i.Ticker, i.RegistrationNumber,
		i.FullName, i.ShortName, i.ISIN, i.FaceValue, i.MaturityDate,
		i.CouponDuration, i.TradePointId)

	err := row.Scan(&instrumentId)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return 0, err
		}
		r.Logger.Error("ошибка сохранения инструмента", zap.String("ticker", i.Ticker), zap.Uint8("trade_point_id", i.TradePointId), zap.Error(err))
		return 0, md.ErrSavingData
	}

	return instrumentId, nil
}
