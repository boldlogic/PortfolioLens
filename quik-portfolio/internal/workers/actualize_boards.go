package workers

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type ActualizeBoardsRunner interface {
	ActualizeBoards(ctx context.Context) error
}

func NewActualizeBoardsWorker(svc ActualizeBoardsRunner, logger *zap.Logger, interval time.Duration) Worker {
	return NewPeriodicWorker(
		"actualize_boards",
		"ошибка сохранения классов",
		interval,
		func(ctx context.Context) error { return svc.ActualizeBoards(ctx) },
		logger,
	)
}
