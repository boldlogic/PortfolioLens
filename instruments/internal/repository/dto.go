package repository

import (
	"database/sql"
	"strings"

	"github.com/boldlogic/PortfolioLens/instruments/internal/models"
)

type instrumentBoard struct {
	InstrumentId      int
	BoardId           uint8
	TypeId            uint8
	SubTypeId         sql.NullInt16
	CurrencyId        sql.NullInt64
	BaseCurrencyId    sql.NullInt64
	QuoteCurrencyId   sql.NullInt64
	CounterCurrencyId sql.NullInt64
	IsPrimary         bool
}

type quoteInstrument struct {
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

func (qi *quoteInstrument) convertToInstrument() models.Instrument {
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
