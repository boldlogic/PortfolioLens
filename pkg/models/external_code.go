package models

type ExternalCode struct {
	ExtId            int32
	ExternalSystemId uint8
	Code             string
	Type             uint8
	IntId            int64
}
