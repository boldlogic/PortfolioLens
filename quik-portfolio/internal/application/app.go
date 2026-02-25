package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	logger "github.com/boldlogic/PortfolioLens/pkg/logger/zap"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/config"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/repository"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/service"
	httpserver "github.com/boldlogic/PortfolioLens/quik-portfolio/internal/transport/http"
	"go.uber.org/zap"
)

type Application struct {
	cfg    *config.Config
	Logger *zap.Logger

	svc *service.Service

	errChan chan error
	wg      sync.WaitGroup
	repo    *repository.Repository

	httpSrv *httpserver.Server
	httpWg  sync.WaitGroup
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

const defaultConfigPath = "quik-portfolio/internal/configs/config.yaml"

func (a *Application) Start(ctx context.Context) error {

	dsn := a.cfg.Db.GetDSN()
	repo, err := repository.NewRepository(ctx, dsn, a.Logger)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	a.repo = repo

	// i, err := a.repo.SelectNewCurrentQuote(ctx)
	// if err != nil {
	// 	return fmt.Errorf("%w", err)
	// }
	a.svc = service.NewService(ctx, a.repo, a.repo, a.repo, a.Logger)
	// for i := 0; i <= 60000; i++ {
	// 	svc.SaveInstrument(ctx)
	// }

	a.Logger.Debug("ok")
	//fmt.Println(i)

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.svc.RollForwardSecurityLimitsOtc(ctx)
	}()

	err = a.InitDictionaries(ctx)
	if err != nil {
		return err
	}

	handler := httpserver.NewHandler(a.svc, a.Logger)
	router := httpserver.NewRouter(handler, a.Logger, a.cfg)
	a.httpSrv = httpserver.New(router.Mux, a.cfg.Http, a.Logger)

	a.httpWg.Add(1)

	go func() {
		defer a.httpWg.Done()
		if err := a.httpSrv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
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
		timeout := time.Duration(a.cfg.Http.Timeout) * time.Second
		if timeout <= 0 {
			timeout = 20 * time.Second
		}
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), timeout)
		_ = a.httpSrv.Shutdown(shutdownCtx)
		shutdownCancel()
	}

	a.httpWg.Wait()
	a.wg.Wait()
	close(a.errChan)
	errWg.Wait()

	return appErr
}

func (a *Application) InitDictionaries(ctx context.Context) error {
	err := a.svc.ActualizeInstrumentTypes(ctx)
	if err != nil {
		a.Logger.Error("ошибка при инициализации справочника типов инструментов", zap.Error(err))
		return err
	}
	err = a.svc.ActualizeInstrumentSubTypes(ctx)
	if err != nil {
		a.Logger.Error("ошибка при инициализации справочника подтипов инструментов", zap.Error(err))
		return err
	}
	err = a.svc.ActualizeBoards(ctx)
	if err != nil {
		a.Logger.Error("ошибка при инициализации справочника классов инструментов", zap.Error(err))
		return err
	}
	return nil
}
