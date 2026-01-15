package application

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/boldlogic/cbr-market-data-worker/internal/config"
	httpserver "github.com/boldlogic/cbr-market-data-worker/internal/transport/http"
	"github.com/boldlogic/cbr-market-data-worker/pkg/logger"
	"github.com/sirupsen/logrus"
)

type Application struct {
	cfg       *config.Config
	Log       *logrus.Logger
	logCloser io.Closer
	errChan   chan error
	wg        sync.WaitGroup
	httpSrv   *httpserver.Server
}

func New() *Application {
	return &Application{
		Log:     logrus.New(),
		errChan: make(chan error, 8),
	}
}

func (a *Application) initConfig() error {
	var err error

	a.cfg, err = config.ParseConfig()
	if err != nil {
		return fmt.Errorf("не удалось прочитать конфиг-файл: %w", err)
	}
	// err = a.cfg.Log.Validate()
	// if err != nil {
	// 	return fmt.Errorf("не удалось загрузить конфигурацию: %w", err)
	// }
	return nil
}

func (a *Application) Start(ctx context.Context) error {
	if err := a.initConfig(); err != nil {
		return fmt.Errorf("не удалось проинициализировать конфиг: %w", err)
	}
	errs := a.cfg.Validate()
	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("некорректный конфиг: %w", err)
	}

	log, logCloser, err := logger.New(a.cfg.Log)
	if err != nil {
		return fmt.Errorf("не удалось инициализировать логгер: %w", err)
	}
	a.Log = log
	a.logCloser = logCloser

	// HTTP transport
	handler := httpserver.NewHandler(a.Log)
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
