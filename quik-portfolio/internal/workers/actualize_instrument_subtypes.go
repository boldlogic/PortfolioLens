package workers

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type ActualizeInstrumentSubTypesRunner interface {
	ActualizeInstrumentSubTypes(ctx context.Context) error
}

func NewActualizeInstrumentSubTypesWorker(svc ActualizeInstrumentSubTypesRunner, logger *zap.Logger, interval time.Duration) Worker {
	return NewPeriodicWorker(
		"actualize_instrument_sub_types",
		"ошибка сохранения типов активов",
		interval,
		func(ctx context.Context) error { return svc.ActualizeInstrumentSubTypes(ctx) },
		logger,
	)
}
