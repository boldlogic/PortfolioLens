package workers

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type RollForwardOtcRunner interface {
	DoRollForwardOtc(ctx context.Context) error
}

func NewRollForwardOtcWorker(svc RollForwardOtcRunner, logger *zap.Logger, interval time.Duration) Worker {
	return NewPeriodicWorker(
		"roll_forward_otc",
		"ошибка переноса OTC-лимитов",
		interval,
		func(ctx context.Context) error { return svc.DoRollForwardOtc(ctx) },
		logger,
	)
}
