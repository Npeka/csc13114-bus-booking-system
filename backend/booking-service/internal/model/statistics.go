package model

import (
	"time"

	"github.com/google/uuid"
)

type TripBookingStats struct {
	TripID        uuid.UUID `json:"trip_id"`
	TotalBookings int64     `json:"total_bookings"`
	TotalRevenue  float64   `json:"total_revenue"`
	AverageRating float64   `json:"average_rating"`
}

type BookingStats struct {
	TotalBookings     int64   `json:"total_bookings"`
	TotalRevenue      float64 `json:"total_revenue"`
	CancelledBookings int64   `json:"cancelled_bookings"`
	CompletedBookings int64   `json:"completed_bookings"`
	AverageRating     float64 `json:"average_rating"`
}

type BookingStatsRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

type BookingStatsResponse struct {
	TotalBookings     int64     `json:"total_bookings"`
	TotalRevenue      float64   `json:"total_revenue"`
	CancelledBookings int64     `json:"cancelled_bookings"`
	CompletedBookings int64     `json:"completed_bookings"`
	AverageRating     float64   `json:"average_rating"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
}

type PopularTripsRequest struct {
	Limit int `form:"limit,default=10"`
	Days  int `form:"days,default=30"`
}

type TripStatsResponse struct {
	TripID        uuid.UUID `json:"trip_id"`
	TotalBookings int64     `json:"total_bookings"`
	TotalRevenue  float64   `json:"total_revenue"`
	AverageRating float64   `json:"average_rating"`
}
