package repository

import (
	"bus-booking/payment-service/internal/model"
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *model.Transaction) error
	UpdateTransaction(ctx context.Context, transaction *model.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Transaction, error)
	GetByBookingID(ctx context.Context, bookingID uuid.UUID) (*model.Transaction, error)
	GetByWebhookData(ctx context.Context, orderCode int, paymentLinkID string) (*model.Transaction, error)
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

func (r *transactionRepositoryImpl) UpdateTransaction(ctx context.Context, transaction *model.Transaction) error {
	if err := r.db.WithContext(ctx).Save(transaction).Error; err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}
	return nil
}

func (r *transactionRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&transaction).Error; err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}
	return &transaction, nil
}

func (r *transactionRepositoryImpl) GetByWebhookData(ctx context.Context, orderCode int, paymentLinkID string) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := r.db.WithContext(ctx).
		Where("order_code = ? and payment_link_id = ?", orderCode, paymentLinkID).
		First(&transaction).Error; err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}
	return &transaction, nil
}

func (r *transactionRepositoryImpl) GetByBookingID(ctx context.Context, bookingID uuid.UUID) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := r.db.WithContext(ctx).
		Where("booking_id = ?", bookingID).
		Order("created_at DESC").
		First(&transaction).Error; err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}
	return &transaction, nil
}
