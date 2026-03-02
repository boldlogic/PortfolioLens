package workers

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type SaveInstrumentTypesFromQuotesRunner interface {
	SaveInstrumentTypesFromQuotes(ctx context.Context) error
}

func NewSaveInstrumentTypesFromQuotesWorker(svc SaveInstrumentTypesFromQuotesRunner, logger *zap.Logger, interval time.Duration) Worker {
	return NewPeriodicWorker(
		"actualize_instrument_types",
		"ошибка сохранения типов активов",
		interval,
		func(ctx context.Context) error { return svc.SaveInstrumentTypesFromQuotes(ctx) },
		logger,
	)
}
