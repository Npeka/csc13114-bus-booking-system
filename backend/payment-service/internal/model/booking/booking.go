package booking

import (
	"bus-booking/payment-service/internal/model"
)

type UpdateBookingStatusRequest struct {
	TransactionStatus model.TransactionStatus `json:"transaction_status"`
}

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "PENDING"
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusCancelled BookingStatus = "CANCELLED"
	BookingStatusExpired   BookingStatus = "EXPIRED"
	BookingStatusFailed    BookingStatus = "FAILED"
)
