package model

import "github.com/google/uuid"

// PassengerResponse represents a passenger in a trip
type PassengerResponse struct {
	UserID           uuid.UUID `json:"user_id"`
	FullName         string    `json:"full_name"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	BookingID        uuid.UUID `json:"booking_id"`
	BookingReference string    `json:"booking_reference"`
	Status           string    `json:"status"`
	Seats            []string  `json:"seats"` // e.g. ["A1", "B2"]
	OriginalPrice    int       `json:"original_price"`
	PaidPrice        int       `json:"paid_price"`
	IsBoarded        bool      `json:"is_boarded"`
}
