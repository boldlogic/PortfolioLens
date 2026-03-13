package cbr

import (
	"go.uber.org/zap"
)

type Parser struct {
	logger *zap.Logger
}

func NewParser(logger *zap.Logger) *Parser {
	return &Parser{logger: logger}
}

const RateDateFormat = "02.01.2006"
