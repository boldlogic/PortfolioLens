package models

type InstrumentBoard struct {
	InstrumentId int
	Instrument   *Instrument

	BoardId uint8
	Board   *Board

	TypeId uint8
	Type   *InstrumentType

	SubTypeId *uint8
	SubType   *InstrumentSubType

	CurrencyId        *int
	BaseCurrencyId    *int
	QuoteCurrencyId   *int
	CounterCurrencyId *int
	IsPrimary         bool
}
