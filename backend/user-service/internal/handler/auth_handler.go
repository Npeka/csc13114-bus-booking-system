package handler

import (
	"github.com/rs/zerolog/log"

	sharedcontext "bus-booking/shared/context"
	"bus-booking/shared/ginext"
	"bus-booking/user-service/internal/model"
	"bus-booking/user-service/internal/service"
)

type AuthHandler interface {
	// gateway verified endpoints
	VerifyToken(r *ginext.Request) (*ginext.Response, error)

	// public endpoints
	FirebaseAuth(r *ginext.Request) (*ginext.Response, error)
	Register(r *ginext.Request) (*ginext.Response, error)
	Login(r *ginext.Request) (*ginext.Response, error)
	Logout(r *ginext.Request) (*ginext.Response, error)
	ForgotPassword(r *ginext.Request) (*ginext.Response, error)
	VerifyOTP(r *ginext.Request) (*ginext.Response, error)
	ResetPassword(r *ginext.Request) (*ginext.Response, error)
	RefreshToken(r *ginext.Request) (*ginext.Response, error)

	// internal
	CreateGuestAccount(r *ginext.Request) (*ginext.Response, error)
}

type AuthHandlerImpl struct {
	as service.AuthService
}

func NewAuthHandler(as service.AuthService) AuthHandler {
	return &AuthHandlerImpl{as: as}
}

// VerifyToken godoc
// @Summary Verify access token
// @Description Verifies the validity of an access token and returns user information
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.TokenVerifyRequest true "Token verification request"
// @Success 200 {object} ginext.Response{data=model.TokenVerifyResponse} "Token verified successfully"
// @Failure 400 {object} ginext.Response "Invalid request data"
// @Failure 401 {object} ginext.Response "Invalid or expired token"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /auth/verify [post]
func (h *AuthHandlerImpl) VerifyToken(r *ginext.Request) (*ginext.Response, error) {
	req := model.TokenVerifyRequest{}
	if err := r.GinCtx.ShouldBind(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	res, err := h.as.VerifyToken(r.Context(), req.AccessToken)
	if err != nil {
		log.Error().Err(err).Msg("Token verification failed")
		return nil, err
	}

	return ginext.NewSuccessResponse(res), nil
}

// FirebaseAuth godoc
// @Summary Authenticate with Firebase
// @Description Authenticates a user using Firebase ID token and returns access/refresh tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.FirebaseAuthRequest true "Firebase authentication request"
// @Success 200 {object} ginext.Response{data=model.AuthResponse} "Firebase authentication successful"
// @Failure 400 {object} ginext.Response "Invalid request data"
// @Failure 401 {object} ginext.Response "Invalid Firebase token"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /auth/firebase [post]
func (h *AuthHandlerImpl) FirebaseAuth(r *ginext.Request) (*ginext.Response, error) {
	req := model.FirebaseAuthRequest{}
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	res, err := h.as.FirebaseAuth(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Firebase auth failed")
		return nil, err
	}

	return ginext.NewSuccessResponse(res), nil
}

// Register godoc
// @Summary Register new user with email and password
// @Description Creates a new user account using email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "Registration request"
// @Success 200 {object} ginext.Response{data=model.AuthResponse} "Registration successful"
// @Failure 400 {object} ginext.Response "Invalid request data or email already registered"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandlerImpl) Register(r *ginext.Request) (*ginext.Response, error) {
	req := model.RegisterRequest{}
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	res, err := h.as.Register(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Registration failed")
		return nil, err
	}

	return ginext.NewSuccessResponse(res), nil
}

// EmailPasswordLogin godoc
// @Summary Login with email and password
// @Description Authenticates a user using email and password and returns access/refresh tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Email and password login request"
// @Success 200 {object} ginext.Response{data=model.AuthResponse} "Login successful"
// @Failure 400 {object} ginext.Response "Invalid request data"
// @Failure 401 {object} ginext.Response "Invalid email or password"
// @Failure 403 {object} ginext.Response "Account is not active"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandlerImpl) Login(r *ginext.Request) (*ginext.Response, error) {
	req := model.LoginRequest{}
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	res, err := h.as.Login(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Email/password login failed")
		return nil, err
	}

	return ginext.NewSuccessResponse(res), nil
}

// Logout godoc
// @Summary Logout user
// @Description Invalidates the user's access and refresh tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.LogoutRequest true "Logout request"
// @Success 200 {object} ginext.Response "User logged out successfully"
// @Failure 400 {object} ginext.Response "Invalid request data"
// @Failure 401 {object} ginext.Response "Unauthorized"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /auth/logout [post]
func (h *AuthHandlerImpl) Logout(r *ginext.Request) (*ginext.Response, error) {
	userID := sharedcontext.GetUserID(r.GinCtx)
	accessToken := sharedcontext.GetAccessToken(r.GinCtx)

	req := model.LogoutRequest{AccessToken: accessToken}
	if err := r.GinCtx.ShouldBind(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Dữ liệu yêu cầu không hợp lệ")
	}

	if err := h.as.Logout(r.Context(), req, userID); err != nil {
		log.Error().Err(err).Msg("Logout failed")
		return nil, err
	}

	return ginext.NewSuccessResponse("Đăng xuất thành công"), nil
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Sends a password reset link to the user's email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.ForgotPasswordRequest true "Forgot password request"
// @Success 200 {object} ginext.Response "Password reset email sent successfully"
// @Failure 400 {object} ginext.Response "Invalid request data"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /auth/forgot-password [post]
func (h *AuthHandlerImpl) ForgotPassword(r *ginext.Request) (*ginext.Response, error) {
	req := model.ForgotPasswordRequest{}
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	if err := h.as.ForgotPassword(r.Context(), &req); err != nil {
		log.Error().Err(err).Msg("Forgot password failed")
		return nil, err
	}

	return ginext.NewSuccessResponse("Nếu email tồn tại, liên kết đặt lại mật khẩu đã được gửi"), nil
}

// ResetPassword godoc
// @Summary Reset password with token
// @Description Resets the user's password using a valid reset token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.ResetPasswordRequest true "Reset password request"
// @Success 200 {object} ginext.Response "Password reset successful"
// @Failure 400 {object} ginext.Response "Invalid request data or token"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /auth/reset-password [post]
func (h *AuthHandlerImpl) ResetPassword(r *ginext.Request) (*ginext.Response, error) {
	req := model.ResetPasswordRequest{}
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	if err := h.as.ResetPassword(r.Context(), &req); err != nil {
		log.Error().Err(err).Msg("Reset password failed")
		return nil, err
	}

	return ginext.NewSuccessResponse("Đặt lại mật khẩu thành công"), nil
}

func (h *AuthHandlerImpl) VerifyOTP(r *ginext.Request) (*ginext.Response, error) {
	var req struct {
		OTP string `json:"otp" binding:"required"`
	}
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	if err := h.as.VerifyOTP(r.Context(), req.OTP); err != nil {
		log.Error().Err(err).Msg("OTP verification failed")
		return nil, err
	}

	return ginext.NewSuccessResponse("Mã OTP hợp lệ"), nil
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generates a new access token using a valid refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} ginext.Response{data=model.AuthResponse} "Token refreshed successfully"
// @Failure 400 {object} ginext.Response "Invalid request data"
// @Failure 401 {object} ginext.Response "Invalid or expired refresh token"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /auth/refresh [post]
func (h *AuthHandlerImpl) RefreshToken(r *ginext.Request) (*ginext.Response, error) {
	req := model.RefreshTokenRequest{}
	if err := r.GinCtx.ShouldBind(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	res, err := h.as.RefreshToken(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Token refresh failed")
		return nil, err
	}

	return ginext.NewSuccessResponse(res), nil
}

// CreateGuestAccount godoc
// @Summary Create guest account
// @Description Creates a guest user account for non-authenticated bookings. Requires either email or phone.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.CreateGuestAccountRequest true "Guest account creation request"
// @Success 200 {object} ginext.Response{data=model.UserResponse} "Guest account created successfully"
// @Failure 400 {object} ginext.Response "Invalid request data or missing contact info"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /auth/guest [post]
func (h *AuthHandlerImpl) CreateGuestAccount(r *ginext.Request) (*ginext.Response, error) {
	req := model.CreateGuestAccountRequest{}
	if err := r.GinCtx.ShouldBind(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	user, err := h.as.CreateGuestAccount(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Guest account creation failed")
		return nil, err
	}

	return ginext.NewSuccessResponse(user), nil
}
