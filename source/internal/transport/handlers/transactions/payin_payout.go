package transactions

import (
	"TestProject/source/internal/entities"
	"errors"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

func (h *Handlers) CreateTransaction(c echo.Context) error {
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

	transaction := entities.NewTransaction()

	transaction, err = FillTransactionFromRequest(transaction, request)
	if err != nil {
		logger.ErrorContext(ctx, "error filling transaction",
			slog.String("error", err.Error()),
		)
		return echo.ErrBadRequest.SetInternal(err)
	}

	transaction, err = h.app.ProcessTransaction(ctx, transaction)
	if err != nil {
		logger.ErrorContext(ctx, "error processing transaction", slog.String("error", err.Error()))
		if errors.Is(err, entities.ErrInsufficientFunds) {
			return echo.NewHTTPError(400, "insufficient funds").SetInternal(err)
		}
		return echo.ErrInternalServerError.SetInternal(err)
	}

	response := EntityToResponse(transaction)

	return c.JSON(201, response)
}

func FillTransactionFromRequest(transaction entities.Transaction, request Request) (entities.Transaction, error) {
	transaction.WalletId = request.WalletId
	transaction.OperationType = request.OperationType
	transaction.Amount = decimal.NewFromFloat32(request.Amount)
	return transaction, nil
}
