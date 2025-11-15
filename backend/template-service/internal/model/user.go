package model

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null" validate:"required,min=3,max=50,alphanum"`
	Password  string         `json:"-" gorm:"not null" validate:"required,password"`
	FirstName string         `json:"first_name" gorm:"not null" validate:"required,min=1,max=50"`
	LastName  string         `json:"last_name" gorm:"not null" validate:"required,min=1,max=50"`
	Phone     string         `json:"phone" gorm:"index" validate:"omitempty,phone"`
	Role      string         `json:"role" gorm:"not null;default:'user'" validate:"required,oneof=admin user moderator"`
	Status    string         `json:"status" gorm:"not null;default:'active'" validate:"required,oneof=active inactive suspended"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserCreateRequest represents the request payload for creating a user
type UserCreateRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Password  string `json:"password" validate:"required,password"`
	FirstName string `json:"first_name" validate:"required,min=1,max=50"`
	LastName  string `json:"last_name" validate:"required,min=1,max=50"`
	Phone     string `json:"phone" validate:"omitempty,phone"`
	Role      string `json:"role" validate:"omitempty,oneof=admin user moderator"`
}

// UserUpdateRequest represents the request payload for updating a user
type UserUpdateRequest struct {
	Email     *string `json:"email" validate:"omitempty,email"`
	Username  *string `json:"username" validate:"omitempty,min=3,max=50,alphanum"`
	FirstName *string `json:"first_name" validate:"omitempty,min=1,max=50"`
	LastName  *string `json:"last_name" validate:"omitempty,min=1,max=50"`
	Phone     *string `json:"phone" validate:"omitempty,phone"`
	Role      *string `json:"role" validate:"omitempty,oneof=admin user moderator"`
	Status    *string `json:"status" validate:"omitempty,oneof=active inactive suspended"`
}

// UserResponse represents the response payload for user data
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts User model to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
		Role:      u.Role,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// TableName returns the table name for User model
func (User) TableName() string {
	return "users"
}
