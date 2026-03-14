package workers

import (
	"context"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/periodic"
	"go.uber.org/zap"
)

type RollForwardOtcRunner interface {
	DoRollForwardOtc(ctx context.Context) error
}

func NewRollForwardOtcWorker(svc RollForwardOtcRunner, logger *zap.Logger, interval time.Duration) periodic.Worker {
	return periodic.NewPeriodicWorker(
		"roll_forward_otc",
		"ошибка переноса OTC-лимитов",
		interval,
		func(ctx context.Context) error { return svc.DoRollForwardOtc(ctx) },
		logger,
	)
}
