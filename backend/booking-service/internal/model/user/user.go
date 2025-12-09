package user

import "github.com/google/uuid"

type User struct {
	ID          uuid.UUID `json:"id"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Role        int       `json:"role"`
}
