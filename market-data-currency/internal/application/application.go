package application

import (
	"context"
	"sync"
	"time"

	"github.com/boldlogic/PortfolioLens/market-data-currency/internal/cbr"
	"github.com/boldlogic/PortfolioLens/market-data-currency/internal/client"
	"github.com/boldlogic/PortfolioLens/market-data-currency/internal/config"
	"github.com/boldlogic/PortfolioLens/market-data-currency/internal/repository"
	"github.com/boldlogic/PortfolioLens/market-data-currency/internal/service"
	currencyserver "github.com/boldlogic/PortfolioLens/market-data-currency/internal/transport/http"
	v1 "github.com/boldlogic/PortfolioLens/market-data-currency/internal/transport/http/v1"
	"github.com/boldlogic/PortfolioLens/market-data-currency/internal/workers"
	"github.com/boldlogic/PortfolioLens/pkg/commonconfig"
	logger "github.com/boldlogic/PortfolioLens/pkg/logger/zap"
	"github.com/boldlogic/PortfolioLens/pkg/metrics"
	"github.com/boldlogic/PortfolioLens/pkg/periodic"
	"github.com/boldlogic/PortfolioLens/pkg/transport/httpclient"
	"github.com/boldlogic/PortfolioLens/pkg/transport/httpclient/clientmetrics"
	"github.com/boldlogic/PortfolioLens/pkg/transport/httpserver"
	"github.com/boldlogic/PortfolioLens/pkg/transport/httpserver/handler"
	"go.uber.org/zap"
)

type Application struct {
	cfg     *config.Config
	logger  *zap.Logger
	svc     *service.Service
	repo    *repository.Repository
	server  *httpserver.Server
	errChan chan error
	wg      sync.WaitGroup
}

const defaultConfigPath = "market-data-currency/internal/configs/config.yaml"

func New() (*Application, error) {
	configPath := commonconfig.GetConfigPath(defaultConfigPath)
	cfg, err := config.Load(configPath)
	if err != nil {
		return &Application{}, err
	}
	log := logger.New(cfg.Log)
	return &Application{
		cfg:     cfg,
		logger:  log,
		errChan: make(chan error, 8),
	}, nil
}

func (a *Application) Start(ctx context.Context) error {
	dsn := a.cfg.Db.GetDSN()
	repo, err := repository.NewRepository(ctx, dsn, a.logger)
	if err != nil {
		return err
	}
	a.repo = repo

	reg := metrics.New()

	commonClient := httpclient.NewClient(a.cfg.Client)
	cbrMetrics := clientmetrics.NewMetrics(reg)
	httpClient := client.NewClient(commonClient, cbrMetrics, "cbr", a.logger)

	cbrParser := cbr.NewParser(a.logger)

	a.svc = service.NewService(httpClient, cbrParser, a.repo, a.repo, a.logger)

	if err = a.svc.InitCurrencyDictionary(ctx); err != nil {
		return err
	}

	runner := periodic.NewRunner(
		workers.NewFetchOneNewTaskWorker(a.svc, a.logger, 1*time.Second),
		workers.NewQuikFxCBRRatesSyncWorker(a.svc, a.logger, 60*time.Second),
	)
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		runner.Run(ctx)
	}()

	commonHandler := handler.NewHandler()

	handler := v1.NewHandler(commonHandler, a.svc, a.logger)
	r := currencyserver.NewRouter(handler, a.logger, reg)
	a.server = httpserver.NewServer(r, a.cfg.Server)

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		if err := a.server.ListenAndServe(); err != nil {
			a.logger.Error("HTTP-сервер завершился с ошибкой", zap.Error(err))
		}
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

	if a.server != nil {
		_ = a.server.Shutdown(context.Background())
	}

	a.wg.Wait()
	close(a.errChan)
	errWg.Wait()

	return appErr
}
