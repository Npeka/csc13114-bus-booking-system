package user

import (
	"bus-booking/shared/constants"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID            `json:"id"`
	Email         string               `json:"email,omitempty"`
	Phone         string               `json:"phone,omitempty"`
	FullName      string               `json:"full_name"`
	Avatar        string               `json:"avatar,omitempty"`
	Role          constants.UserRole   `json:"role"`
	Status        constants.UserStatus `json:"status"`
	EmailVerified bool                 `json:"email_verified"`
	PhoneVerified bool                 `json:"phone_verified"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}
