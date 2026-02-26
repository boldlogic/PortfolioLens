package models

type Board struct {
	Id           uint8
	Code         string
	Name         string
	IsTraded     bool
	TradePointId *uint8
	TradePoint   *TradePoint
}
