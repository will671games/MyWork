package wallet

import (
	"TestProject/source/internal/entities"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID      string          `gorm:"primaryKey;type:uuid"`
	Balance decimal.Decimal `gorm:"type:decimal(15,2);default:0;not null"`
}

func FromEntity(entity entities.Wallet) (Wallet, error) {
	return Wallet{
		ID:      entity.ID,
		Balance: entity.Balance,
	}, nil
}

func ToEntity(dto Wallet) (entities.Wallet, error) {
	return entities.Wallet{
		ID:      dto.ID,
		Balance: dto.Balance,
	}, nil
}
