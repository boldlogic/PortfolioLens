package models

import (
	"fmt"
	"strings"
	"time"
)

type CurrentQuote struct {
	InstrumentClass    string
	Ticker             string     // Код инструмента
	RegistrationNumber *string    // Рег.номер инструмента
	FullName           *string    // Полное название инструмента
	ShortName          string     // Краткое название
	ClassCode          string     // Код класса
	ClassName          string     // Наименование класса
	InstrumentType     string     // Тип инструмента
	InstrumentSubtype  *string    // Подтип инструмента
	ISIN               *string    // Международный идентификатор
	FaceValue          *float64   // Номинал
	BaseCurrency       string     // Валюта номинала / базовая валюта
	QuoteCurrency      *string    // Валюта котировки
	CounterCurrency    *string    // Сопряженная валюта
	MaturityDate       *time.Time // Дата погашения
	CouponDuration     *int       // Длительность купона
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
	Id    int16
	Title string
}

type InstrumentSubType struct {
	SubTypeId int16
	TypeId    int16
	Title     string
}

type Board struct {
	Id   int16
	Code string
	Name string
}

type Firm struct {
	Id   int8
	Code string
	Name string
}

type Instrument struct {
	Id                 int
	Ticker             string  `gorm:"column:ticker;type:char(15);not null"`      // Код инструмента
	RegistrationNumber *string `gorm:"column:registration_number;type:char(250)"` // Рег.номер инструмента
	FullName           *string `gorm:"column:full_name;type:char(250);not null"`  // Полное название инструмента
	ShortName          string  `gorm:"column:short_name;type:char(100)"`          // Краткое название
	ClassCode          string  `gorm:"column:class_code;type:char(20)"`           // Код класса
	ClassName          string  `gorm:"column:class_name;type:char(200)"`          // Наименование класса
	TypeId             int16
	Type               InstrumentType // Тип инструмента
	SubTypeId          *int16
	SubType            InstrumentSubType // Подтип инструмента
	//InstrumentType     string     `gorm:"column:instrument_type;type:char(100)"`
	//InstrumentSubtype  string     `gorm:"column:instrument_subtype;type:char(100)"`
	//AssetClass      string
	//AssetSubClass   string
	ISIN            *string    `gorm:"column:isin;type:char(15)"`            // Международный идентификатор
	FaceValue       *float64   `gorm:"column:face_value;type:float"`         // Номинал
	BaseCurrency    string     `gorm:"column:base_currency;type:char(3)"`    // Валюта номинала / базовая валюта
	QuoteCurrency   *string    `gorm:"column:quote_currency;type:char(3)"`   // Валюта котировки
	CounterCurrency *string    `gorm:"column:counter_currency;type:char(3)"` // Сопряженная валюта
	MaturityDate    *time.Time `gorm:"column:maturity_date;type:date"`       // Дата погашения
	CouponDuration  *int       `gorm:"column:coupon_duration;type:int"`      // Длительность купона
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
