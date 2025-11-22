package repository

import (
	"bus-booking/payment-service/internal/model"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *model.Transaction) error
}

type transactionRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepositoryImpl{db: db}
}

func (r *transactionRepositoryImpl) CreateTransaction(ctx context.Context, transaction *model.Transaction) error {
	if err := r.db.WithContext(ctx).Create(transaction).Error; err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}
