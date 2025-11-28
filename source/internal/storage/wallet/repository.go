package wallet

import (
	"TestProject/source/internal/entities"
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, wallet entities.Wallet) (entities.Wallet, error) {
	dto, err := FromEntity(wallet)
	if err != nil {
		return entities.Wallet{}, fmt.Errorf("wallet to entity error: %w", err)

	}

	if err := r.db.WithContext(ctx).Create(&dto).Error; err != nil {
		return entities.Wallet{}, fmt.Errorf("failed to create wallet: %w", err)
	}

	entity, err := ToEntity(dto)
	if err != nil {
		return entities.Wallet{}, fmt.Errorf("wallet to entity error: %w", err)

	}
	return entity, nil
}

func (r *Repo) GetByID(ctx context.Context, walletID string) (entities.Wallet, error) {
	var dto Wallet

	err := r.db.WithContext(ctx).Where("id = ?", walletID).First(&dto).Error
	if err != nil {
		return entities.Wallet{}, fmt.Errorf("db.First error: %w", err)
	}

	entity, err := ToEntity(dto)
	if err != nil {
		return entities.Wallet{}, fmt.Errorf("wallet.ToEntity error: %w", err)

	}

	return entity, nil
}

func (r *Repo) GetByIDForUpdate(ctx context.Context, walletID string) (entities.Wallet, error) {
	var dto Wallet

	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", walletID).
		First(&dto).Error
	if err != nil {
		return entities.Wallet{}, fmt.Errorf("db.First error: %w", err)
	}

	entity, err := ToEntity(dto)
	if err != nil {
		return entities.Wallet{}, fmt.Errorf("wallet.ToEntity error: %w", err)
	}

	return entity, nil
}

func (r *Repo) UpdateWithLock(ctx context.Context, walletID string, updateFn func(*entities.Wallet) error) (entities.Wallet, error) {
	var result entities.Wallet

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var dto Wallet
		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", walletID).
			First(&dto).Error
		if err != nil {
			return fmt.Errorf("db.First error: %w", err)
		}

		entity, err := ToEntity(dto)
		if err != nil {
			return fmt.Errorf("wallet.ToEntity error: %w", err)
		}

		err = updateFn(&entity)
		if err != nil {
			return err
		}

		dto, err = FromEntity(entity)
		if err != nil {
			return fmt.Errorf("wallet from entity error: %w", err)
		}

		err = tx.Save(&dto).Error
		if err != nil {
			return fmt.Errorf("db.Save error: %w", err)
		}

		result = entity
		return nil
	})

	if err != nil {
		return entities.Wallet{}, err
	}

	return result, nil
}

func (r *Repo) Update(ctx context.Context, wallet entities.Wallet) (entities.Wallet, error) {
	dto, err := FromEntity(wallet)
	if err != nil {
		return entities.Wallet{}, fmt.Errorf("wallet from entity error: %w", err)
	}

	err = r.db.WithContext(ctx).Save(&dto).Error
	if err != nil {
		return entities.Wallet{}, fmt.Errorf("db.Save error: %w", err)
	}

	entity, err := ToEntity(dto)
	if err != nil {
		return entities.Wallet{}, fmt.Errorf("wallet to entity error: %w", err)
	}

	return entity, nil
}
