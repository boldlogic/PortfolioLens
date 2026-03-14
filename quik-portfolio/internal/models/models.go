package models

import "time"

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
