package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/boldlogic/PortfolioLens/pkg/commonconfig"
	logger "github.com/boldlogic/PortfolioLens/pkg/logger/zap"
	"github.com/boldlogic/PortfolioLens/pkg/metrics"
	"github.com/boldlogic/PortfolioLens/pkg/transport/httpserver"
	"github.com/boldlogic/PortfolioLens/pkg/transport/httpserver/handler"
	"github.com/boldlogic/PortfolioLens/task-manager/internal/config"
	"github.com/boldlogic/PortfolioLens/task-manager/internal/repository"
	"github.com/boldlogic/PortfolioLens/task-manager/internal/service"
	taskserver "github.com/boldlogic/PortfolioLens/task-manager/internal/transport/http"
	v1 "github.com/boldlogic/PortfolioLens/task-manager/internal/transport/http/v1"
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

const (
	defaultConfigPath  = "task-manager/internal/configs/config.yaml"
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
		logger:  log,
		errChan: make(chan error, errChanBufSize),
	}, nil
}

func (a *Application) Start(ctx context.Context) error {
	dsn := a.cfg.Db.GetDSN()
	repo, err := repository.NewRepository(ctx, dsn, a.logger)
	if err != nil {
		return err
	}
	a.repo = repo

	a.svc = service.NewService(ctx, a.repo, a.logger)

	reg := metrics.New()
	commonHandler := handler.NewHandler()

	handler := v1.NewHandler(commonHandler, a.svc, a.logger)
	r := taskserver.NewRouter(handler, a.logger, reg)
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
