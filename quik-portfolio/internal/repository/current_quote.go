package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	"go.uber.org/zap"
)

const (
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
			q.coupon_duration
		FROM
			quik.current_quotes q
			JOIN quik.boards b ON b.code = q.class_code
			JOIN quik.trade_points p ON p.point_id = b.trade_point_id
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
}

func (r *Repository) SetInstrument(ctx context.Context, id int, ic string) error {
	r.logger.Debug("привязка инструмента к котировке", zap.Int("id", id), zap.String("instrument_class", ic))
	result, err := r.db.ExecContext(ctx, setInstrCurrentQuote, id, ic)
	if err != nil {
		if IsExceeded(err) {
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

func (r *Repository) SelectInstrumentFromNewCurrentQuote(ctx context.Context) (models.Instrument, string, error) {
	qi := QuoteInstrument{}
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
		&qi.CouponDuration)

	if err != nil {
		if IsExceeded(err) {
			return models.Instrument{}, "", err
		}
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("котировок без инструмента не найдено")
			return models.Instrument{}, "", apperrors.ErrNotFound
		}
		r.logger.Error("ошибка получения котировки", zap.Error(err))
		return models.Instrument{}, "", err
	}

	i := qi.convertToInstrument()

	return i, qi.InstrumentClass, nil
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
	return i
}
