package repository

import (
	"bus-booking/booking-service/internal/model"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *model.Review) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Review, error)
	GetByBookingID(ctx context.Context, bookingID uuid.UUID) (*model.Review, error)
	GetUserReviews(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*model.Review, int64, error)
	GetTripReviews(ctx context.Context, req *model.GetTripReviewsRequest) ([]*model.Review, int64, error)
	GetTripReviewSummary(ctx context.Context, tripID uuid.UUID) (*model.TripReviewSummary, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ReviewRepositoryImpl struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &ReviewRepositoryImpl{db: db}
}

func (r *ReviewRepositoryImpl) Create(ctx context.Context, review *model.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *ReviewRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.Review, error) {
	var review model.Review
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepositoryImpl) GetByBookingID(ctx context.Context, bookingID uuid.UUID) (*model.Review, error) {
	var review model.Review
	err := r.db.WithContext(ctx).Where("booking_id = ?", bookingID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepositoryImpl) GetUserReviews(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*model.Review, int64, error) {
	var reviews []*model.Review
	var total int64

	query := r.db.WithContext(ctx).
		Model(&model.Review{}).
		Where("user_id = ?", userID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&reviews).Error

	return reviews, total, err
}

func (r *ReviewRepositoryImpl) GetTripReviews(ctx context.Context, req *model.GetTripReviewsRequest) ([]*model.Review, int64, error) {
	var reviews []*model.Review
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Review{})

	// Apply filters
	if req.TripID != nil {
		query = query.Where("trip_id = ?", *req.TripID)
	}
	if req.MinRating != nil {
		query = query.Where("rating >= ?", *req.MinRating)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	} else {
		// Default to active reviews only
		query = query.Where("status = ?", model.ReviewStatusActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (req.Page - 1) * req.PageSize
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&reviews).Error

	return reviews, total, err
}

func (r *ReviewRepositoryImpl) GetTripReviewSummary(ctx context.Context, tripID uuid.UUID) (*model.TripReviewSummary, error) {
	var summary model.TripReviewSummary
	summary.TripID = tripID

	// Get rating distribution
	type ratingCount struct {
		Rating int
		Count  int
	}
	var counts []ratingCount

	err := r.db.WithContext(ctx).
		Model(&model.Review{}).
		Select("rating, COUNT(*) as count").
		Where("trip_id = ?", tripID).
		Where("status = ?", model.ReviewStatusActive).
		Group("rating").
		Scan(&counts).Error

	if err != nil {
		return nil, err
	}

	// Populate summary
	totalRating := 0
	for _, rc := range counts {
		summary.TotalReviews += rc.Count
		totalRating += rc.Rating * rc.Count

		switch rc.Rating {
		case 1:
			summary.Rating1Count = rc.Count
		case 2:
			summary.Rating2Count = rc.Count
		case 3:
			summary.Rating3Count = rc.Count
		case 4:
			summary.Rating4Count = rc.Count
		case 5:
			summary.Rating5Count = rc.Count
		}
	}

	if summary.TotalReviews > 0 {
		summary.AverageRating = float64(totalRating) / float64(summary.TotalReviews)
	}

	return &summary, nil
}

func (r *ReviewRepositoryImpl) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&model.Review{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *ReviewRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.Review{}).Error
}
