package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole uint8

const (
	RolePassenger UserRole = 1 << iota // bit 0: 1
	RoleAdmin                          // bit 1: 2
	RoleOperator                       // bit 2: 4
	RoleSupport                        // bit 3: 8
)

type User struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email         string         `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Username      string         `json:"username" gorm:"uniqueIndex;not null" validate:"required,min=3,max=50,alphanum"`
	Password      string         `json:"-" gorm:"not null"`
	FirstName     string         `json:"first_name" gorm:"not null" validate:"required,min=1,max=50"`
	LastName      string         `json:"last_name" gorm:"not null" validate:"required,min=1,max=50"`
	Phone         string         `json:"phone" gorm:"index" validate:"omitempty,phone"`
	Role          UserRole       `json:"role" gorm:"not null;default:1"`
	Status        string         `json:"status" gorm:"not null;default:'active'" validate:"required,oneof=active inactive suspended verified"`
	FirebaseUID   string         `json:"-" gorm:"index"`
	EmailVerified bool           `json:"email_verified" gorm:"default:false"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type UserCreateRequest struct {
	Email     string   `json:"email" validate:"required,email"`
	Username  string   `json:"username" validate:"required,min=3,max=50,alphanum"`
	Password  string   `json:"password" validate:"required,password"`
	FirstName string   `json:"first_name" validate:"required,min=1,max=50"`
	LastName  string   `json:"last_name" validate:"required,min=1,max=50"`
	Phone     string   `json:"phone" validate:"omitempty,phone"`
	Role      UserRole `json:"role" validate:"omitempty"`
}

type UserUpdateRequest struct {
	Email     *string   `json:"email" validate:"omitempty,email"`
	Username  *string   `json:"username" validate:"omitempty,min=3,max=50,alphanum"`
	FirstName *string   `json:"first_name" validate:"omitempty,min=1,max=50"`
	LastName  *string   `json:"last_name" validate:"omitempty,min=1,max=50"`
	Phone     *string   `json:"phone" validate:"omitempty,phone"`
	Role      *UserRole `json:"role" validate:"omitempty"`
	Status    *string   `json:"status" validate:"omitempty,oneof=active inactive suspended verified"`
}

type UserResponse struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Phone         string    `json:"phone"`
	Role          UserRole  `json:"role"`
	Status        string    `json:"status"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		Username:      u.Username,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Phone:         u.Phone,
		Role:          u.Role,
		Status:        u.Status,
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

// BeforeCreate hook to generate UUID
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

func (User) TableName() string {
	return "users"
}
