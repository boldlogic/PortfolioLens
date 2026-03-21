package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/boldlogic/PortfolioLens/instruments/internal/config"
	"github.com/boldlogic/PortfolioLens/instruments/internal/repository"
	"github.com/boldlogic/PortfolioLens/instruments/internal/service"
	instrumentserver "github.com/boldlogic/PortfolioLens/instruments/internal/transport/http"
	v1 "github.com/boldlogic/PortfolioLens/instruments/internal/transport/http/v1"
	"github.com/boldlogic/PortfolioLens/instruments/internal/workers"
	"github.com/boldlogic/PortfolioLens/pkg/commonconfig"
	logger "github.com/boldlogic/PortfolioLens/pkg/logger/zap"
	"github.com/boldlogic/PortfolioLens/pkg/metrics"
	"github.com/boldlogic/PortfolioLens/pkg/periodic"
	"github.com/boldlogic/PortfolioLens/pkg/transport/httpserver"
	"github.com/boldlogic/PortfolioLens/pkg/transport/httpserver/handler"
	"go.uber.org/zap"
)

type Application struct {
	cfg    *config.Config
	Logger *zap.Logger

	svc     *service.Service
	errChan chan error
	wg      sync.WaitGroup
	repo    *repository.Repository
	server  *httpserver.Server
}

const (
	defaultConfigPath  = "instruments/internal/configs/config.yaml"
	errChanBufSize     = 1
)

func New() (*Application, error) {
	configPath := commonconfig.GetConfigPath(defaultConfigPath)
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}
	log := logger.New(cfg.Log)
	return &Application{
		cfg:     cfg,
		Logger:  log,
		errChan: make(chan error, errChanBufSize),
	}, nil
}

func (a *Application) Start(ctx context.Context) error {
	dsn := a.cfg.Db.GetDSN()
	repo, err := repository.NewRepository(ctx, dsn, a.Logger)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	a.repo = repo

	a.svc = service.NewService(a.repo, a.repo, a.repo, a.Logger)

	if err = a.initDictionaries(ctx); err != nil {
		return err
	}

	runner := periodic.NewRunner(
		workers.NewActualizeRefsWorker(a.svc, a.Logger, 60*time.Second),
		workers.NewSaveInstrumentsWorker(a.svc, a.Logger, 1*time.Second),
	)
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		runner.Run(ctx)
	}()

	reg := metrics.New()
	commonHandler := handler.NewHandler()
	h := v1.NewHandler(commonHandler, a.svc, a.Logger)
	r := instrumentserver.NewRouter(h, a.Logger, reg)
	a.server = httpserver.NewServer(r, a.cfg.Server)

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.errChan <- fmt.Errorf("http server остановлен с ошибкой: %w", err)
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

func (a *Application) initDictionaries(ctx context.Context) error {
	if err := a.svc.ActualizeRefs(ctx); err != nil {
		a.Logger.Error("ошибка при инициализации справочников", zap.Error(err))
		return err
	}
	return nil
}
