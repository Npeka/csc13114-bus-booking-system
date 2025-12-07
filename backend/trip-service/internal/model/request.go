package model

import (
	"github.com/google/uuid"
)

// ListBusesRequest represents query parameters for listing buses
type ListBusesRequest struct {
	PaginationRequest
}

// ListRoutesRequest represents query parameters for listing routes
type ListRoutesRequest struct {
	PaginationRequest
}

// SearchRoutesQueryRequest represents query parameters for searching routes
type SearchRoutesQueryRequest struct {
	Origin      string `form:"origin" binding:"required"`
	Destination string `form:"destination" binding:"required"`
}

// ListTripsByRouteRequest represents query parameters for listing trips by route
type ListTripsByRouteRequest struct {
	Date string `form:"date" binding:"required"`
}

// SeatAvailabilityRequest represents seat availability check request
type SeatAvailabilityRequest struct {
	TripID uuid.UUID `json:"trip_id" validate:"required"`
}

type CreateBookingRequest struct {
	TripID     uuid.UUID          `json:"trip_id" validate:"required"`
	SeatIDs    []uuid.UUID        `json:"seat_ids" validate:"required,min=1,max=10"`
	Passengers []PassengerRequest `json:"passengers" validate:"required,dive"`
	IsGuest    bool               `json:"is_guest"`
	GuestInfo  *GuestInfo         `json:"guest_info,omitempty"`
}

type PassengerRequest struct {
	SeatID      uuid.UUID `json:"seat_id" validate:"required"`
	FullName    string    `json:"full_name" validate:"required,min=2,max=255"`
	IDNumber    string    `json:"id_number,omitempty" validate:"omitempty,min=5,max=50"`
	PhoneNumber string    `json:"phone_number,omitempty" validate:"omitempty,min=10,max=20"`
}

type GuestInfo struct {
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required,min=10,max=20"`
	Name  string `json:"name" validate:"required,min=2,max=255"`
}

type UnlockSeatsRequest struct {
	SessionID string `json:"session_id" validate:"required"`
}

type CancelBookingRequest struct {
	Reason string `json:"reason,omitempty"`
}

// Booking query parameters
type BookingLookupRequest struct {
	Reference string `form:"reference" binding:"required"`
	Email     string `form:"email" binding:"required,email"`
}

type ListBookingsRequest struct {
	PaginationRequest
	Status string `form:"status"`
}
