package workers

import (
	"context"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/periodic"
	"go.uber.org/zap"
)

type ActualizeRefsRunner interface {
	ActualizeRefs(ctx context.Context) error
}

func NewActualizeRefsWorker(svc ActualizeRefsRunner, logger *zap.Logger, interval time.Duration) periodic.Worker {
	return periodic.NewPeriodicWorker(
		"actualize_refs",
		"ошибка актуализации справочников",
		interval,
		func(ctx context.Context) error { return svc.ActualizeRefs(ctx) },
		logger,
	)
}
