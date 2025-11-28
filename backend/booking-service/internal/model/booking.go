package model

import (
	"time"

	"github.com/google/uuid"
)

type BookingStatus string
type PaymentStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusExpired   BookingStatus = "expired"

	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusRefunded PaymentStatus = "refunded"
	PaymentStatusFailed   PaymentStatus = "failed"
)

type Booking struct {
	BaseModel
	BookingReference   string        `json:"booking_reference" gorm:"type:varchar(10);unique;not null"`
	TripID             uuid.UUID     `json:"trip_id" gorm:"type:uuid;not null"`
	UserID             *uuid.UUID    `json:"user_id,omitempty" gorm:"type:uuid"`
	GuestEmail         string        `json:"guest_email,omitempty" gorm:"type:varchar(255)"`
	GuestPhone         string        `json:"guest_phone,omitempty" gorm:"type:varchar(20)"`
	GuestName          string        `json:"guest_name,omitempty" gorm:"type:varchar(255)"`
	TotalAmount        float64       `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	Status             BookingStatus `json:"status" gorm:"type:varchar(20);not null"`
	PaymentStatus      PaymentStatus `json:"payment_status" gorm:"type:varchar(20);not null"`
	PaymentMethod      string        `json:"payment_method,omitempty" gorm:"type:varchar(50)"`
	PaymentID          string        `json:"payment_id,omitempty" gorm:"type:varchar(255)"`
	ExpiresAt          *time.Time    `json:"expires_at,omitempty" gorm:"type:timestamptz"`
	ConfirmedAt        *time.Time    `json:"confirmed_at,omitempty" gorm:"type:timestamptz"`
	CancelledAt        *time.Time    `json:"cancelled_at,omitempty" gorm:"type:timestamptz"`
	CancellationReason string        `json:"cancellation_reason,omitempty" gorm:"type:text"`

	Passengers []Passenger `json:"passengers,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Booking) TableName() string { return "bookings" }
