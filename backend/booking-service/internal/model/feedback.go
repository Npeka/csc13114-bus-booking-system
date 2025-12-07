package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Feedback struct {
	BaseModel
	BookingID   uuid.UUID `gorm:"type:uuid;not null" json:"booking_id" validate:"required"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id" validate:"required"`
	TripID      uuid.UUID `gorm:"type:uuid;not null" json:"trip_id" validate:"required"`
	Rating      int       `gorm:"type:integer;not null" json:"rating" validate:"required,min=1,max=5"`
	Comment     string    `gorm:"type:text" json:"comment"`
	SubmittedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"submitted_at"`

	Booking *Booking `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"booking,omitempty"`
}

func (Feedback) TableName() string { return "feedbacks" }

func (f *Feedback) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

type CreateFeedbackRequest struct {
	BookingID uuid.UUID `json:"booking_id" validate:"required"`
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	Rating    int       `json:"rating" validate:"required,min=1,max=5"`
	Comment   string    `json:"comment,omitempty"`
}

type GetTripFeedbacksRequest struct {
	PaginationRequest
	TripID uuid.UUID `json:"trip_id" validate:"required"`
}
