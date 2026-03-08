package application

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/boldlogic/PortfolioLens/instruments/internal/config"
	"github.com/boldlogic/PortfolioLens/instruments/internal/repository"
	"github.com/boldlogic/PortfolioLens/instruments/internal/service"
	"github.com/boldlogic/PortfolioLens/instruments/workers"
	logger "github.com/boldlogic/PortfolioLens/pkg/logger/zap"
	"github.com/boldlogic/PortfolioLens/pkg/periodic"
	"go.uber.org/zap"
)

type Application struct {
	cfg    *config.Config
	Logger *zap.Logger

	svc *service.Service

	errChan chan error
	wg      sync.WaitGroup
	repo    *repository.Repository

	httpWg sync.WaitGroup
}

func New() (*Application, error) {
	config, err := config.LoadConfig(defaultConfigPath)
	if err != nil {
		return &Application{}, err
	}
	log := logger.New(config.Log)
	return &Application{
		cfg:     config,
		Logger:  log,
		errChan: make(chan error, 8),
	}, nil
}

const defaultConfigPath = "instruments/internal/configs/config.yaml"

func (a *Application) Start(ctx context.Context) error {

	dsn := a.cfg.Db.GetDSN()
	repo, err := repository.NewRepository(ctx, dsn, a.Logger)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	a.repo = repo

	a.svc = service.NewService(ctx, a.repo, a.Logger)

	runner := periodic.NewRunner(
		workers.NewSaveInstrumentsWorker(a.svc, a.Logger, 1*time.Second),
	)
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		runner.Run(ctx)
	}()
	return nil
}

func (a *Application) Wait(ctx context.Context, cancel context.CancelFunc) error {
	var appErr error

	errWg := sync.WaitGroup{}

	errWg.Add(1)

	go func() {
		defer errWg.Done()

		for err := range a.errChan {
			cancel()
			appErr = err
		}
	}()

	<-ctx.Done()

	a.wg.Wait()
	close(a.errChan)
	errWg.Wait()

	return appErr
}
