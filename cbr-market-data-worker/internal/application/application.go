package application

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/client"
	"github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/config"
	"github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/service"
	"github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/service/request_catalog"
	"github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/storage"
	httpserver "github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/transport/http"
	logger "github.com/boldlogic/PortfolioLens/pkg/logger/logrus"
	"github.com/sirupsen/logrus"
)

type Application struct {
	cfg       *config.Config
	Log       *logrus.Logger
	logCloser io.Closer
	errChan   chan error
	wg        sync.WaitGroup
	httpSrv   *httpserver.Server
	storage   *storage.Storage
}

func New() *Application {
	return &Application{
		Log:     logrus.New(),
		errChan: make(chan error, 8),
	}
}

const defaultConfigPath = "cbr-market-data-worker/internal/configs/config.yaml"

func (a *Application) Start(ctx context.Context) error {

	var err error
	a.cfg, err = config.LoadConfig(defaultConfigPath)
	if err != nil {
		return err
	}
	log, logCloser, err := logger.New(a.cfg.Log)
	if err != nil {
		return fmt.Errorf("не удалось инициализировать логгер: %w", err)
	}
	a.Log = log
	a.logCloser = logCloser

	dsn := a.cfg.Db.GetDSN()
	db, err := storage.NewStorage(dsn)
	if err != nil {
		defer a.Close()
		return fmt.Errorf("%w", err)
	}

	err = db.Migrate()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	httpClient := client.NewClient(a.cfg.Client)
	registry := request_catalog.NewProvider(a.cfg.Client)

	svc := service.NewService(httpClient, registry, db, db, db, a.Log)
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		svc.StartWorker(ctx)
	}()
	handler := httpserver.NewHandler(a.Log, *svc)
	router := httpserver.NewRouter(handler, a.Log, a.cfg)
	a.httpSrv = httpserver.NewServer(router.Mux, a.cfg.Server, a.Log)

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		if err := a.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

	if a.httpSrv != nil {
		timeout := time.Duration(a.cfg.Server.Timeout) * time.Second
		if timeout <= 0 {
			timeout = 10 * time.Second
		}
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), timeout)
		_ = a.httpSrv.Shutdown(shutdownCtx)
		shutdownCancel()
	}

	a.wg.Wait()
	close(a.errChan)
	errWg.Wait()

	return appErr
}

func (a *Application) Close() {
	if a.logCloser != nil {
		_ = a.logCloser.Close()
		a.logCloser = nil
	}
}
