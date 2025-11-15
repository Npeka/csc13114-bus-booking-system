package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"csc13114-bus-ticket-booking-system/shared/constants"
	sharedcontext "csc13114-bus-ticket-booking-system/shared/context"
	"csc13114-bus-ticket-booking-system/shared/response"
	"csc13114-bus-ticket-booking-system/shared/validator"
	"csc13114-bus-ticket-booking-system/template-service/internal/model"
	"csc13114-bus-ticket-booking-system/template-service/internal/service"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService service.AuthServiceInterface
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	requestID := sharedcontext.GetRequestID(c)

	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Str("request_id", requestID).Str("error", err.Error()).Msg(MsgFailedToBindRequest + " login request")
		response.BadRequestResponse(c, constants.ErrInvalidRequestData)
		return
	}

	// Validate request
	if validationErrors := validator.ValidatorInstance.ValidateStructDetailed(&req); len(validationErrors) > 0 {
		log.Error().Str("request_id", requestID).Interface("errors", validationErrors).Msg("Login " + MsgValidationFailed)
		response.ValidationErrorResponse(c, validationErrors)
		return
	}

	// Authenticate user
	loginResponse, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		log.Error().
			Str("request_id", requestID).
			Str("email", req.Email).
			Str("error", err.Error()).
			Msg(MsgLoginFailed)

		if err.Error() == "invalid email or password" || err.Error() == "account is not active" {
			response.UnauthorizedResponse(c, err.Error())
			return
		}

		response.InternalServerErrorResponse(c, constants.ErrInternalServer)
		return
	}

	log.Info().
		Str("request_id", requestID).
		Uint("user_id", loginResponse.User.ID).
		Str("email", loginResponse.User.Email).
		Msg(MsgLoginSuccess)

	response.SuccessResponse(c, "Login successful", loginResponse)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	requestID := sharedcontext.GetRequestID(c)

	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Str("request_id", requestID).Str("error", err.Error()).Msg(MsgFailedToBindRequest + " refresh token request")
		response.BadRequestResponse(c, constants.ErrInvalidRequestData)
		return
	}

	// Validate request
	if validationErrors := validator.ValidatorInstance.ValidateStructDetailed(&req); len(validationErrors) > 0 {
		log.Error().Str("request_id", requestID).Interface("errors", validationErrors).Msg("Refresh token " + MsgValidationFailed)
		response.ValidationErrorResponse(c, validationErrors)
		return
	}

	// Refresh tokens
	loginResponse, err := h.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		log.Error().Str("request_id", requestID).Str("error", err.Error()).Msg(MsgTokenRefreshFailed)

		if err.Error() == "invalid refresh token" || err.Error() == "user not found" || err.Error() == "account is not active" {
			response.UnauthorizedResponse(c, err.Error())
			return
		}

		response.InternalServerErrorResponse(c, constants.ErrInternalServer)
		return
	}

	log.Info().
		Str("request_id", requestID).
		Uint("user_id", loginResponse.User.ID).
		Msg(MsgTokenRefreshSuccess)

	response.SuccessResponse(c, "Token refreshed successfully", loginResponse)
}

// ChangePassword handles password change
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	requestID := sharedcontext.GetRequestID(c)
	userID := sharedcontext.GetUserID(c)

	if userID == "" {
		response.UnauthorizedResponse(c, constants.ErrUnauthorized)
		return
	}

	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		log.Error().
			Str("request_id", requestID).
			Str("user_id", userID).
			Str("error", err.Error()).
			Msg(MsgInvalidUserID)
		response.BadRequestResponse(c, "Invalid user ID")
		return
	}

	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Str("request_id", requestID).Str("error", err.Error()).Msg(MsgFailedToBindRequest + " change password request")
		response.BadRequestResponse(c, constants.ErrInvalidRequestData)
		return
	}

	// Validate request
	if validationErrors := validator.ValidatorInstance.ValidateStructDetailed(&req); len(validationErrors) > 0 {
		log.Error().Str("request_id", requestID).Interface("errors", validationErrors).Msg("Change password " + MsgValidationFailed)
		response.ValidationErrorResponse(c, validationErrors)
		return
	}

	// Change password
	err = h.authService.ChangePassword(c.Request.Context(), uint(userIDUint), &req)
	if err != nil {
		log.Error().
			Str("request_id", requestID).
			Uint64("user_id", userIDUint).
			Str("error", err.Error()).
			Msg(MsgPasswordChangeFailed)

		if err.Error() == constants.ErrNotFound {
			response.NotFoundResponse(c, "User not found")
			return
		}

		if err.Error() == "current password is incorrect" {
			response.BadRequestResponse(c, err.Error())
			return
		}

		response.InternalServerErrorResponse(c, constants.ErrInternalServer)
		return
	}

	log.Info().
		Str("request_id", requestID).
		Uint64("user_id", userIDUint).
		Msg(MsgPasswordChangeSuccess)

	response.SuccessResponse(c, "Password changed successfully", nil)
}

// ResetPassword handles password reset request
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	requestID := sharedcontext.GetRequestID(c)

	var req model.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Str("request_id", requestID).Str("error", err.Error()).Msg(MsgFailedToBindRequest + " reset password request")
		response.BadRequestResponse(c, constants.ErrInvalidRequestData)
		return
	}

	// Validate request
	if validationErrors := validator.ValidatorInstance.ValidateStructDetailed(&req); len(validationErrors) > 0 {
		log.Error().Str("request_id", requestID).Interface("errors", validationErrors).Msg("Reset password " + MsgValidationFailed)
		response.ValidationErrorResponse(c, validationErrors)
		return
	}

	// Request password reset
	err := h.authService.ResetPassword(c.Request.Context(), &req)
	if err != nil {
		log.Error().
			Str("request_id", requestID).
			Str("email", req.Email).
			Str("error", err.Error()).
			Msg(MsgPasswordResetFailed)

		response.InternalServerErrorResponse(c, constants.ErrInternalServer)
		return
	}

	log.Info().
		Str("request_id", requestID).
		Str("email", req.Email).
		Msg(MsgPasswordResetRequested)

	response.SuccessResponse(c, "Password reset instructions sent to your email", nil)
}

// ConfirmResetPassword handles password reset confirmation
func (h *AuthHandler) ConfirmResetPassword(c *gin.Context) {
	requestID := sharedcontext.GetRequestID(c)

	var req model.ConfirmResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Str("request_id", requestID).Str("error", err.Error()).Msg(MsgFailedToBindRequest + " confirm reset password request")
		response.BadRequestResponse(c, constants.ErrInvalidRequestData)
		return
	}

	// Validate request
	if validationErrors := validator.ValidatorInstance.ValidateStructDetailed(&req); len(validationErrors) > 0 {
		log.Error().Str("request_id", requestID).Interface("errors", validationErrors).Msg("Confirm reset password " + MsgValidationFailed)
		response.ValidationErrorResponse(c, validationErrors)
		return
	}

	// Confirm password reset
	err := h.authService.ConfirmResetPassword(c.Request.Context(), &req)
	if err != nil {
		log.Error().Str("request_id", requestID).Str("error", err.Error()).Msg(MsgConfirmResetFailed)

		if err.Error() == "not implemented" {
			response.ErrorResponse(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Feature not implemented")
			return
		}

		response.InternalServerErrorResponse(c, constants.ErrInternalServer)
		return
	}

	response.SuccessResponse(c, "Password reset successfully", nil)
}
