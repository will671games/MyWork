package main

import (
	"TestProject/source/config"
	"TestProject/source/internal/application"
	"TestProject/source/internal/storage"
	"TestProject/source/internal/transport"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/robbert229/fxslog"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {
	fx.New(fx.WithLogger(func(log *slog.Logger) fxevent.Logger {
		return &fxslog.SlogLogger{Logger: log}
	}), CreateApp()).Run()
}

func NewLogger() *slog.Logger {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(log)
	return log
}

func CreateApp() fx.Option {
	return fx.Options(
		application.Module,
		transport.Module,
		storage.Module,

		fx.Provide(
			NewLogger,
			config.NewDBConfig,
			config.NewHttpConfig,
		),
		fx.Invoke(
			func(echo *echo.Echo) {},
		),
		fx.Invoke(
			func(*application.Application) {},
		),
	)
}
