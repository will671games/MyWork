package application

import (
	"TestProject/source/internal/entities"
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/shopspring/decimal"
)

type mockWalletRepo struct {
	mu      sync.RWMutex
	wallets map[string]entities.Wallet
	errors  map[string]error
}

func newMockWalletRepo() *mockWalletRepo {
	return &mockWalletRepo{
		wallets: make(map[string]entities.Wallet),
		errors:  make(map[string]error),
	}
}

func (m *mockWalletRepo) Create(ctx context.Context, wallet entities.Wallet) (entities.Wallet, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err, ok := m.errors["create"]; ok {
		return entities.Wallet{}, err
	}
	m.wallets[wallet.ID] = wallet
	return wallet, nil
}

func (m *mockWalletRepo) GetByID(ctx context.Context, walletID string) (entities.Wallet, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if err, ok := m.errors["get"]; ok {
		return entities.Wallet{}, err
	}
	wallet, ok := m.wallets[walletID]
	if !ok {
		return entities.Wallet{}, errors.New("wallet not found")
	}
	return wallet, nil
}

func (m *mockWalletRepo) GetByIDForUpdate(ctx context.Context, walletID string) (entities.Wallet, error) {
	return m.GetByID(ctx, walletID)
}

func (m *mockWalletRepo) Update(ctx context.Context, wallet entities.Wallet) (entities.Wallet, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err, ok := m.errors["update"]; ok {
		return entities.Wallet{}, err
	}
	m.wallets[wallet.ID] = wallet
	return wallet, nil
}

func (m *mockWalletRepo) UpdateWithLock(ctx context.Context, walletID string, updateFn func(*entities.Wallet) error) (entities.Wallet, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err, ok := m.errors["updateWithLock"]; ok {
		return entities.Wallet{}, err
	}
	wallet, ok := m.wallets[walletID]
	if !ok {
		return entities.Wallet{}, errors.New("wallet not found")
	}

	err := updateFn(&wallet)
	if err != nil {
		return entities.Wallet{}, err
	}

	m.wallets[walletID] = wallet
	return wallet, nil
}

func TestApplication_ProcessTransaction_Deposit(t *testing.T) {
	repo := newMockWalletRepo()
	app := &Application{WalletRepo: repo}
	ctx := context.Background()

	wallet := entities.NewWallet()
	wallet.Balance = decimal.NewFromFloat(1000.0)
	repo.wallets[wallet.ID] = wallet

	transaction := entities.Transaction{
		WalletId:      wallet.ID,
		OperationType: "DEPOSIT",
		Amount:        decimal.NewFromFloat(500.0),
	}

	result, err := app.ProcessTransaction(ctx, transaction)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedBalance := decimal.NewFromFloat(1500.0)
	if !result.Balance.Equal(expectedBalance) {
		t.Errorf("expected balance 1500.0, got %s", result.Balance.String())
	}

	updatedWallet, _ := repo.GetByID(ctx, wallet.ID)
	if !updatedWallet.Balance.Equal(expectedBalance) {
		t.Errorf("expected wallet balance 1500.0, got %s", updatedWallet.Balance.String())
	}
}

func TestApplication_ProcessTransaction_Withdraw(t *testing.T) {
	repo := newMockWalletRepo()
	app := &Application{WalletRepo: repo}
	ctx := context.Background()

	wallet := entities.NewWallet()
	wallet.Balance = decimal.NewFromFloat(1000.0)
	repo.wallets[wallet.ID] = wallet

	transaction := entities.Transaction{
		WalletId:      wallet.ID,
		OperationType: "WITHDRAW",
		Amount:        decimal.NewFromFloat(300.0),
	}

	result, err := app.ProcessTransaction(ctx, transaction)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedBalance := decimal.NewFromFloat(700.0)
	if !result.Balance.Equal(expectedBalance) {
		t.Errorf("expected balance 700.0, got %s", result.Balance.String())
	}
}

func TestApplication_ProcessTransaction_InsufficientFunds(t *testing.T) {
	repo := newMockWalletRepo()
	app := &Application{WalletRepo: repo}
	ctx := context.Background()

	wallet := entities.NewWallet()
	wallet.Balance = decimal.NewFromFloat(100.0)
	repo.wallets[wallet.ID] = wallet

	transaction := entities.Transaction{
		WalletId:      wallet.ID,
		OperationType: "WITHDRAW",
		Amount:        decimal.NewFromFloat(500.0),
	}

	_, err := app.ProcessTransaction(ctx, transaction)
	if err == nil {
		t.Fatal("expected error for insufficient funds")
	}

	if !errors.Is(err, entities.ErrInsufficientFunds) {
		t.Errorf("expected ErrInsufficientFunds, got %v", err)
	}

	updatedWallet, _ := repo.GetByID(ctx, wallet.ID)
	expectedBalance := decimal.NewFromFloat(100.0)
	if !updatedWallet.Balance.Equal(expectedBalance) {
		t.Errorf("expected wallet balance to remain 100.0, got %s", updatedWallet.Balance.String())
	}
}

func TestApplication_ProcessTransaction_InvalidOperation(t *testing.T) {
	repo := newMockWalletRepo()
	app := &Application{WalletRepo: repo}
	ctx := context.Background()

	wallet := entities.NewWallet()
	wallet.Balance = decimal.NewFromFloat(1000.0)
	repo.wallets[wallet.ID] = wallet

	transaction := entities.Transaction{
		WalletId:      wallet.ID,
		OperationType: "INVALID",
		Amount:        decimal.NewFromFloat(100.0),
	}

	_, err := app.ProcessTransaction(ctx, transaction)
	if err == nil {
		t.Fatal("expected error for invalid operation")
	}

	if !errors.Is(err, entities.ErrInvalidOperation) {
		t.Errorf("expected ErrInvalidOperation, got %v", err)
	}
}

func TestApplication_ProcessTransaction_WalletNotFound(t *testing.T) {
	repo := newMockWalletRepo()
	app := &Application{WalletRepo: repo}
	ctx := context.Background()

	transaction := entities.Transaction{
		WalletId:      "non-existent-id",
		OperationType: "DEPOSIT",
		Amount:        decimal.NewFromFloat(100.0),
	}

	_, err := app.ProcessTransaction(ctx, transaction)
	if err == nil {
		t.Fatal("expected error for wallet not found")
	}
}
