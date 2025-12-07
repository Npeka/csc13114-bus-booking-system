package booking

import "github.com/google/uuid"

type SeatStatus struct {
	SeatID   uuid.UUID `json:"seat_id"`
	IsBooked bool      `json:"is_booked"`
	IsLocked bool      `json:"is_locked"`
}

type SeatStatusResponse struct {
	IsBooked bool `json:"is_booked"`
	IsLocked bool `json:"is_locked"`
}
