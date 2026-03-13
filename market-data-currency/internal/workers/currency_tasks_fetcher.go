package workers

import (
	"context"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/periodic"
	"go.uber.org/zap"
)

type FetchOneNewTaskRunner interface {
	FetchOneNewTask(ctx context.Context) error
}

func NewFetchOneNewTaskWorker(svc FetchOneNewTaskRunner, logger *zap.Logger, interval time.Duration) periodic.Worker {
	return periodic.NewPeriodicWorker(
		"currency_tasks_fetcher",
		"ошибка обработки задачи",
		interval,
		func(ctx context.Context) error { return svc.FetchOneNewTask(ctx) },
		logger,
	)
}
