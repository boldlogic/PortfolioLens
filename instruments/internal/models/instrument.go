package models

import (
	"time"

	md "github.com/boldlogic/PortfolioLens/pkg/models"
)

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

	TradePointId uint8
	TradePoint   *md.TradePoint
}
