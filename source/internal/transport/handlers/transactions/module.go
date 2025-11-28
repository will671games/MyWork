package transactions

import (
	"TestProject/source/internal/application"
	"go.uber.org/fx"
	"log/slog"
)

const moduleName = "transactions_handler"

var Module = fx.Module(
	moduleName,
	fx.Provide(
		NewHandlers,
	),
	fx.Decorate(
		func(log *slog.Logger) *slog.Logger {
			return log.With("module", moduleName)
		},
	),
)

type Handlers struct {
	log *slog.Logger
	app *application.Application
}

func NewHandlers(log *slog.Logger, app *application.Application) *Handlers {
	return &Handlers{
		log: log,
		app: app,
	}
}
