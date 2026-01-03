package booking

import "github.com/google/uuid"

type Booking struct {
	ID                uuid.UUID `json:"id"`
	TotalAmount       int       `json:"total_amount"`
	TransactionStatus string    `json:"transaction_status"`
	Status            string    `json:"status"`
}

type CancelBookingRequest struct {
	Reason string `json:"reason"`
}

// BookingListResponse represents the response when getting bookings for a trip
type BookingListResponse struct {
	Data  []*Booking `json:"data"`
	Page  int        `json:"page"`
	Total int64      `json:"total"`
}
