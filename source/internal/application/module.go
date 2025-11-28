package application

import (
	"TestProject/source/internal/entities"
	"context"
	"log/slog"

	"go.uber.org/fx"
)

const moduleName = "application"

var Module = fx.Module(
	moduleName,
	fx.Provide(
		New,
	),
	fx.Decorate(
		func(log *slog.Logger) *slog.Logger {
			return log.With(slog.String("module", moduleName))
		},
	),
)

type Application struct {
	log        *slog.Logger
	WalletRepo WalletRepo
}

func New(
	log *slog.Logger,
	walletRepo WalletRepo,
) *Application {
	return &Application{
		log:        log,
		WalletRepo: walletRepo,
	}
}

type WalletRepo interface {
	Create(ctx context.Context, wallet entities.Wallet) (entities.Wallet, error)
	GetByID(ctx context.Context, walletID string) (entities.Wallet, error)
	GetByIDForUpdate(ctx context.Context, walletID string) (entities.Wallet, error)
	Update(ctx context.Context, wallet entities.Wallet) (entities.Wallet, error)
	UpdateWithLock(ctx context.Context, walletID string, updateFn func(*entities.Wallet) error) (entities.Wallet, error)
}
