package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"bus-booking/shared/ginext"
	"bus-booking/user-service/internal/service"
)

type AdminHandler interface {
	BlacklistUserTokens(r *ginext.Request) (*ginext.Response, error)
	ForceLogoutUser(r *ginext.Request) (*ginext.Response, error)
}

type AdminHandlerImpl struct {
	tokenBlacklistMgr service.TokenBlacklistManager
	authService       service.AuthService
}

func NewAdminHandler(
	tokenBlacklistMgr service.TokenBlacklistManager,
	authService service.AuthService,
) AdminHandler {
	return &AdminHandlerImpl{
		tokenBlacklistMgr: tokenBlacklistMgr,
		authService:       authService,
	}
}

type BlacklistUserTokensRequest struct {
	UserID string `json:"user_id" form:"user_id" binding:"required" validate:"required,uuid"`
	Reason string `json:"reason" form:"reason" binding:"required" validate:"required,min=5,max=200"`
}

type ForceLogoutUserRequest struct {
	UserID string `json:"user_id" form:"user_id" binding:"required" validate:"required,uuid"`
	Reason string `json:"reason" form:"reason" binding:"required" validate:"required,min=5,max=200"`
}

func (h *AdminHandlerImpl) BlacklistUserTokens(r *ginext.Request) (*ginext.Response, error) {
	log := log.With().Str("handler", "AdminHandler.BlacklistUserTokens").Logger()

	if h.tokenBlacklistMgr == nil {
		return nil, ginext.NewInternalServerError("Token blacklist manager not available")
	}

	req := BlacklistUserTokensRequest{}
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	// Validate request using ginext validator
	if err := ginext.ValidateRequest(&req); err != nil {
		log.Debug().Err(err).Msg("Validation failed")
		return nil, err
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, ginext.NewBadRequestError("Invalid user ID format")
	}

	// Blacklist all tokens for the user - đơn giản
	h.tokenBlacklistMgr.BlacklistUserTokens(r.Context(), userID)

	log.Info().
		Str("user_id", req.UserID).
		Str("reason", req.Reason).
		Msg("Admin blacklisted all user tokens")

	return ginext.NewSuccessResponse(gin.H{
		"user_id": req.UserID,
		"action":  "tokens_blacklisted",
		"reason":  req.Reason,
	}), nil
}

func (h *AdminHandlerImpl) ForceLogoutUser(r *ginext.Request) (*ginext.Response, error) {
	log := log.With().Str("handler", "AdminHandler.ForceLogoutUser").Logger()

	if h.tokenBlacklistMgr == nil {
		return nil, ginext.NewInternalServerError("Token blacklist manager not available")
	}

	req := ForceLogoutUserRequest{}
	if err := r.GinCtx.ShouldBindJSON(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	// Validate request using ginext validator
	if err := ginext.ValidateRequest(&req); err != nil {
		log.Debug().Err(err).Msg("Validation failed")
		return nil, err
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, ginext.NewBadRequestError("Invalid user ID format")
	}

	// Blacklist all tokens for the user - đơn giản
	h.tokenBlacklistMgr.BlacklistUserTokens(r.Context(), userID)

	log.Warn().
		Str("user_id", req.UserID).
		Str("reason", req.Reason).
		Msg("Admin forced user logout")

	return ginext.NewSuccessResponse(gin.H{
		"user_id": req.UserID,
		"action":  "force_logout",
		"reason":  req.Reason,
	}), nil
}
