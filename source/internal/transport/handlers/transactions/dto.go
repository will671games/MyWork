package transactions

import "TestProject/source/internal/entities"

type Request struct {
	WalletId      string  `json:"wallet_id" validate:"required,uuid"`
	OperationType string  `json:"operation_type" validate:"required,oneof=DEPOSIT WITHDRAW"`
	Amount        float32 `json:"amount" validate:"required,gt=0"`
}

type Response struct {
	WalletId      string  `json:"wallet_id"`
	OperationType string  `json:"operation_type"`
	Amount        float32 `json:"amount"`
	Balance       float64 `json:"balance,omitempty"`
}

func EntityToResponse(transaction entities.Transaction) Response {
	amount64, _ := transaction.Amount.Float64()
	balance, _ := transaction.Balance.Float64()
	return Response{
		WalletId:      transaction.WalletId,
		OperationType: transaction.OperationType,
		Amount:        float32(amount64),
		Balance:       balance,
	}
}
