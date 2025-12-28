package repository

import (
	"bus-booking/payment-service/internal/model"
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefundRepository interface {
	// Core CRUD
	Create(ctx context.Context, refund *model.Refund) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Refund, error)
	GetByBookingID(ctx context.Context, bookingID uuid.UUID) (*model.Refund, error)
	Update(ctx context.Context, refund *model.Refund) error

	// List & Filter
	List(ctx context.Context, query *model.RefundListQuery) ([]*model.Refund, int64, error)
	ListByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.Refund, error)

	// Stats
	GetPendingRefundsStats(ctx context.Context) (totalAmount int, count int, err error)
	GetCompletedRefundsTotal(ctx context.Context) (int, error)
}

type RefundRepositoryImpl struct {
	db *gorm.DB
}

func NewRefundRepository(db *gorm.DB) RefundRepository {
	return &RefundRepositoryImpl{db: db}
}

// Create creates a new refund
func (r *RefundRepositoryImpl) Create(ctx context.Context, refund *model.Refund) error {
	if err := r.db.WithContext(ctx).Create(refund).Error; err != nil {
		return fmt.Errorf("failed to create refund: %w", err)
	}
	return nil
}

// GetByID retrieves a refund by ID
func (r *RefundRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.Refund, error) {
	var refund model.Refund
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&refund).Error; err != nil {
		return nil, fmt.Errorf("refund not found: %w", err)
	}
	return &refund, nil
}

// GetByBookingID retrieves a refund by booking ID
func (r *RefundRepositoryImpl) GetByBookingID(ctx context.Context, bookingID uuid.UUID) (*model.Refund, error) {
	var refund model.Refund
	if err := r.db.WithContext(ctx).
		Where("booking_id = ?", bookingID).
		First(&refund).Error; err != nil {
		return nil, fmt.Errorf("refund not found: %w", err)
	}
	return &refund, nil
}

// Update updates a refund
func (r *RefundRepositoryImpl) Update(ctx context.Context, refund *model.Refund) error {
	if err := r.db.WithContext(ctx).Save(refund).Error; err != nil {
		return fmt.Errorf("failed to update refund: %w", err)
	}
	return nil
}

// List retrieves refunds with filters and pagination
func (r *RefundRepositoryImpl) List(ctx context.Context, query *model.RefundListQuery) ([]*model.Refund, int64, error) {
	db := r.db.WithContext(ctx).Model(&model.Refund{})

	// Apply filters
	if query.Status != nil {
		db = db.Where("refund_status = ?", *query.Status)
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
		return nil, 0, fmt.Errorf("failed to count refunds: %w", err)
	}

	// Normalize pagination
	query.Normalize()

	// Calculate offset
	offset := (query.Page - 1) * query.PageSize

	// Get results with transaction preload
	var refunds []*model.Refund
	if err := db.
		Preload("Transaction").
		Offset(offset).
		Limit(query.PageSize).
		Order("created_at DESC").
		Find(&refunds).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list refunds: %w", err)
	}

	return refunds, total, nil
}

// ListByIDs retrieves multiple refunds by their IDs (for export)
func (r *RefundRepositoryImpl) ListByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.Refund, error) {
	var refunds []*model.Refund
	if err := r.db.WithContext(ctx).
		Preload("Transaction").
		Where("id IN ?", ids).
		Order("created_at DESC").
		Find(&refunds).Error; err != nil {
		return nil, fmt.Errorf("failed to get refunds by IDs: %w", err)
	}
	return refunds, nil
}

// GetPendingRefundsStats gets total amount and count of pending refunds
func (r *RefundRepositoryImpl) GetPendingRefundsStats(ctx context.Context) (totalAmount int, count int, err error) {
	// Get total amount
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Refund{}).
		Where("refund_status = ?", model.RefundStatusPending).
		Select("COALESCE(SUM(refund_amount), 0)").
		Scan(&total).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to sum pending refunds: %w", err)
	}

	// Get count
	var pendingCount int64
	if err := r.db.WithContext(ctx).Model(&model.Refund{}).
		Where("refund_status = ?", model.RefundStatusPending).
		Count(&pendingCount).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to count pending refunds: %w", err)
	}

	return int(total), int(pendingCount), nil
}

// GetCompletedRefundsTotal gets total amount of completed refunds
func (r *RefundRepositoryImpl) GetCompletedRefundsTotal(ctx context.Context) (int, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Refund{}).
		Where("refund_status = ?", model.RefundStatusCompleted).
		Select("COALESCE(SUM(refund_amount), 0)").
		Scan(&total).Error; err != nil {
		return 0, fmt.Errorf("failed to sum completed refunds: %w", err)
	}
	return int(total), nil
}
