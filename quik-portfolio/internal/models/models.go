package models

import (
	"fmt"
	"strings"
	"time"
)

type CurrentQuote struct {
	QuoteDate          *time.Time // Дата торгов
	InstrumentClass    string     // Код инструмента+Борд
	Ticker             string     // Код инструмента
	ISIN               *string    // Международный идентификатор
	RegistrationNumber *string    // Рег.номер инструмента
	FullName           *string    // Полное название инструмента
	ShortName          string     // Краткое название
	FaceValue          *float64   // Номинал
	MaturityDate       *time.Time // Дата погашения
	CouponDuration     *int       // Длительность купона
	ClassCode          string     // Код класса / Борд
	ClassName          string     // Наименование класса
	InstrumentType     string     // Тип инструмента
	InstrumentSubtype  *string    // Подтип инструмента
	Currency           string     // Валюта
	BaseCurrency       string     // Базовая валюта
	QuoteCurrency      *string    // Валюта котировки
	CounterCurrency    *string    // Сопряженная валюта
	InstrumentId       int

	LastPrice     *float64
	ClosePrice    *float64
	AccruedInt    *float64
	TradingStatus *string
}

func (q *CurrentQuote) Clear() {
	q.InstrumentClass = strings.TrimSpace(q.InstrumentClass)
	q.Ticker = strings.TrimSpace(q.Ticker)
	if q.RegistrationNumber != nil {
		trimmedRn := strings.TrimSpace(*q.RegistrationNumber)
		q.RegistrationNumber = &trimmedRn
	}
	if q.FullName != nil {
		trimmedFn := strings.TrimSpace(*q.FullName)
		q.FullName = &trimmedFn
	}
	q.ShortName = strings.TrimSpace(q.ShortName)
	q.ClassCode = strings.TrimSpace(q.ClassCode)
	q.ClassName = strings.TrimSpace(q.ClassName)
	q.InstrumentType = strings.TrimSpace(q.InstrumentType)
	if q.InstrumentSubtype != nil {
		trimmedSt := strings.TrimSpace(*q.InstrumentSubtype)
		q.InstrumentSubtype = &trimmedSt
	}
	if q.ISIN != nil {
		trimmedIsin := strings.TrimSpace(*q.ISIN)
		q.ISIN = &trimmedIsin
	}
	if q.QuoteCurrency != nil {
		trimmedQc := strings.TrimSpace(*q.QuoteCurrency)
		q.QuoteCurrency = &trimmedQc
	}
	if q.CounterCurrency != nil {
		trimmedCc := strings.TrimSpace(*q.CounterCurrency)
		q.CounterCurrency = &trimmedCc
	}
}

func (q CurrentQuote) String() string {
	faceVal := "nil"
	if q.FaceValue != nil {
		faceVal = fmt.Sprintf("%g", *q.FaceValue)
	}
	matDate := "nil"
	if q.MaturityDate != nil {
		matDate = q.MaturityDate.Format(time.DateOnly)
	}
	isin := "nil"
	if q.ISIN != nil {
		isin = fmt.Sprintf("%q", *q.ISIN)
	}
	return fmt.Sprintf("CurrentQuote{Ticker:%q ShortName:%q ClassCode:%q ISIN:%s FaceValue:%s BaseCurrency:%q MaturityDate:%s}",
		q.Ticker, q.ShortName, q.ClassCode, isin, faceVal, q.BaseCurrency, matDate)
}

type InstrumentType struct {
	Id    uint8
	Title string
}

type InstrumentSubType struct {
	SubTypeId uint8
	Title     string
	TypeId    uint8
}

type Firm struct {
	Id   uint8
	Code string
	Name string
}

type Instrument struct {
	Id                 int
	Ticker             string     // Код инструмента
	ISIN               *string    // Международный идентификатор
	RegistrationNumber *string    // Рег.номер инструмента
	FullName           *string    // Полное название инструмента
	ShortName          string     // Краткое название
	MaturityDate       *time.Time // Дата погашения
	CouponDuration     *int       // Длительность купона
	FaceValue          *float64   // Номинал

	// ClassCode string // Код класса
	// ClassName string // Наименование класса
	// TypeId    uint8
	// Type      InstrumentType // Тип инструмента
	// SubTypeId *uint8
	// SubType   InstrumentSubType // Подтип инструмента
	// //InstrumentType     string     `gorm:"column:instrument_type;type:char(100)"`
	// //InstrumentSubtype  string     `gorm:"column:instrument_subtype;type:char(100)"`
	// //AssetClass      string
	// //AssetSubClass   string

	// BaseCurrency    string  // Валюта номинала / базовая валюта
	// QuoteCurrency   *string // Валюта котировки
	// CounterCurrency *string // Сопряженная валюта

}

type MoneyLimit struct {
	LoadDate     time.Time
	ClientCode   string
	Currency     string
	PositionCode string
	FirmCode     string
	FirmName     string
	Balance      float64
}

type SecurityLimit struct {
	LoadDate       time.Time
	ClientCode     string
	Ticker         string
	TradeAccount   string
	SettleCode     string
	FirmCode       string
	FirmName       string
	Balance        float64
	AcquisitionCcy string
	ISIN           *string
}

type Limit struct {
	LoadDate       time.Time
	ClientCode     string
	Ticker         string
	ISIN           *string
	FirmCode       string
	FirmName       string
	Balance        float64
	AcquisitionCcy string
}

// PortfolioItem — позиция с рыночной стоимостью в рублях (по скрипту portfolio: limits + котировки + fx).
type PortfolioItem struct {
	LoadDate       time.Time
	ClientCode     string
	Ticker         string
	TradeAccount   string
	FirmCode       string
	FirmName       string
	Balance        float64
	AcquisitionCcy string
	ISIN           *string
	MvCurrency     *string // валюта рыночной стоимости (из котировки)
	MvRub          float64 // рыночная стоимость в рублях (с учётом НКД)
	ShortName      *string // краткое имя инструмента из котировки
}
