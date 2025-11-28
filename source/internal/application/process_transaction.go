package application

import (
	"TestProject/source/internal/entities"
	"context"
	"fmt"
)

const (
	deposit  = "DEPOSIT"
	withdraw = "WITHDRAW"
)

func (a *Application) ProcessTransaction(ctx context.Context, transaction entities.Transaction) (entities.Transaction, error) {
	wallet, err := a.WalletRepo.UpdateWithLock(ctx, transaction.WalletId, func(w *entities.Wallet) error {
		switch transaction.OperationType {
		case deposit:
			w.Balance = w.Balance.Add(transaction.Amount)
		case withdraw:
			if w.Balance.LessThan(transaction.Amount) {
				return entities.ErrInsufficientFunds
			}
			w.Balance = w.Balance.Sub(transaction.Amount)
		default:
			return fmt.Errorf("%w: %s", entities.ErrInvalidOperation, transaction.OperationType)
		}
		return nil
	})
	if err != nil {
		return entities.Transaction{}, fmt.Errorf("error processing transaction: %w", err)
	}

	transaction.Balance = wallet.Balance

	return transaction, nil
}
