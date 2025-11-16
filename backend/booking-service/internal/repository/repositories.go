package repository

import (
	"gorm.io/gorm"
)

// Repositories holds all repository instances
type Repositories struct {
	Booking       BookingRepository
	PaymentMethod PaymentMethodRepository
	Feedback      FeedbackRepository
	BookingStats  BookingStatsRepository
	SeatStatus    SeatStatusRepository
}

// NewRepositories creates a new repositories instance with all implementations
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Booking:       NewBookingRepository(db),
		PaymentMethod: NewPaymentMethodRepository(db),
		Feedback:      NewFeedbackRepository(db),
		BookingStats:  NewBookingStatsRepository(db),
		SeatStatus:    NewSeatStatusRepository(db),
	}
}
