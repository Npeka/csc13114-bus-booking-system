package model

type SignupRequest struct {
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required,min=8"`
	Username  string   `json:"username" validate:"required,min=3,max=50,alphanum"`
	FirstName string   `json:"first_name" validate:"required,min=1,max=50"`
	LastName  string   `json:"last_name" validate:"required,min=1,max=50"`
	Phone     string   `json:"phone" validate:"omitempty,phone"`
	Role      UserRole `json:"role" validate:"omitempty"`
}

type SigninRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type OAuth2SigninRequest struct {
	Provider    string `json:"provider" validate:"required,oneof=firebase"`
	IDToken     string `json:"id_token" validate:"required"`
	AccessToken string `json:"access_token,omitempty"`
}

type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int64         `json:"expires_in"`
}

type TokenVerifyRequest struct {
	Token string `json:"token" validate:"required"`
}

type TokenVerifyResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"user_id,omitempty"`
	Email  string `json:"email,omitempty"`
	Role   string `json:"role,omitempty"`
	Name   string `json:"name,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type SignoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,password"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ConfirmResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,password"`
}
