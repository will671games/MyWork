package transport

import (
	"TestProject/source/internal/transport/handlers/transactions"
	"TestProject/source/internal/transport/handlers/wallet"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handlers struct {
	logger       *slog.Logger
	transactions *transactions.Handlers
	wallet       *wallet.Handlers
}

func NewHandlers(
	logger *slog.Logger,
	transactionHandlers *transactions.Handlers,
	walletHandlers *wallet.Handlers,
) *Handlers {
	return &Handlers{
		logger:       logger,
		transactions: transactionHandlers,
		wallet:       walletHandlers,
	}
}

func (h *Handlers) Registry(e *echo.Echo) {
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	api := e.Group("/api/v1")

	api.POST("/wallets", h.wallet.CreateWallet)

	api.GET("/wallets/:walletId", h.wallet.GetBalance)

	api.POST("/wallet", h.transactions.CreateTransaction)
}
