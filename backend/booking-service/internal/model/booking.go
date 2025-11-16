package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusRefunded  BookingStatus = "refunded"
)

type SeatState string

const (
	SeatStateAvailable SeatState = "available"
	SeatStateBooked    SeatState = "booked"
	SeatStateLocked    SeatState = "locked"
)

type Booking struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID             uuid.UUID      `gorm:"type:uuid;not null" json:"user_id" validate:"required"`
	TripID             uuid.UUID      `gorm:"type:uuid;not null" json:"trip_id" validate:"required"`
	PaymentMethodID    uuid.UUID      `gorm:"type:uuid;not null" json:"payment_method_id" validate:"required"`
	Status             string         `gorm:"type:varchar(50);not null;default:'pending'" json:"status"`
	TotalAmount        float64        `gorm:"type:decimal(10,2);not null" json:"total_amount" validate:"required,min=0"`
	PassengerName      string         `gorm:"type:varchar(255);not null" json:"passenger_name" validate:"required"`
	PassengerPhone     string         `gorm:"type:varchar(20);not null" json:"passenger_phone" validate:"required"`
	PassengerEmail     string         `gorm:"type:varchar(255)" json:"passenger_email"`
	SpecialRequests    string         `gorm:"type:text" json:"special_requests"`
	CancellationReason string         `gorm:"type:text" json:"cancellation_reason"`
	CreatedAt          time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
	CompletedAt        *time.Time     `gorm:"type:timestamptz" json:"completed_at"`
	CancelledAt        *time.Time     `gorm:"type:timestamptz" json:"cancelled_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	BookingSeats  []BookingSeat `gorm:"foreignKey:BookingID" json:"booking_seats,omitempty"`
	PaymentMethod PaymentMethod `gorm:"foreignKey:PaymentMethodID" json:"payment_method,omitempty"`
}

func (Booking) TableName() string { return "bookings" }

func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

type BookingStats struct {
	TotalBookings     int64   `json:"total_bookings"`
	TotalRevenue      float64 `json:"total_revenue"`
	CancelledBookings int64   `json:"cancelled_bookings"`
	CompletedBookings int64   `json:"completed_bookings"`
	AverageRating     float64 `json:"average_rating"`
}

type TripBookingStats struct {
	TripID        uuid.UUID `json:"trip_id"`
	TotalBookings int64     `json:"total_bookings"`
	TotalRevenue  float64   `json:"total_revenue"`
	AverageRating float64   `json:"average_rating"`
}
