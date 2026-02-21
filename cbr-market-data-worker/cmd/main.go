package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/boldlogic/PortfolioLens/cbr-market-data-worker/internal/application"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
	)
	defer cancel()

	app := application.New()
	defer app.Close()

	err := app.Start(ctx)
	if err != nil {
		app.Log.Fatalf("не удалось запустить приложение: %v", err)
	}

	app.Log.Infof("Приложение запущено. Ожидаю сигнал завершения...")

	err = app.Wait(ctx, cancel)
	if err != nil {
		app.Log.Fatalf("Приложение завершилось с ошибкой. Последняя ошибка: %v", err)
	}

	app.Log.Infof("Приложение завершилось без ошибок")
}
