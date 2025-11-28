package wallet

import "TestProject/source/internal/entities"

type Request struct {
	Balance float64 `json:"balance" validate:"required,gte=0"`
}

type Response struct {
	WalletId string  `json:"walletId"`
	Balance  float64 `json:"balance"`
}

func EntityToResponse(wallet entities.Wallet) Response {
	balance, _ := wallet.Balance.Float64()
	return Response{
		WalletId: wallet.ID,
		Balance:  balance,
	}
}
