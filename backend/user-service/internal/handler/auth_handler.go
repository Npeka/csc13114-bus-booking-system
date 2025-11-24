package handler

import (
	"github.com/rs/zerolog/log"

	sharedcontext "bus-booking/shared/context"
	"bus-booking/shared/ginext"
	"bus-booking/user-service/internal/model"
	"bus-booking/user-service/internal/service"
)

type AuthHandler interface {
	VerifyToken(r *ginext.Request) (*ginext.Response, error)
	FirebaseAuth(r *ginext.Request) (*ginext.Response, error)
	RefreshToken(r *ginext.Request) (*ginext.Response, error)
	Logout(r *ginext.Request) (*ginext.Response, error)
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

	verifyResp, err := h.as.VerifyToken(r.Context(), req.AccessToken)
	if err != nil {
		log.Error().Err(err).Msg("Token verification failed")
		return nil, err
	}

	return ginext.NewSuccessResponse(verifyResp, "Token verified successfully"), nil
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

	authResp, err := h.as.FirebaseAuth(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Firebase auth failed")
		return nil, err
	}

	return ginext.NewSuccessResponse(authResp, "Firebase authentication successful"), nil
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

	userID := sharedcontext.GetUserID(r.GinCtx)

	authResp, err := h.as.RefreshToken(r.Context(), &req, userID)
	if err != nil {
		log.Error().Err(err).Msg("Token refresh failed")
		return nil, err
	}

	return ginext.NewSuccessResponse(authResp, "Token refreshed successfully"), nil
}

// Logout godoc
// @Summary Logout user
// @Description Invalidates the user's access and refresh tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.SignoutRequest true "Logout request"
// @Success 200 {object} ginext.Response "User logged out successfully"
// @Failure 400 {object} ginext.Response "Invalid request data"
// @Failure 401 {object} ginext.Response "Unauthorized"
// @Failure 500 {object} ginext.Response "Internal server error"
// @Router /auth/logout [post]
func (h *AuthHandlerImpl) Logout(r *ginext.Request) (*ginext.Response, error) {
	userID := sharedcontext.GetUserID(r.GinCtx)
	accessToken := sharedcontext.GetAccessToken(r.GinCtx)

	req := model.SignoutRequest{AccessToken: accessToken}
	if err := r.GinCtx.ShouldBind(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	if err := h.as.Logout(r.Context(), req, userID); err != nil {
		log.Error().Err(err).Msg("Logout failed")
		return nil, err
	}

	return ginext.NewSuccessResponse(nil, "User logged out successfully"), nil
}
