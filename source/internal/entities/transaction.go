package entities

import "github.com/shopspring/decimal"

type Transaction struct {
	WalletId      string          `json:"wallet_id"`
	OperationType string          `json:"operation_type"`
	Amount        decimal.Decimal `json:"amount"`
	Balance       decimal.Decimal `json:"balance,omitempty"`
}

func NewTransaction() Transaction {
	return Transaction{}
}
