package workers

import (
	"context"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/periodic"
	"go.uber.org/zap"
)

type QuikFxCBRRatesSyncRunner interface {
	SaveFxCBRRatesFromQuik(ctx context.Context) error
}

func NewQuikFxCBRRatesSyncWorker(svc QuikFxCBRRatesSyncRunner, logger *zap.Logger, interval time.Duration) periodic.Worker {
	return periodic.NewPeriodicWorker(
		"quik_fx_cbr_rates_sync",
		"ошибка загрузки кросс-курсов из квик",
		interval,
		func(ctx context.Context) error { return svc.SaveFxCBRRatesFromQuik(ctx) },
		logger,
	)
}
