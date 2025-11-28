package transport

import (
	"TestProject/source/config"
	"TestProject/source/internal/transport/handlers/transactions"
	"TestProject/source/internal/transport/handlers/wallet"
	"context"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

const moduleName = "http_api"

var Module = fx.Module(
	moduleName,
	transactions.Module,
	wallet.Module,
	fx.Provide(
		NewHandlers,
		NewEchoServer,
	),
	fx.Decorate(
		func(log *slog.Logger) *slog.Logger {
			return log.With("module", moduleName)
		},
	),
)

func NewEchoServer(
	lc fx.Lifecycle,
	handlers *Handlers,
	cfg config.HttpConfig,
) *echo.Echo {
	echoRouter := echo.New()

	echoRouter.HideBanner = true

	v := validator.New()
	echoRouter.Validator = &CustomValidator{validator: v}

	handlers.Registry(echoRouter)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := echoRouter.Start(fmt.Sprintf(":%v", cfg.Port))
				if err != nil {
					slog.Error("starting http echo server", slog.Any("error", err))
					return
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return echoRouter.Shutdown(ctx)
		},
	})

	return echoRouter
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
