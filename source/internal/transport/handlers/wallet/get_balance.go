package wallet

import (
	"github.com/labstack/echo/v4"
	"log/slog"
)

func (h *Handlers) GetBalance(c echo.Context) error {
	ctx := c.Request().Context()

	logger := h.log

	walletID := c.Param("walletId")

	if walletID == "" {
		logger.ErrorContext(ctx, "error wallet id is not provided")
		return echo.ErrBadRequest
	}

	wallet, err := h.app.WalletRepo.GetByID(ctx, walletID)
	if err != nil {
		logger.ErrorContext(ctx, "wallet not found", slog.String("error", err.Error()))
		return echo.ErrNotFound.SetInternal(err)
	}

	response := EntityToResponse(wallet)

	return c.JSON(200, response)
}
