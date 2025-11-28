package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/repository"
)

type FeedbackService interface {
	CreateFeedback(ctx context.Context, req *model.CreateFeedbackRequest) (*model.FeedbackResponse, error)
	GetBookingFeedback(ctx context.Context, bookingID uuid.UUID) (*model.FeedbackResponse, error)
	GetTripFeedbacks(ctx context.Context, tripID uuid.UUID, page, limit int) (*model.PaginatedFeedbackResponse, error)
}

type FeedbackServiceImpl struct {
	repositories *repository.Repositories
}

func NewFeedbackService(repositories *repository.Repositories) FeedbackService {
	return &FeedbackServiceImpl{
		repositories: repositories,
	}
}

// CreateFeedback creates feedback for a booking
func (s *FeedbackServiceImpl) CreateFeedback(ctx context.Context, req *model.CreateFeedbackRequest) (*model.FeedbackResponse, error) {
	// Verify booking exists and belongs to user
	booking, err := s.repositories.Booking.GetBookingByID(ctx, req.BookingID)
	if err != nil {
		return nil, err
	}

	if booking.UserID == nil || *booking.UserID != req.UserID {
		return nil, fmt.Errorf("booking does not belong to user")
	}

	if booking.Status != "completed" {
		return nil, fmt.Errorf("can only provide feedback for completed bookings")
	}

	// Check if feedback already exists
	if _, err := s.repositories.Feedback.GetFeedbackByBookingID(ctx, req.BookingID); err == nil {
		return nil, fmt.Errorf("feedback already exists for this booking")
	}

	feedback := &model.Feedback{
		ID:        uuid.New(),
		UserID:    req.UserID,
		BookingID: req.BookingID,
		TripID:    booking.TripID,
		Rating:    req.Rating,
		Comment:   req.Comment,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.repositories.Feedback.CreateFeedback(ctx, feedback); err != nil {
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
	feedback, err := s.repositories.Feedback.GetFeedbackByBookingID(ctx, bookingID)
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
func (s *FeedbackServiceImpl) GetTripFeedbacks(ctx context.Context, tripID uuid.UUID, page, limit int) (*model.PaginatedFeedbackResponse, error) {
	offset := (page - 1) * limit
	feedbacks, total, err := s.repositories.Feedback.GetFeedbacksByTripID(ctx, tripID, limit, offset)
	if err != nil {
		return nil, err
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

	return &model.PaginatedFeedbackResponse{
		Data:       feedbackResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: (total + int64(limit) - 1) / int64(limit),
	}, nil
}
