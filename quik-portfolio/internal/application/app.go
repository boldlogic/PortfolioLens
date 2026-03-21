package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/commonconfig"
	logger "github.com/boldlogic/PortfolioLens/pkg/logger/zap"
	"github.com/boldlogic/PortfolioLens/pkg/metrics"
	"github.com/boldlogic/PortfolioLens/pkg/periodic"
	"github.com/boldlogic/PortfolioLens/pkg/transport/httpserver"
	"github.com/boldlogic/PortfolioLens/pkg/transport/httpserver/handler"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/config"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/repository"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/service"
	portfolioserver "github.com/boldlogic/PortfolioLens/quik-portfolio/internal/transport/http"
	v1 "github.com/boldlogic/PortfolioLens/quik-portfolio/internal/transport/http/v1"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/workers"
	"go.uber.org/zap"
)

const defaultConfigPath = "quik-portfolio/internal/configs/config.yaml"

type Application struct {
	cfg    *config.Config
	Logger *zap.Logger

	svc *service.Service

	errChan chan error
	wg      sync.WaitGroup
	repo    *repository.Repository

	server *httpserver.Server
}

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
		errChan: make(chan error, 1),
	}, nil
}

func (a *Application) Start(ctx context.Context) error {

	dsn := a.cfg.Db.GetDSN()
	repo, err := repository.NewRepository(ctx, dsn, a.Logger)
	if err != nil {
		return err
	}
	a.repo = repo

	a.svc = service.NewService(a.repo, a.Logger)

	runner := periodic.NewRunner(
		workers.NewRollForwardMoneyLimitsWorker(a.svc, a.Logger, 60*time.Second),
		workers.NewRollForwardSecurityLimitsWorker(a.svc, a.Logger, 60*time.Second),
		workers.NewRollForwardOtcWorker(a.svc, a.Logger, 60*time.Second),
		workers.NewActualizeFirmsWorker(a.svc, a.Logger, 60*time.Second),
	)
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		runner.Run(ctx)
	}()

	reg := metrics.New()
	commonHandler := handler.NewHandler()
	handler := v1.NewHandler(commonHandler, a.svc, a.Logger)
	r := portfolioserver.NewRouter(handler, a.Logger, reg)
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
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		_ = a.server.Shutdown(shutdownCtx)
	}

	a.wg.Wait()
	close(a.errChan)
	errWg.Wait()

	return appErr
}
