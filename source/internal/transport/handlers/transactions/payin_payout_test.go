package transactions

import (
	"TestProject/source/internal/entities"
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockTransactionApp struct {
	processFunc func(ctx interface{}, transaction entities.Transaction) (entities.Transaction, error)
}

func (m *mockTransactionApp) ProcessTransaction(ctx interface{}, transaction entities.Transaction) (entities.Transaction, error) {
	if m.processFunc != nil {
		return m.processFunc(ctx, transaction)
	}
	return transaction, nil
}

type testHandlers struct {
	log *slog.Logger
	app *mockTransactionApp
}

func (h *testHandlers) CreateTransaction(c echo.Context) error {
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
	transaction.WalletId = request.WalletId
	transaction.OperationType = request.OperationType
	transaction.Amount = decimal.NewFromFloat32(request.Amount)

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

func setupTestHandler() (*testHandlers, *echo.Echo) {
	mockApp := &mockTransactionApp{}
	handlers := &testHandlers{
		log: slog.Default(),
		app: mockApp,
	}

	e := echo.New()
	v := &testValidator{}
	e.Validator = v

	return handlers, e
}

type testValidator struct{}

func (v *testValidator) Validate(i interface{}) error {
	return nil
}

func TestHandlers_CreateTransaction_Success(t *testing.T) {
	handlers, e := setupTestHandler()

	expectedBalance := 1500.0
	handlers.app.processFunc = func(ctx interface{}, transaction entities.Transaction) (entities.Transaction, error) {
		transaction.Balance = decimal.NewFromFloat(expectedBalance)
		return transaction, nil
	}

	reqBody := Request{
		WalletId:      "123e4567-e89b-12d3-a456-426614174000",
		OperationType: "DEPOSIT",
		Amount:        500.0,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.CreateTransaction(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response Response
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, reqBody.WalletId, response.WalletId)
	assert.Equal(t, reqBody.OperationType, response.OperationType)
	assert.Equal(t, reqBody.Amount, response.Amount)
	assert.Equal(t, expectedBalance, response.Balance)
}

func TestHandlers_CreateTransaction_InsufficientFunds(t *testing.T) {
	handlers, e := setupTestHandler()

	handlers.app.processFunc = func(ctx interface{}, transaction entities.Transaction) (entities.Transaction, error) {
		return entities.Transaction{}, entities.ErrInsufficientFunds
	}

	reqBody := Request{
		WalletId:      "123e4567-e89b-12d3-a456-426614174000",
		OperationType: "WITHDRAW",
		Amount:        1000.0,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.CreateTransaction(c)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok, "expected HTTPError")
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
}

func TestHandlers_CreateTransaction_InvalidRequest(t *testing.T) {
	handlers, e := setupTestHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewReader([]byte("invalid json")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.CreateTransaction(c)
	require.Error(t, err)
}
