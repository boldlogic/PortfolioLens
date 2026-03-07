package application

import (
	"context"
	"fmt"
	"sync"

	"github.com/boldlogic/PortfolioLens/market-data-currency/internal/client"
	"github.com/boldlogic/PortfolioLens/market-data-currency/internal/config"
	"github.com/boldlogic/PortfolioLens/market-data-currency/internal/repository"
	"github.com/boldlogic/PortfolioLens/market-data-currency/internal/service"
	"github.com/boldlogic/PortfolioLens/market-data-currency/internal/service/request_catalog"
	logger "github.com/boldlogic/PortfolioLens/pkg/logger/zap"
	"go.uber.org/zap"
)

type Application struct {
	cfg     *config.Config
	logger  *zap.Logger
	svc     *service.Service
	repo    *repository.Repository
	errChan chan error
	wg      sync.WaitGroup
}

const defaultConfigPath = "market-data-currency/internal/configs/config.yaml"

func New() (*Application, error) {
	config, err := config.LoadConfig(defaultConfigPath)
	if err != nil {
		return &Application{}, err
	}
	log := logger.New(config.Log)
	return &Application{
		cfg:     config,
		logger:  log,
		errChan: make(chan error, 8),
	}, nil
}

func (a *Application) Start(ctx context.Context) error {

	dsn := a.cfg.Db.GetDSN()
	repo, err := repository.NewRepository(ctx, dsn, a.logger)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	a.repo = repo
	httpClient := client.NewClient(a.cfg.Client)
	registry := request_catalog.NewProvider(a.cfg.Client)

	a.svc = service.NewService(ctx, httpClient, a.repo, a.repo, registry, a.logger)

	// err = a.svc.InitCurrencyDictionary(ctx)
	// if err != nil {
	// 	return err
	// }
	_ = a.svc.FetchOneNewTask(ctx)
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
