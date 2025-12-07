package service

import (
	"context"
	"time"

	"bus-booking/shared/ginext"

	"github.com/google/uuid"

	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/repository"
)

type FeedbackService interface {
	CreateFeedback(ctx context.Context, req *model.CreateFeedbackRequest) (*model.FeedbackResponse, error)
	GetBookingFeedback(ctx context.Context, bookingID uuid.UUID) (*model.FeedbackResponse, error)
	GetTripFeedbacks(ctx context.Context, tripID uuid.UUID, page, limit int) ([]*model.FeedbackResponse, int64, error)
}

type FeedbackServiceImpl struct {
	bookingRepo  repository.BookingRepository
	feedbackRepo repository.FeedbackRepository
}

func NewFeedbackService(bookingRepo repository.BookingRepository, feedbackRepo repository.FeedbackRepository) FeedbackService {
	return &FeedbackServiceImpl{
		bookingRepo:  bookingRepo,
		feedbackRepo: feedbackRepo,
	}
}

// CreateFeedback creates feedback for a booking
func (s *FeedbackServiceImpl) CreateFeedback(ctx context.Context, req *model.CreateFeedbackRequest) (*model.FeedbackResponse, error) {
	// Verify booking exists and belongs to user
	booking, err := s.bookingRepo.GetBookingByID(ctx, req.BookingID)
	if err != nil {
		return nil, err
	}

	if booking.UserID != req.UserID {
		return nil, ginext.NewForbiddenError("booking does not belong to user")
	}

	if booking.Status != "completed" {
		return nil, ginext.NewBadRequestError("can only provide feedback for completed bookings")
	}

	// Check if feedback already exists
	if _, err := s.feedbackRepo.GetFeedbackByBookingID(ctx, req.BookingID); err == nil {
		return nil, ginext.NewConflictError("feedback already exists for this booking")
	}

	feedback := &model.Feedback{
		BaseModel: model.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		UserID:    req.UserID,
		BookingID: req.BookingID,
		TripID:    booking.TripID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}

	if err := s.feedbackRepo.CreateFeedback(ctx, feedback); err != nil {
		return nil, err
	}

	return &model.FeedbackResponse{
		ID:        feedback.ID,
		UserID:    feedback.UserID,
		BookingID: feedback.BookingID,
		TripID:    feedback.TripID,
		Rating:    feedback.Rating,
		Comment:   feedback.Comment,
		CreatedAt: feedback.CreatedAt,
	}, nil
}

// GetBookingFeedback retrieves feedback for a booking
func (s *FeedbackServiceImpl) GetBookingFeedback(ctx context.Context, bookingID uuid.UUID) (*model.FeedbackResponse, error) {
	feedback, err := s.feedbackRepo.GetFeedbackByBookingID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	return &model.FeedbackResponse{
		ID:        feedback.ID,
		UserID:    feedback.UserID,
		BookingID: feedback.BookingID,
		TripID:    feedback.TripID,
		Rating:    feedback.Rating,
		Comment:   feedback.Comment,
		CreatedAt: feedback.CreatedAt,
	}, nil
}

// GetTripFeedbacks retrieves feedbacks for a trip with pagination
func (s *FeedbackServiceImpl) GetTripFeedbacks(ctx context.Context, tripID uuid.UUID, page, limit int) ([]*model.FeedbackResponse, int64, error) {
	offset := (page - 1) * limit
	feedbacks, total, err := s.feedbackRepo.GetFeedbacksByTripID(ctx, tripID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var feedbackResponses []*model.FeedbackResponse
	for _, feedback := range feedbacks {
		response := &model.FeedbackResponse{
			ID:        feedback.ID,
			UserID:    feedback.UserID,
			BookingID: feedback.BookingID,
			TripID:    feedback.TripID,
			Rating:    feedback.Rating,
			Comment:   feedback.Comment,
			CreatedAt: feedback.CreatedAt,
		}
		feedbackResponses = append(feedbackResponses, response)
	}

	return feedbackResponses, total, nil
}
