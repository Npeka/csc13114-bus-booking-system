package model

import (
	"time"

	"github.com/google/uuid"
)

type SeatLock struct {
	BaseModel
	TripID    uuid.UUID `json:"trip_id" gorm:"type:uuid;not null"`
	SeatID    uuid.UUID `json:"seat_id" gorm:"type:uuid;not null"`
	SessionID string    `json:"session_id" gorm:"type:varchar(255);not null"`
	LockedAt  time.Time `json:"locked_at" gorm:"type:timestamptz"`
	ExpiresAt time.Time `json:"expires_at" gorm:"type:timestamptz;not null"`
}

func (SeatLock) TableName() string { return "seat_locks" }
