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
	// Поклассовый перебор, для дальнейшей оптимизации (важно зафиксировать метрики перед этим)
	selectInstrumentFromNewCurrentQuote = `
		SELECT TOP (1)
			q.instrument_class,
			q.ticker,
			q.isin,
			q.registration_number,
			q.full_name,
			q.short_name,
			q.face_value,
			q.maturity_date,
			q.coupon_duration,
			p.point_id,
			b.board_id,
			it.type_id,
			ist.subtype_id,
			curr.iso_code,
			base_curr.iso_code,
			quote_curr.iso_code,
			counter_curr.iso_code
		FROM
			quik.current_quotes q
			JOIN quik.boards b ON b.code = q.class_code
			JOIN quik.trade_points p ON p.point_id = b.trade_point_id
			JOIN quik.instrument_types it ON it.title = q.instrument_type
    		LEFT JOIN quik.instrument_subtypes ist ON ist.type_id = it.type_id AND ist.title = q.instrument_subtype
			LEFT JOIN dbo.currencies curr ON curr.iso_char_code = RTRIM(CASE WHEN q.currency IN ('SUR', 'RUR') THEN 'RUB' ELSE q.currency END)
    		LEFT JOIN dbo.currencies base_curr ON base_curr.iso_char_code = RTRIM(CASE WHEN q.base_currency IN ('SUR', 'RUR') THEN 'RUB' ELSE q.base_currency END)
    		LEFT JOIN dbo.currencies quote_curr ON quote_curr.iso_char_code = RTRIM(CASE WHEN q.quote_currency IN ('SUR', 'RUR') THEN 'RUB' ELSE q.quote_currency END)
    		LEFT JOIN dbo.currencies counter_curr ON counter_curr.iso_char_code = RTRIM(CASE WHEN q.counter_currency IN ('SUR', 'RUR') THEN 'RUB' ELSE q.counter_currency END)
		WHERE
			q.instrument_id IS NULL
			`
	setInstrCurrentQuote = `
			update quik.current_quotes
			set instrument_id=@p1
			where instrument_class=@p2
	`
)

func (r *Repository) SetInstrument(ctx context.Context, id int, ic string) error {
	result, err := r.Db.ExecContext(ctx, setInstrCurrentQuote, id, ic)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}
		r.Logger.Error("ошибка сохранения инструмента", zap.String("instrument_class", ic), zap.Error(err))
		return md.ErrSavingData
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		r.Logger.Warn("котировка не найдена, обновление не выполнено", zap.String("instrument_class", ic))
		return md.ErrNotFound
	}
	r.Logger.Debug("инструмент обновлен в котировках", zap.String("instrument_class", ic), zap.Int("instrument_id", id))
	return nil
}

func (r *Repository) SelectInstrumentFromNewCurrentQuote(ctx context.Context) (models.Instrument, models.InstrumentBoard, string, error) {
	qi := quoteInstrument{}
	qib := instrumentBoard{}

	row := r.Db.QueryRowContext(ctx, selectInstrumentFromNewCurrentQuote)

	err := row.Scan(&qi.InstrumentClass,
		&qi.Ticker,
		&qi.ISIN,
		&qi.RegistrationNumber,
		&qi.FullName,
		&qi.ShortName,
		&qi.FaceValue,
		&qi.MaturityDate,
		&qi.CouponDuration,
		&qi.TradePointId,
		&qib.BoardId,
		&qib.TypeId,
		&qib.SubTypeId,
		&qib.CurrencyId,
		&qib.BaseCurrencyId,
		&qib.QuoteCurrencyId,
		&qib.CounterCurrencyId)

	if err != nil {
		if shutdown.IsExceeded(err) {
			return models.Instrument{}, models.InstrumentBoard{}, "", err
		}
		if errors.Is(err, sql.ErrNoRows) {
			r.Logger.Debug("котировок без инструмента не найдено")
			return models.Instrument{}, models.InstrumentBoard{}, "", md.ErrNotFound
		}
		r.Logger.Error("ошибка получения котировки", zap.Error(err))
		return models.Instrument{}, models.InstrumentBoard{}, "", err
	}

	return qi.convertToInstrument(), qib.convertToInstrumentBoard(), qi.InstrumentClass, nil
}
