package workers

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type SaveInstrumentsRunner interface {
	SaveInstrument(ctx context.Context) error
}

func NewSaveInstrumentsWorker(svc SaveInstrumentsRunner, logger *zap.Logger, interval time.Duration) Worker {
	return NewPeriodicWorker(
		"save_instrument",
		"ошибка сохранения инструментов",
		interval,
		func(ctx context.Context) error { return svc.SaveInstrument(ctx) },
		logger,
	)
}
