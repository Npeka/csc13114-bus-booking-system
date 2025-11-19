package model

import "bus-booking/shared/constants"

type TokenVerifyRequest struct {
	AccessToken string `json:"access_token" validate:"required,min=1"`
}

type TokenVerifyResponse struct {
	UserID string             `json:"user_id,omitempty"`
	Email  string             `json:"email,omitempty"`
	Role   constants.UserRole `json:"role,omitempty"`
	Name   string             `json:"name,omitempty"`
}

type FirebaseAuthRequest struct {
	IDToken string `json:"id_token" validate:"required,min=1"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=1"`
}

type SignoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=1"`
	AccessToken  string `json:"access_token"` // set after middleware
}

type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int64         `json:"expires_in"`
}
