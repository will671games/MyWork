package wallet

import (
	"TestProject/source/internal/entities"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

func (h *Handlers) CreateWallet(c echo.Context) error {
	ctx := c.Request().Context()

	logger := h.log

	var request Request

	err := c.Bind(&request)
	if err != nil {
		logger.ErrorContext(ctx, "error binding request", slog.String("error", err.Error()))

		return echo.ErrBadRequest.SetInternal(err)
	}

	err = c.Validate(&request)
	if err != nil {
		logger.ErrorContext(ctx, "error validating request", slog.String("error", err.Error()))
		return echo.ErrBadRequest.SetInternal(err)
	}

	wallet := entities.NewWallet()

	wallet, err = FillWalletFromRequest(wallet, request)
	if err != nil {
		logger.ErrorContext(ctx, "error when filling in the wallet's fields", slog.String("error", err.Error()))

		return echo.ErrConflict.SetInternal(err)
	}

	wallet, err = h.app.WalletRepo.Create(ctx, wallet)
	if err != nil {
		logger.ErrorContext(ctx, "error creating wallet", slog.String("error", err.Error()))

		return echo.ErrInternalServerError.SetInternal(err)
	}

	response := EntityToResponse(wallet)

	return c.JSON(200, response)
}

func FillWalletFromRequest(wallet entities.Wallet, request Request) (entities.Wallet, error) {
	wallet.Balance = decimal.NewFromFloat(request.Balance)
	return wallet, nil
}
