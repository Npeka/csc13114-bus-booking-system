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
	return &AuthHandlerImpl{
		as: as,
	}
}

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

func (h *AuthHandlerImpl) RefreshToken(r *ginext.Request) (*ginext.Response, error) {
	req := model.RefreshTokenRequest{}
	if err := r.GinCtx.ShouldBind(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	userId := sharedcontext.GetUserID(r.GinCtx)

	authResp, err := h.as.RefreshToken(r.Context(), &req, userId)
	if err != nil {
		log.Error().Err(err).Msg("Token refresh failed")
		return nil, err
	}

	return ginext.NewSuccessResponse(authResp, "Token refreshed successfully"), nil
}

func (h *AuthHandlerImpl) Logout(r *ginext.Request) (*ginext.Response, error) {
	userId := sharedcontext.GetUserID(r.GinCtx)
	accessToken := sharedcontext.GetAccessToken(r.GinCtx)

	req := model.SignoutRequest{AccessToken: accessToken}
	if err := r.GinCtx.ShouldBind(&req); err != nil {
		log.Debug().Err(err).Msg("JSON binding failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	if err := h.as.Logout(r.Context(), req, userId); err != nil {
		log.Error().Err(err).Msg("Logout failed")
		return nil, err
	}

	return ginext.NewSuccessResponse(nil, "User logged out successfully"), nil
}
