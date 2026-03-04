package workers

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Worker interface {
	Name() string
	Run(ctx context.Context)
}

type Job func(ctx context.Context) error

type PeriodicWorker struct {
	name       string
	interval   time.Duration
	retryDelay time.Duration
	job        Job
	logger     *zap.Logger
	errMsg     string
}

func NewPeriodicWorker(name, errMsg string, interval time.Duration, job Job, logger *zap.Logger) *PeriodicWorker {
	if interval <= 0 {
		interval = 60 * time.Second
	}
	return &PeriodicWorker{
		name:       name,
		interval:   interval,
		retryDelay: 5 * time.Second,
		job:        job,
		logger:     logger,
		errMsg:     errMsg,
	}
}

func (w *PeriodicWorker) Name() string { return w.name }

func (w *PeriodicWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := w.job(ctx)
			if err != nil {
				w.logger.Error(w.errMsg, zap.Error(err))
				select {
				case <-time.After(w.retryDelay):
				case <-ctx.Done():
					return
				}
			}
		}
	}
}

type Runner struct {
	workers []Worker
	wg      sync.WaitGroup
}

func NewRunner(workers ...Worker) *Runner {
	return &Runner{workers: workers}
}

func (r *Runner) Run(ctx context.Context) {
	for _, w := range r.workers {
		w := w
		r.wg.Add(1)
		go func() {
			defer r.wg.Done()
			w.Run(ctx)
		}()
	}
	<-ctx.Done()
	r.wg.Wait()
}
