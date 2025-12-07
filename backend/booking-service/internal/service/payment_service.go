package service

import (
	"context"
	"time"

	"bus-booking/shared/ginext"

	"github.com/google/uuid"

	"bus-booking/booking-service/internal/model"
	"bus-booking/booking-service/internal/repository"
)

type PaymentService interface {
	GetPaymentMethods(ctx context.Context) ([]*model.PaymentMethodResponse, error)
	ProcessPayment(ctx context.Context, req *model.ProcessPaymentRequest) (*model.PaymentResponse, error)
}

type PaymentServiceImpl struct {
	paymentMethodRepo repository.PaymentMethodRepository
	bookingRepo       repository.BookingRepository
}

func NewPaymentService(paymentMethodRepo repository.PaymentMethodRepository, bookingRepo repository.BookingRepository) PaymentService {
	return &PaymentServiceImpl{
		paymentMethodRepo: paymentMethodRepo,
		bookingRepo:       bookingRepo,
	}
}

// GetPaymentMethods retrieves all available payment methods
func (s *PaymentServiceImpl) GetPaymentMethods(ctx context.Context) ([]*model.PaymentMethodResponse, error) {
	paymentMethods, err := s.paymentMethodRepo.GetPaymentMethods(ctx)
	if err != nil {
		return nil, err
	}

	var responses []*model.PaymentMethodResponse
	for _, pm := range paymentMethods {
		response := &model.PaymentMethodResponse{
			ID:          pm.ID,
			Name:        pm.Name,
			Code:        pm.Code,
			Description: pm.Description,
			IsActive:    pm.IsActive,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// ProcessPayment processes payment for a booking
func (s *PaymentServiceImpl) ProcessPayment(ctx context.Context, req *model.ProcessPaymentRequest) (*model.PaymentResponse, error) {
	// Get booking
	booking, err := s.bookingRepo.GetBookingByID(ctx, req.BookingID)
	if err != nil {
		return nil, err
	}

	if booking.Status != "pending" {
		return nil, ginext.NewBadRequestError("booking is not in pending status")
	}

	// In a real implementation, you would integrate with payment gateway here
	// For now, we'll simulate payment processing

	response := &model.PaymentResponse{
		BookingID:       req.BookingID,
		Amount:          booking.TotalAmount,
		PaymentMethodID: uuid.Nil, // PaymentMethodID removed from new Booking model
		Status:          "completed",
		TransactionID:   uuid.New().String(),
		ProcessedAt:     time.Now().UTC(),
	}

	// Update booking status to confirmed would be done here
	// This is handled by booking service UpdateBookingStatus method

	return response, nil
}
