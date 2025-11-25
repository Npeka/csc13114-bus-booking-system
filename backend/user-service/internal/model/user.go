package model

import (
	"time"

	"bus-booking/shared/constants"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID          `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email         string             `json:"email" gorm:"index;uniqueIndex:,type:NULLS NOT DISTINCT"`
	Phone         string             `json:"phone" gorm:"index"`
	FullName      string             `json:"full_name" gorm:"not null"`
	Avatar        string             `json:"avatar" gorm:"type:text"`
	Role          constants.UserRole `json:"role" gorm:"not null;default:1"`
	Status        string             `json:"status" gorm:"not null;default:'active'"`
	FirebaseUID   *string            `json:"-" gorm:"uniqueIndex:,type:NULLS NOT DISTINCT"`
	PasswordHash  *string            `json:"-" gorm:"type:text"`
	EmailVerified bool               `json:"email_verified" gorm:"default:false"`
	PhoneVerified bool               `json:"phone_verified" gorm:"default:false"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	DeletedAt     gorm.DeletedAt     `json:"-" gorm:"index"`
}

type UserCreateRequest struct {
	FirebaseUID string             `json:"firebase_uid" form:"firebase_uid" binding:"required"`
	Email       string             `json:"email" form:"email" binding:"omitempty,email"`
	Phone       string             `json:"phone" form:"phone" binding:"omitempty"`
	FullName    string             `json:"full_name" form:"full_name" binding:"required,min=1,max=100"`
	Avatar      string             `json:"avatar" form:"avatar" binding:"omitempty,url"`
	Role        constants.UserRole `json:"role" form:"role" binding:"omitempty"`
}

type UserUpdateRequest struct {
	Email    *string             `json:"email" form:"email" binding:"omitempty,email"`
	Phone    *string             `json:"phone" form:"phone" binding:"omitempty"`
	FullName *string             `json:"full_name" form:"full_name" binding:"omitempty,min=1,max=100"`
	Avatar   *string             `json:"avatar" form:"avatar" binding:"omitempty,url"`
	Role     *constants.UserRole `json:"role" form:"role" binding:"omitempty"`
	Status   *string             `json:"status" form:"status" binding:"omitempty,oneof=active inactive suspended verified"`
}

type UserResponse struct {
	ID            uuid.UUID          `json:"id"`
	Email         string             `json:"email,omitempty"`
	Phone         string             `json:"phone,omitempty"`
	FullName      string             `json:"full_name"`
	Avatar        string             `json:"avatar,omitempty"`
	Role          constants.UserRole `json:"role"`
	Status        string             `json:"status"`
	EmailVerified bool               `json:"email_verified"`
	PhoneVerified bool               `json:"phone_verified"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		Phone:         u.Phone,
		FullName:      u.FullName,
		Avatar:        u.Avatar,
		Role:          u.Role,
		Status:        u.Status,
		EmailVerified: u.EmailVerified,
		PhoneVerified: u.PhoneVerified,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

func (User) TableName() string {
	return "users"
}

type UserListQuery struct {
	Page     int    `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit    int    `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Search   string `form:"search" binding:"omitempty,max=100" json:"search"`
	Role     string `form:"role" binding:"omitempty,oneof=1 2 4 8" json:"role"`
	Status   string `form:"status" binding:"omitempty,oneof=active inactive suspended verified" json:"status"`
	SortBy   string `form:"sort_by" binding:"omitempty,oneof=created_at updated_at email phone full_name" json:"sort_by"`
	SortDesc bool   `form:"sort_desc" json:"sort_desc"`
}

type UserStatusUpdateRequest struct {
	Status string `json:"status" form:"status" binding:"required,oneof=active inactive suspended verified"`
}
