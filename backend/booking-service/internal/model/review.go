package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReviewStatus string

const (
	ReviewStatusActive  ReviewStatus = "active"
	ReviewStatusHidden  ReviewStatus = "hidden"
	ReviewStatusFlagged ReviewStatus = "flagged"
	ReviewStatusRemoved ReviewStatus = "removed"
)

type Review struct {
	BaseModel
	// Trip reference
	TripID uuid.UUID `gorm:"type:uuid;not null;index" json:"trip_id" validate:"required"`

	// User & booking verification
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id" validate:"required"`
	BookingID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uq_reviews_booking" json:"booking_id" validate:"required"`

	// Review content
	Rating  int    `gorm:"type:int;not null" json:"rating" validate:"required,min=1,max=5"`
	Comment string `gorm:"type:text" json:"comment,omitempty"`

	// Moderation
	IsVerified bool         `gorm:"type:boolean;not null;default:true" json:"is_verified"`
	Status     ReviewStatus `gorm:"type:varchar(20);not null;default:'active';index" json:"status" validate:"required"`
	AdminNotes string       `gorm:"type:text" json:"admin_notes,omitempty"`
}

func (Review) TableName() string {
	return "reviews"
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	r.IsVerified = true // Always verified via booking
	if r.Status == "" {
		r.Status = ReviewStatusActive
	}
	return nil
}

// CreateReviewRequest for creating trip review
type CreateReviewRequest struct {
	BookingID uuid.UUID `json:"booking_id" validate:"required"`
	Rating    int       `json:"rating" validate:"required,min=1,max=5"`
	Comment   string    `json:"comment,omitempty" validate:"omitempty,max=1000"`
}

// UpdateReviewRequest for updating review
type UpdateReviewRequest struct {
	Rating  *int    `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
	Comment *string `json:"comment,omitempty" validate:"omitempty,max=1000"`
}

// ReviewResponse for API responses
type ReviewResponse struct {
	ID         uuid.UUID    `json:"id"`
	TripID     uuid.UUID    `json:"trip_id"`
	UserID     uuid.UUID    `json:"user_id"`
	BookingID  uuid.UUID    `json:"booking_id"`
	Rating     int          `json:"rating"`
	Comment    string       `json:"comment,omitempty"`
	IsVerified bool         `json:"is_verified"`
	Status     ReviewStatus `json:"status"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

// GetUserReviewsRequest for listing user reviews
type GetUserReviewsRequest struct {
	PaginationRequest
	UserID uuid.UUID     `json:"user_id" validate:"required"`
	Status *ReviewStatus `form:"status" json:"status,omitempty"`
}

// GetTripReviewsRequest for listing trip reviews
type GetTripReviewsRequest struct {
	PaginationRequest
	TripID    *uuid.UUID    `form:"trip_id" json:"trip_id,omitempty"`
	MinRating *int          `form:"min_rating" json:"min_rating,omitempty" validate:"omitempty,min=1,max=5"`
	Status    *ReviewStatus `form:"status" json:"status,omitempty"`
}

// TripReviewSummary for aggregated review stats
type TripReviewSummary struct {
	TripID        uuid.UUID `json:"trip_id"`
	TotalReviews  int       `json:"total_reviews"`
	AverageRating float64   `json:"average_rating"`
	Rating1Count  int       `json:"rating_1_count"`
	Rating2Count  int       `json:"rating_2_count"`
	Rating3Count  int       `json:"rating_3_count"`
	Rating4Count  int       `json:"rating_4_count"`
	Rating5Count  int       `json:"rating_5_count"`
}
