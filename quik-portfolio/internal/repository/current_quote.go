package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/boldlogic/PortfolioLens/pkg/shutdown"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	"go.uber.org/zap"
)

const (
	//Поклассовый перебор, для дальнейшей оптимизации (важно зафиксировать метрики перед этим)
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

type QuoteInstrument struct {
	InstrumentClass    string          // Код инструмента+Борд
	Ticker             string          // Код инструмента
	ISIN               sql.NullString  // Международный идентификатор
	RegistrationNumber sql.NullString  // Рег.номер инструмента
	FullName           sql.NullString  // Полное название инструмента
	ShortName          string          // Краткое название
	MaturityDate       sql.NullTime    // Дата погашения
	CouponDuration     sql.NullInt64   // Длительность купона
	FaceValue          sql.NullFloat64 // Номинал
	TradePointId       uint8
}

func (r *Repository) SetInstrument(ctx context.Context, id int, ic string) error {
	r.logger.Debug("привязка инструмента к котировке", zap.Int("id", id), zap.String("instrument_class", ic))
	result, err := r.db.ExecContext(ctx, setInstrCurrentQuote, id, ic)
	if err != nil {
		if shutdown.IsExceeded(err) {
			return err
		}

		r.logger.Error("ошибка сохранения инструмента", zap.String("instrument_class", ic), zap.Error(err))
		return apperrors.ErrSavingData
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		r.logger.Warn("котировка не найдена, обновление не выполнено", zap.String("instrument_class", ic))
		return apperrors.ErrNotFound
	}
	r.logger.Debug("инструмент обновлен в котировках", zap.String("instrument_class", ic))
	return nil
}

func (r *Repository) SelectInstrumentFromNewCurrentQuote(ctx context.Context) (models.Instrument, models.InstrumentBoard, string, error) {
	qi := QuoteInstrument{}
	qib := instrumentBoard{}

	r.logger.Debug("выбор котировки без привязанного инструмента")
	row := r.db.QueryRowContext(ctx, selectInstrumentFromNewCurrentQuote)

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
			r.logger.Debug("котировок без инструмента не найдено")
			return models.Instrument{}, models.InstrumentBoard{}, "", apperrors.ErrNotFound
		}
		r.logger.Error("ошибка получения котировки", zap.Error(err))
		return models.Instrument{}, models.InstrumentBoard{}, "", err
	}

	i := qi.convertToInstrument()

	ib := qib.convertToInstrumentBoard()

	return i, ib, qi.InstrumentClass, nil
}

func (qi *QuoteInstrument) convertToInstrument() models.Instrument {
	i := models.Instrument{
		Ticker:    strings.TrimSpace(qi.Ticker),
		ShortName: strings.TrimSpace(qi.ShortName),
	}
	if qi.ISIN.Valid {
		isin := strings.TrimSpace(qi.ISIN.String)
		i.ISIN = &isin
	}
	if qi.RegistrationNumber.Valid {
		registrationNumber := strings.TrimSpace(qi.RegistrationNumber.String)
		i.RegistrationNumber = &registrationNumber
	}
	if qi.FullName.Valid {
		fullName := strings.TrimSpace(qi.FullName.String)
		i.FullName = &fullName
	}

	if qi.MaturityDate.Valid {
		maturityDate := qi.MaturityDate.Time
		i.MaturityDate = &maturityDate
	}

	if qi.CouponDuration.Valid {
		couponDuration := int(qi.CouponDuration.Int64)
		i.CouponDuration = &couponDuration
	}
	if qi.FaceValue.Valid {
		faceValue := qi.FaceValue.Float64
		i.FaceValue = &faceValue
	}
	i.TradePointId = qi.TradePointId
	return i
}

func (qib *instrumentBoard) convertToInstrumentBoard() models.InstrumentBoard {

	ib := models.InstrumentBoard{}
	ib.BoardId = qib.BoardId
	ib.TypeId = qib.TypeId

	if qib.SubTypeId.Valid {
		sid := uint8(qib.SubTypeId.Int16)
		ib.SubTypeId = &sid
	}
	if qib.CurrencyId.Valid {
		cid := int(qib.CurrencyId.Int64)
		ib.CurrencyId = &cid
	}
	if qib.BaseCurrencyId.Valid {
		cid := int(qib.BaseCurrencyId.Int64)
		ib.BaseCurrencyId = &cid
	}
	if qib.QuoteCurrencyId.Valid {
		cid := int(qib.QuoteCurrencyId.Int64)
		ib.QuoteCurrencyId = &cid
	}
	if qib.CounterCurrencyId.Valid {
		cid := int(qib.CounterCurrencyId.Int64)
		ib.CounterCurrencyId = &cid
	}
	return ib
}
