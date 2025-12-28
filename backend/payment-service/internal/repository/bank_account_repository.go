package repository

import (
	"bus-booking/payment-service/internal/model"
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BankAccountRepository interface {
	CreateBankAccount(ctx context.Context, account *model.BankAccount) error
	GetUserBankAccounts(ctx context.Context, userID uuid.UUID) ([]*model.BankAccount, error)
	GetPrimaryBankAccount(ctx context.Context, userID uuid.UUID) (*model.BankAccount, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.BankAccount, error)
	UpdateBankAccount(ctx context.Context, account *model.BankAccount) error
	SetPrimaryBankAccount(ctx context.Context, accountID uuid.UUID, userID uuid.UUID) error
	DeleteBankAccount(ctx context.Context, accountID uuid.UUID, userID uuid.UUID) error
}

type bankAccountRepositoryImpl struct {
	db *gorm.DB
}

func NewBankAccountRepository(db *gorm.DB) BankAccountRepository {
	return &bankAccountRepositoryImpl{db: db}
}

func (r *bankAccountRepositoryImpl) CreateBankAccount(ctx context.Context, account *model.BankAccount) error {
	if err := r.db.WithContext(ctx).Create(account).Error; err != nil {
		return fmt.Errorf("failed to create bank account: %w", err)
	}
	return nil
}

func (r *bankAccountRepositoryImpl) GetUserBankAccounts(ctx context.Context, userID uuid.UUID) ([]*model.BankAccount, error) {
	var accounts []*model.BankAccount
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_primary DESC, created_at DESC").
		Find(&accounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get user bank accounts: %w", err)
	}
	return accounts, nil
}

func (r *bankAccountRepositoryImpl) GetPrimaryBankAccount(ctx context.Context, userID uuid.UUID) (*model.BankAccount, error) {
	var account model.BankAccount
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_primary = ?", userID, true).
		First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no primary bank account found")
		}
		return nil, fmt.Errorf("failed to get primary bank account: %w", err)
	}
	return &account, nil
}

func (r *bankAccountRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.BankAccount, error) {
	var account model.BankAccount
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("bank account not found")
		}
		return nil, fmt.Errorf("failed to get bank account: %w", err)
	}
	return &account, nil
}

func (r *bankAccountRepositoryImpl) UpdateBankAccount(ctx context.Context, account *model.BankAccount) error {
	if err := r.db.WithContext(ctx).Save(account).Error; err != nil {
		return fmt.Errorf("failed to update bank account: %w", err)
	}
	return nil
}

func (r *bankAccountRepositoryImpl) SetPrimaryBankAccount(ctx context.Context, accountID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Unset all primary flags for this user
		if err := tx.Model(&model.BankAccount{}).
			Where("user_id = ?", userID).
			Update("is_primary", false).Error; err != nil {
			return fmt.Errorf("failed to unset primary flags: %w", err)
		}

		// Set the new primary
		if err := tx.Model(&model.BankAccount{}).
			Where("id = ? AND user_id = ?", accountID, userID).
			Update("is_primary", true).Error; err != nil {
			return fmt.Errorf("failed to set primary: %w", err)
		}

		return nil
	})
}

func (r *bankAccountRepositoryImpl) DeleteBankAccount(ctx context.Context, accountID uuid.UUID, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", accountID, userID).
		Delete(&model.BankAccount{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete bank account: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("bank account not found")
	}

	return nil
}
