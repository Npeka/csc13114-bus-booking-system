package repository

import (
	"bus-booking/payment-service/internal/model"
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	GetList(ctx context.Context, query *model.TransactionListQuery) ([]*model.Transaction, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Transaction, error)
	GetByBookingID(ctx context.Context, bookingID uuid.UUID) (*model.Transaction, error)
	GetByWebhookData(ctx context.Context, orderCode int, paymentLinkID string) (*model.Transaction, error)
	GetStats(ctx context.Context) (*model.TransactionStats, error)
	CreateTransaction(ctx context.Context, transaction *model.Transaction) error
	UpdateTransaction(ctx context.Context, transaction *model.Transaction) error
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

// GetList lists all transactions with filters (admin)
func (r *transactionRepositoryImpl) GetList(ctx context.Context, query *model.TransactionListQuery) ([]*model.Transaction, int64, error) {
	db := r.db.WithContext(ctx).Model(&model.Transaction{})

	// Apply filters
	if query.TransactionType != nil {
		db = db.Where("transaction_type = ?", *query.TransactionType)
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}
	if query.RefundStatus != nil {
		db = db.Where("refund_status = ?", *query.RefundStatus)
	}
	if query.StartDate != nil {
		db = db.Where("created_at >= ?", *query.StartDate)
	}
	if query.EndDate != nil {
		db = db.Where("created_at <= ?", *query.EndDate)
	}

	// Get total count
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	// Normalize pagination
	query.Normalize()

	// Calculate offset
	offset := (query.Page - 1) * query.PageSize

	// Get results
	var transactions []*model.Transaction
	if err := db.Offset(offset).Limit(query.PageSize).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list transactions: %w", err)
	}

	return transactions, total, nil
}

func (r *transactionRepositoryImpl) GetStats(ctx context.Context) (*model.TransactionStats, error) {
	stats := &model.TransactionStats{}

	// Total transactions
	var totalTx int64
	if err := r.db.WithContext(ctx).Model(&model.Transaction{}).Count(&totalTx).Error; err != nil {
		return nil, fmt.Errorf("failed to count total transactions: %w", err)
	}
	stats.TotalTransactions = int(totalTx)

	// Total IN (revenue)
	var totalIn int64
	if err := r.db.WithContext(ctx).Model(&model.Transaction{}).
		Where("transaction_type = ? AND status = ?", model.TransactionTypeIn, model.TransactionStatusPaid).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIn).Error; err != nil {
		return nil, fmt.Errorf("failed to sum IN transactions: %w", err)
	}
	stats.TotalIn = int(totalIn)

	// Total OUT (refunds completed)
	var totalOut int64
	if err := r.db.WithContext(ctx).Model(&model.Transaction{}).
		Where("transaction_type = ? AND refund_status = ?", model.TransactionTypeOut, model.RefundStatusCompleted).
		Select("COALESCE(SUM(refund_amount), 0)").
		Scan(&totalOut).Error; err != nil {
		return nil, fmt.Errorf("failed to sum OUT transactions: %w", err)
	}
	stats.TotalOut = int(totalOut)

	// Pending refunds amount
	var pendingRefunds int64
	if err := r.db.WithContext(ctx).Model(&model.Transaction{}).
		Where("transaction_type = ? AND refund_status = ?", model.TransactionTypeOut, model.RefundStatusPending).
		Select("COALESCE(SUM(refund_amount), 0)").
		Scan(&pendingRefunds).Error; err != nil {
		return nil, fmt.Errorf("failed to sum pending refunds: %w", err)
	}
	stats.PendingRefunds = int(pendingRefunds)

	// Pending refunds count
	var pendingCount int64
	if err := r.db.WithContext(ctx).Model(&model.Transaction{}).
		Where("transaction_type = ? AND refund_status = ?", model.TransactionTypeOut, model.RefundStatusPending).
		Count(&pendingCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count pending refunds: %w", err)
	}
	stats.PendingRefundCount = int(pendingCount)

	return stats, nil
}
