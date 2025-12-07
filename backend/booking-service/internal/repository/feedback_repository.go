package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"bus-booking/booking-service/internal/model"
)

type FeedbackRepository interface {
	CreateFeedback(ctx context.Context, feedback *model.Feedback) error
	GetFeedbackByBookingID(ctx context.Context, bookingID uuid.UUID) (*model.Feedback, error)
	GetFeedbacksByTripID(ctx context.Context, tripID uuid.UUID, limit, offset int) ([]*model.Feedback, int64, error)
	UpdateFeedback(ctx context.Context, feedback *model.Feedback) error
	DeleteFeedback(ctx context.Context, id uuid.UUID) error
}

type feedbackRepositoryImpl struct {
	db *gorm.DB
}

func NewFeedbackRepository(db *gorm.DB) FeedbackRepository {
	return &feedbackRepositoryImpl{db: db}
}

func (r *feedbackRepositoryImpl) CreateFeedback(ctx context.Context, feedback *model.Feedback) error {
	if err := r.db.WithContext(ctx).Create(feedback).Error; err != nil {
		return fmt.Errorf("failed to create feedback: %w", err)
	}
	return nil
}

func (r *feedbackRepositoryImpl) GetFeedbackByBookingID(ctx context.Context, bookingID uuid.UUID) (*model.Feedback, error) {
	var feedback model.Feedback
	err := r.db.WithContext(ctx).First(&feedback, "booking_id = ?", bookingID).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("feedback not found")
		}
		return nil, fmt.Errorf("failed to get feedback: %w", err)
	}

	return &feedback, nil
}

func (r *feedbackRepositoryImpl) GetFeedbacksByTripID(ctx context.Context, tripID uuid.UUID, limit, offset int) ([]*model.Feedback, int64, error) {
	var feedbacks []*model.Feedback
	var total int64

	if err := r.db.WithContext(ctx).
		Model(&model.Feedback{}).
		Where("trip_id = ?", tripID).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count feedbacks: %w", err)
	}

	err := r.db.WithContext(ctx).
		Where("trip_id = ?", tripID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&feedbacks).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get feedbacks: %w", err)
	}

	return feedbacks, total, nil
}

func (r *feedbackRepositoryImpl) UpdateFeedback(ctx context.Context, feedback *model.Feedback) error {
	if err := r.db.WithContext(ctx).Save(feedback).Error; err != nil {
		return fmt.Errorf("failed to update feedback: %w", err)
	}
	return nil
}

func (r *feedbackRepositoryImpl) DeleteFeedback(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&model.Feedback{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete feedback: %w", err)
	}
	return nil
}
