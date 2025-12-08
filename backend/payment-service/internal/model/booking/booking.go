package booking

import (
	"bus-booking/payment-service/internal/model"

	"github.com/google/uuid"
)

type UpdatePaymentStatusRequest struct {
	PaymentStatus  model.TransactionStatus `json:"payment_status"`
	BookingStatus  BookingStatus           `json:"booking_status"`
	TransactionID  uuid.UUID               `json:"transaction_id"`
	PaymentOrderID string                  `json:"payment_order_id"`
}

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "PENDING"
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusCancelled BookingStatus = "CANCELLED"
	BookingStatusExpired   BookingStatus = "EXPIRED"
	BookingStatusFailed    BookingStatus = "FAILED"
)
