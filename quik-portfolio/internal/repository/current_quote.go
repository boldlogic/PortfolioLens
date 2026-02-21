package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	"go.uber.org/zap"
)

const (
	selectNewCurrentQuote = `
	SELECT TOP (1) 
	  instrument_class
	  ,ticker
      ,registration_number
      ,full_name
      ,short_name
      ,class_code
      ,class_name
      ,instrument_type
      ,instrument_subtype
      ,isin
      ,face_value
      ,base_currency
      ,quote_currency
      ,counter_currency
      ,maturity_date
      ,coupon_duration
  	FROM quik.current_quotes
  	where instrument_id is null
	`
	setInstrCurrentQuote = `
	update quik.current_quotes
	set instrument_id=@p1
	where instrument_class=@p2
	`
)

func (r *Repository) SetInstrument(ctx context.Context, id int, ic string) error {
	r.logger.Debug("инструмента", zap.Any("instrument_class", ic))
	_, err := r.db.ExecContext(ctx, setInstrCurrentQuote, id, ic)

	if err != nil {
		r.logger.Error("ошибка сохранения  инструмента", zap.Any("instrument_class", ic), zap.Error(err))
		return models.ErrInstrumentCreating
	}
	r.logger.Debug("инструмент обновлен в котировках", zap.Any("instrument_class", ic))

	return nil
}

func (r *Repository) SelectNewCurrentQuote(ctx context.Context) (models.CurrentQuote, error) {
	//to-do добавить подсчет времени
	qt := models.CurrentQuote{}
	r.logger.Debug("сохранение типа инструмента")
	row := r.db.QueryRowContext(ctx, selectNewCurrentQuote)

	err := row.Scan(&qt.InstrumentClass, &qt.Ticker,
		&qt.RegistrationNumber,
		&qt.FullName,
		&qt.ShortName,
		&qt.ClassCode,
		&qt.ClassName,
		&qt.InstrumentType,
		&qt.InstrumentSubtype,
		&qt.ISIN,
		&qt.FaceValue,
		&qt.BaseCurrency,
		&qt.QuoteCurrency,
		&qt.CounterCurrency,
		&qt.MaturityDate,
		&qt.CouponDuration)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Error("типа инструмента не найден", zap.String("title", "title"))
			return models.CurrentQuote{}, models.ErrInstrumentTypeNotFound
		}

		r.logger.Error("ошибка получения типа инструмента", zap.String("title", "title"), zap.Error(err))
		return models.CurrentQuote{}, err
	}

	qt.Ticker = strings.TrimSpace(qt.Ticker)

	return qt, nil
}

// type CurrentQuote struct {
// 	Ticker             string     // Код инструмента
// 	RegistrationNumber *string    // Рег.номер инструмента
// 	FullName           *string    // Полное название инструмента
// 	ShortName          string     // Краткое название
// 	ClassCode          string     // Код класса
// 	ClassName          string     // Наименование класса
// 	InstrumentType     string     // Тип инструмента
// 	InstrumentSubtype  *string    // Подтип инструмента
// 	ISIN               *string    // Международный идентификатор
// 	FaceValue          *float64   // Номинал
// 	BaseCurrency       string     // Валюта номинала / базовая валюта
// 	QuoteCurrency      *string    // Валюта котировки
// 	CounterCurrency    *string    // Сопряженная валюта
// 	MaturityDate       *time.Time // Дата погашения
// 	CouponDuration     *int       // Длительность купона
// }
