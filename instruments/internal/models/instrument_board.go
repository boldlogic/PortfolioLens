package models

import (
	"github.com/boldlogic/PortfolioLens/pkg/models/quik"
)

type InstrumentBoard struct {
	InstrumentId int
	Instrument   *Instrument

	BoardId uint8
	Board   *quik.Board

	TypeId uint8
	Type   *quik.InstrumentType

	SubTypeId *uint8
	SubType   *quik.InstrumentSubType

	CurrencyId        *int
	BaseCurrencyId    *int
	QuoteCurrencyId   *int
	CounterCurrencyId *int
	IsPrimary         bool
}
