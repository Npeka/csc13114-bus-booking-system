package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BookingStatus represents booking status enum
type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusRefunded  BookingStatus = "refunded"
)

// SeatState represents seat booking state
type SeatState string

const (
	SeatStateAvailable SeatState = "available"
	SeatStateBooked    SeatState = "booked"
	SeatStateLocked    SeatState = "locked"
)

// Booking represents a bus ticket booking
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

	// Relationships
	BookingSeats  []BookingSeat `gorm:"foreignKey:BookingID" json:"booking_seats,omitempty"`
	PaymentMethod PaymentMethod `gorm:"foreignKey:PaymentMethodID" json:"payment_method,omitempty"`
}

// BeforeCreate sets UUID for new bookings
func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// BookingSeat represents the relationship between booking and seats
type BookingSeat struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BookingID uuid.UUID `gorm:"type:uuid;not null" json:"booking_id" validate:"required"`
	SeatID    uuid.UUID `gorm:"type:uuid;not null" json:"seat_id" validate:"required"`
	Price     float64   `gorm:"type:decimal(10,2);not null" json:"price" validate:"required,min=0"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`

	// Relationships
	Booking *Booking `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"booking,omitempty"`
}

// BeforeCreate sets UUID for new booking seats
func (bs *BookingSeat) BeforeCreate(tx *gorm.DB) error {
	if bs.ID == uuid.Nil {
		bs.ID = uuid.New()
	}
	return nil
}

// SeatStatus represents seat availability status for specific trips
type SeatStatus struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TripID        uuid.UUID  `gorm:"type:uuid;not null" json:"trip_id" validate:"required"`
	SeatID        uuid.UUID  `gorm:"type:uuid;not null" json:"seat_id" validate:"required"`
	SeatNumber    string     `gorm:"type:varchar(10);not null" json:"seat_number" validate:"required"`
	Status        string     `gorm:"type:varchar(50);not null;default:'available'" json:"status"`
	UserID        *uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	BookingID     *uuid.UUID `gorm:"type:uuid" json:"booking_id,omitempty"`
	ReservedUntil *time.Time `gorm:"type:timestamptz" json:"reserved_until,omitempty"`
	CreatedAt     time.Time  `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`

	// Relationships
	Booking *Booking `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"booking,omitempty"`
}

// BeforeCreate sets UUID for new seat status
func (ss *SeatStatus) BeforeCreate(tx *gorm.DB) error {
	if ss.ID == uuid.Nil {
		ss.ID = uuid.New()
	}
	return nil
}

// PaymentMethod represents available payment methods
type PaymentMethod struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name" validate:"required"`
	Code        string         `gorm:"type:varchar(50);not null;unique" json:"code" validate:"required"`
	Description string         `gorm:"type:text" json:"description"`
	IsActive    bool           `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedAt   time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate sets UUID for new payment methods
func (pm *PaymentMethod) BeforeCreate(tx *gorm.DB) error {
	if pm.ID == uuid.Nil {
		pm.ID = uuid.New()
	}
	return nil
}

// Feedback represents user feedback for completed trips
type Feedback struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BookingID   uuid.UUID      `gorm:"type:uuid;not null" json:"booking_id" validate:"required"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null" json:"user_id" validate:"required"`
	TripID      uuid.UUID      `gorm:"type:uuid;not null" json:"trip_id" validate:"required"`
	Rating      int            `gorm:"type:integer;not null" json:"rating" validate:"required,min=1,max=5"`
	Comment     string         `gorm:"type:text" json:"comment"`
	SubmittedAt time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"submitted_at"`
	CreatedAt   time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Booking *Booking `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"booking,omitempty"`
}

// BeforeCreate sets UUID for new feedback
func (f *Feedback) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// BookingStats represents booking statistics
type BookingStats struct {
	TotalBookings     int64   `json:"total_bookings"`
	TotalRevenue      float64 `json:"total_revenue"`
	CancelledBookings int64   `json:"cancelled_bookings"`
	CompletedBookings int64   `json:"completed_bookings"`
	AverageRating     float64 `json:"average_rating"`
}

// TripBookingStats represents trip booking statistics
type TripBookingStats struct {
	TripID        uuid.UUID `json:"trip_id"`
	TotalBookings int64     `json:"total_bookings"`
	TotalRevenue  float64   `json:"total_revenue"`
	AverageRating float64   `json:"average_rating"`
}

// TableName overrides to ensure consistent naming
func (Booking) TableName() string       { return "bookings" }
func (BookingSeat) TableName() string   { return "booking_seats" }
func (SeatStatus) TableName() string    { return "seat_statuses" }
func (PaymentMethod) TableName() string { return "payment_methods" }
func (Feedback) TableName() string      { return "feedbacks" }
