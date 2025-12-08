package user

import "github.com/google/uuid"

type CreateGuestRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

type GuestResponse struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email,omitempty"`
	Phone    string    `json:"phone,omitempty"`
	Role     int       `json:"role"`
}
