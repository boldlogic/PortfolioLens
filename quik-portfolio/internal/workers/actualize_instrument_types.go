package workers

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type ActualizeInstrumentTypesRunner interface {
	ActualizeInstrumentTypes(ctx context.Context) error
}

func NewActualizeInstrumentTypesWorker(svc ActualizeInstrumentTypesRunner, logger *zap.Logger, interval time.Duration) Worker {
	return NewPeriodicWorker(
		"actualize_instrument_types",
		"ошибка сохранения типов активов",
		interval,
		func(ctx context.Context) error { return svc.ActualizeInstrumentTypes(ctx) },
		logger,
	)
}
