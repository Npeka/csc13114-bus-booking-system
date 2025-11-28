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
	GetTransactionByID(ctx context.Context, id uuid.UUID) (*model.Transaction, error)
	GetTransactionByOrderCode(ctx context.Context, orderCode int64) (*model.Transaction, error)
	GetTransactionByBookingID(ctx context.Context, bookingID uuid.UUID) (*model.Transaction, error)
	GetTransactionsByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*model.Transaction, error)
}

type transactionRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepositoryImpl{db: db}
}

// CreateTransaction creates a new transaction
func (r *transactionRepositoryImpl) CreateTransaction(ctx context.Context, transaction *model.Transaction) error {
	if err := r.db.WithContext(ctx).Create(transaction).Error; err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}

// UpdateTransaction updates an existing transaction
func (r *transactionRepositoryImpl) UpdateTransaction(ctx context.Context, transaction *model.Transaction) error {
	if err := r.db.WithContext(ctx).Save(transaction).Error; err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}
	return nil
}

// GetTransactionByID retrieves a transaction by ID
func (r *transactionRepositoryImpl) GetTransactionByID(ctx context.Context, id uuid.UUID) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&transaction).Error; err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}
	return &transaction, nil
}

// GetTransactionByOrderCode retrieves a transaction by PayOS order code
func (r *transactionRepositoryImpl) GetTransactionByOrderCode(ctx context.Context, orderCode int64) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := r.db.WithContext(ctx).Where("order_code = ?", orderCode).First(&transaction).Error; err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}
	return &transaction, nil
}

// GetTransactionByBookingID retrieves the latest transaction for a booking
func (r *transactionRepositoryImpl) GetTransactionByBookingID(ctx context.Context, bookingID uuid.UUID) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := r.db.WithContext(ctx).
		Where("booking_id = ?", bookingID).
		Order("created_at DESC").
		First(&transaction).Error; err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}
	return &transaction, nil
}

// GetTransactionsByBookingID retrieves all transactions for a booking
func (r *transactionRepositoryImpl) GetTransactionsByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	if err := r.db.WithContext(ctx).
		Where("booking_id = ?", bookingID).
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	return transactions, nil
}
