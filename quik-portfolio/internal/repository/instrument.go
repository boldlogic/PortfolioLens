package repository

import (
	"context"
	"database/sql"
	"errors"

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
	INSERT INTO quik.instruments
           (ticker
           ,registration_number
           ,full_name 
           ,short_name
           ,class_code
           ,class_name
           ,isin
           ,face_value
           ,base_currency
           ,quote_currency
           ,counter_currency
           ,maturity_date
           ,coupon_duration
           ,type_id
           ,subtype_id)
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
		,@p10
		,@p11
		,@p12
		,@p13
		,@p14
		,@p15)
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
			return 0, models.ErrInstrumentNotFound
		}
		r.logger.Error("ошибка получения инструмента", zap.String("ticker", ticker), zap.Error(err))
		return 0, err
	}
	return instrumentId, nil
}

func (r *Repository) InsInstrument(ctx context.Context, i models.Instrument) (int, error) {
	var instrumentId int
	r.logger.Debug("сохранение подтипа инструмента", zap.Any("Ticker", i.Ticker))
	row := r.db.QueryRowContext(ctx, insInstrument, i.Ticker, i.RegistrationNumber,
		i.FullName, i.ShortName, i.ClassCode, i.ClassCode, i.ISIN, i.FaceValue,
		i.BaseCurrency, i.QuoteCurrency, i.CounterCurrency, i.MaturityDate,
		i.CouponDuration, i.TypeId, i.SubTypeId)

	err := row.Scan(&instrumentId)
	if err != nil {
		r.logger.Error("ошибка сохранения подтипа инструмента", zap.Any("title", i), zap.Error(err))
		return 0, models.ErrInstrumentCreating
	}

	return instrumentId, nil
}
