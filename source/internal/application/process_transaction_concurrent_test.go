package application

import (
	"TestProject/source/internal/entities"
	"context"
	"sync"
	"testing"

	"github.com/shopspring/decimal"
)

func TestApplication_ProcessTransaction_ConcurrentDeposits(t *testing.T) {
	repo := newMockWalletRepo()
	app := &Application{WalletRepo: repo}
	ctx := context.Background()

	wallet := entities.NewWallet()
	wallet.Balance = decimal.NewFromFloat(1000.0)
	repo.wallets[wallet.ID] = wallet

	const numGoroutines = 100
	depositAmount := decimal.NewFromFloat(10.0)
	expectedBalance := decimal.NewFromFloat(1000.0).Add(depositAmount.Mul(decimal.NewFromInt(numGoroutines)))

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			transaction := entities.Transaction{
				WalletId:      wallet.ID,
				OperationType: "DEPOSIT",
				Amount:        depositAmount,
			}
			_, err := app.ProcessTransaction(ctx, transaction)
			if err != nil {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("unexpected error in concurrent transaction: %v", err)
	}

	finalWallet, _ := repo.GetByID(ctx, wallet.ID)
	if !finalWallet.Balance.Equal(expectedBalance) {
		t.Errorf("expected balance %s, got %s", expectedBalance.String(), finalWallet.Balance.String())
	}
}

func TestApplication_ProcessTransaction_ConcurrentMixed(t *testing.T) {
	repo := newMockWalletRepo()
	app := &Application{WalletRepo: repo}
	ctx := context.Background()

	wallet := entities.NewWallet()
	wallet.Balance = decimal.NewFromFloat(1000.0)
	repo.wallets[wallet.ID] = wallet

	const numDeposits = 50
	const numWithdraws = 30
	amount := decimal.NewFromFloat(10.0)

	var wg sync.WaitGroup
	errors := make(chan error, numDeposits+numWithdraws)

	for i := 0; i < numDeposits; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			transaction := entities.Transaction{
				WalletId:      wallet.ID,
				OperationType: "DEPOSIT",
				Amount:        amount,
			}
			_, err := app.ProcessTransaction(ctx, transaction)
			if err != nil {
				errors <- err
			}
		}()
	}

	for i := 0; i < numWithdraws; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			transaction := entities.Transaction{
				WalletId:      wallet.ID,
				OperationType: "WITHDRAW",
				Amount:        amount,
			}
			_, err := app.ProcessTransaction(ctx, transaction)
			if err != nil {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("unexpected error in concurrent transaction: %v", err)
	}

	expectedBalance := decimal.NewFromFloat(1000.0).Add(amount.Mul(decimal.NewFromInt(numDeposits - numWithdraws)))
	finalWallet, _ := repo.GetByID(ctx, wallet.ID)
	if !finalWallet.Balance.Equal(expectedBalance) {
		t.Errorf("expected balance %s, got %s", expectedBalance.String(), finalWallet.Balance.String())
	}
}
