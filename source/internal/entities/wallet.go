package entities

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID      string          `json:"id"`
	Balance decimal.Decimal `json:"balance"`
}

func NewWallet() Wallet {
	return Wallet{
		ID:      uuid.NewString(),
		Balance: decimal.Zero,
	}
}
