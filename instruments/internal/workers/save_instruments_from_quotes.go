package workers

import (
	"context"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/periodic"
	"go.uber.org/zap"
)

type SaveInstrumentsRunner interface {
	SaveInstrument(ctx context.Context) error
}

func NewSaveInstrumentsWorker(svc SaveInstrumentsRunner, logger *zap.Logger, interval time.Duration) periodic.Worker {
	return periodic.NewPeriodicWorker(
		"save_instrument",
		"ошибка сохранения инструментов",
		interval,
		func(ctx context.Context) error { return svc.SaveInstrument(ctx) },
		logger,
	)
}
