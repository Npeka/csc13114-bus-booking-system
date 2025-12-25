package service

import (
	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type ReviewService interface {
	CreateReview(ctx context.Context, userID uuid.UUID, req *model.CreateReviewRequest) (*model.Review, error)
	GetReview(ctx context.Context, id uuid.UUID) (*model.Review, error)
	GetReviewByBooking(ctx context.Context, bookingID uuid.UUID) (*model.Review, error)
	GetUserReviews(ctx context.Context, req *model.GetUserReviewsRequest) ([]*model.ReviewResponse, int64, error)
	GetTripReviews(ctx context.Context, req *model.GetTripReviewsRequest) ([]*model.ReviewResponse, int64, error)
	GetTripReviewSummary(ctx context.Context, tripID uuid.UUID) (*model.TripReviewSummary, error)
	UpdateReview(ctx context.Context, userID, reviewID uuid.UUID, req *model.UpdateReviewRequest) (*model.Review, error)
	DeleteReview(ctx context.Context, userID, reviewID uuid.UUID) error
	ModerateReview(ctx context.Context, reviewID uuid.UUID, status model.ReviewStatus, adminNotes string) error
}

type ReviewServiceImpl struct {
	reviewRepo  repository.ReviewRepository
	bookingRepo repository.BookingRepository
}

func NewReviewService(
	reviewRepo repository.ReviewRepository,
	bookingRepo repository.BookingRepository,
) ReviewService {
	return &ReviewServiceImpl{
		reviewRepo:  reviewRepo,
		bookingRepo: bookingRepo,
	}
}

func (s *ReviewServiceImpl) CreateReview(ctx context.Context, userID uuid.UUID, req *model.CreateReviewRequest) (*model.Review, error) {
	// Get booking to verify
	booking, err := s.bookingRepo.GetBookingByID(ctx, req.BookingID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("booking not found")
		}
		log.Error().Err(err).Str("booking_id", req.BookingID.String()).Msg("failed to get booking")
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	// Verify booking belongs to user
	if booking.UserID != userID {
		log.Warn().
			Str("user_id", userID.String()).
			Str("booking_id", req.BookingID.String()).
			Msg("unauthorized review attempt")
		return nil, fmt.Errorf("unauthorized: booking does not belong to user")
	}

	// Verify booking is confirmed
	if booking.Status != model.BookingStatusConfirmed {
		return nil, fmt.Errorf("can only review confirmed bookings, current status: %s", booking.Status)
	}

	// Check if review already exists for this booking
	existing, err := s.reviewRepo.GetByBookingID(ctx, req.BookingID)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error().Err(err).Msg("failed to check existing review")
		return nil, fmt.Errorf("failed to check existing review: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("review already exists for this booking")
	}

	// Create review using trip_id from booking
	review := &model.Review{
		TripID:     booking.TripID,
		UserID:     userID,
		BookingID:  req.BookingID,
		Rating:     req.Rating,
		Comment:    req.Comment,
		IsVerified: true,
		Status:     model.ReviewStatusActive,
	}

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		log.Error().Err(err).Msg("failed to create review")
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	log.Info().
		Str("review_id", review.ID.String()).
		Str("user_id", userID.String()).
		Str("booking_id", req.BookingID.String()).
		Int("rating", req.Rating).
		Msg("created review")

	return review, nil
}

func (s *ReviewServiceImpl) GetReview(ctx context.Context, id uuid.UUID) (*model.Review, error) {
	return s.reviewRepo.GetByID(ctx, id)
}

func (s *ReviewServiceImpl) GetReviewByBooking(ctx context.Context, bookingID uuid.UUID) (*model.Review, error) {
	return s.reviewRepo.GetByBookingID(ctx, bookingID)
}

func (s *ReviewServiceImpl) GetUserReviews(ctx context.Context, req *model.GetUserReviewsRequest) ([]*model.ReviewResponse, int64, error) {
	reviews, total, err := s.reviewRepo.GetUserReviews(ctx, req.UserID, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*model.ReviewResponse, len(reviews))
	for i, review := range reviews {
		responses[i] = s.toResponse(review)
	}

	return responses, total, nil
}

func (s *ReviewServiceImpl) GetTripReviews(ctx context.Context, req *model.GetTripReviewsRequest) ([]*model.ReviewResponse, int64, error) {
	reviews, total, err := s.reviewRepo.GetTripReviews(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*model.ReviewResponse, len(reviews))
	for i, review := range reviews {
		responses[i] = s.toResponse(review)
	}

	return responses, total, nil
}

func (s *ReviewServiceImpl) GetTripReviewSummary(ctx context.Context, tripID uuid.UUID) (*model.TripReviewSummary, error) {
	return s.reviewRepo.GetTripReviewSummary(ctx, tripID)
}

func (s *ReviewServiceImpl) UpdateReview(ctx context.Context, userID, reviewID uuid.UUID, req *model.UpdateReviewRequest) (*model.Review, error) {
	// Get existing review
	review, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("review not found")
		}
		return nil, fmt.Errorf("failed to get review: %w", err)
	}

	// Verify ownership
	if review.UserID != userID {
		return nil, fmt.Errorf("unauthorized: review does not belong to user")
	}

	// Build updates
	updates := make(map[string]interface{})
	if req.Rating != nil {
		updates["rating"] = *req.Rating
	}
	if req.Comment != nil {
		updates["comment"] = *req.Comment
	}

	if len(updates) > 0 {
		if err := s.reviewRepo.Update(ctx, reviewID, updates); err != nil {
			log.Error().Err(err).Str("review_id", reviewID.String()).Msg("failed to update review")
			return nil, fmt.Errorf("failed to update review: %w", err)
		}
	}

	return s.reviewRepo.GetByID(ctx, reviewID)
}

func (s *ReviewServiceImpl) DeleteReview(ctx context.Context, userID, reviewID uuid.UUID) error {
	// Get review to verify ownership
	review, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("review not found")
		}
		return fmt.Errorf("failed to get review: %w", err)
	}

	// Verify ownership
	if review.UserID != userID {
		return fmt.Errorf("unauthorized: review does not belong to user")
	}

	return s.reviewRepo.Delete(ctx, reviewID)
}

func (s *ReviewServiceImpl) ModerateReview(ctx context.Context, reviewID uuid.UUID, status model.ReviewStatus, adminNotes string) error {
	updates := map[string]interface{}{
		"status":      status,
		"admin_notes": adminNotes,
	}

	if err := s.reviewRepo.Update(ctx, reviewID, updates); err != nil {
		log.Error().Err(err).Str("review_id", reviewID.String()).Msg("failed to moderate review")
		return fmt.Errorf("failed to moderate review: %w", err)
	}

	log.Info().
		Str("review_id", reviewID.String()).
		Str("status", string(status)).
		Msg("moderated review")

	return nil
}

// Helper function
func (s *ReviewServiceImpl) toResponse(review *model.Review) *model.ReviewResponse {
	return &model.ReviewResponse{
		ID:         review.ID,
		TripID:     review.TripID,
		UserID:     review.UserID,
		BookingID:  review.BookingID,
		Rating:     review.Rating,
		Comment:    review.Comment,
		IsVerified: review.IsVerified,
		Status:     review.Status,
		CreatedAt:  review.CreatedAt,
		UpdatedAt:  review.UpdatedAt,
	}
}
