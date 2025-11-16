package handler

import (
	"github.com/rs/zerolog/log"

	"bus-booking/shared/ginext"
	"bus-booking/user-service/internal/model"
	"bus-booking/user-service/internal/service"
)

type AuthHandler interface {
	Signup(r *ginext.Request) (*ginext.Response, error)
	Signin(r *ginext.Request) (*ginext.Response, error)
	OAuth2Signin(r *ginext.Request) (*ginext.Response, error)
	Signout(r *ginext.Request) (*ginext.Response, error)
	VerifyToken(r *ginext.Request) (*ginext.Response, error)
	RefreshToken(r *ginext.Request) (*ginext.Response, error)
}

type AuthHandlerImpl struct {
	as service.AuthService
}

func NewAuthHandler(as service.AuthService) AuthHandler {
	return &AuthHandlerImpl{
		as: as,
	}
}

func (h *AuthHandlerImpl) Signup(r *ginext.Request) (*ginext.Response, error) {
	log := log.With().Str("handler", "AuthHandler.Signup").Logger()

	req := model.SignupRequest{}
	r.MustBind(&req)

	if err := ginext.ValidateRequest(&req); err != nil {
		log.Error().Err(err).Msg("Validation failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	authResp, err := h.as.Signup(r.Context(), &req)
	if err != nil {
		log.Warn().Err(err).Msg("Signup failed")
		return nil, err
	}

	return ginext.NewCreatedResponse(authResp, "User registered successfully"), nil
}

func (h *AuthHandlerImpl) Signin(r *ginext.Request) (*ginext.Response, error) {
	log := log.With().Str("handler", "AuthHandler.Signin").Logger()

	req := model.SigninRequest{}
	r.MustBind(&req)

	if err := ginext.ValidateRequest(&req); err != nil {
		log.Debug().Err(err).Msg("Validation failed")
		return nil, err
	}

	authResp, err := h.as.Signin(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Signin failed")
		if err.Error() == "invalid credentials" || err.Error() == "user not found" {
			return nil, ginext.NewUnauthorizedError("Invalid email or password")
		}
		return nil, ginext.NewInternalServerError("Sign in failed")
	}

	return ginext.NewSuccessResponse(authResp, "User signed in successfully"), nil
}

func (h *AuthHandlerImpl) OAuth2Signin(r *ginext.Request) (*ginext.Response, error) {
	log := log.With().Str("handler", "AuthHandler.OAuth2Signin").Logger()

	req := model.OAuth2SigninRequest{}
	r.MustBind(&req)

	if err := ginext.ValidateRequest(&req); err != nil {
		log.Debug().Err(err).Msg("Validation failed")
		return nil, err
	}

	authResp, err := h.as.OAuth2Signin(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("OAuth2 signin failed")
		if err.Error() == "invalid token" || err.Error() == "invalid firebase token" {
			return nil, ginext.NewUnauthorizedError("Invalid token")
		}
		return nil, ginext.NewInternalServerError("OAuth2 sign in failed")
	}

	return ginext.NewSuccessResponse(authResp, "OAuth2 sign in successful"), nil
}

func (h *AuthHandlerImpl) Signout(r *ginext.Request) (*ginext.Response, error) {
	log := log.With().Str("handler", "AuthHandler.Signout").Logger()

	req := model.SignoutRequest{}
	r.MustBind(&req)

	if err := ginext.ValidateRequest(&req); err != nil {
		log.Debug().Err(err).Msg("Validation failed")
		return nil, err
	}

	err := h.as.Signout(r.Context(), req.RefreshToken)
	if err != nil {
		log.Error().Err(err).Msg("Signout failed")
		return nil, ginext.NewInternalServerError("Sign out failed")
	}

	return ginext.NewSuccessResponse(nil, "User signed out successfully"), nil
}

func (h *AuthHandlerImpl) VerifyToken(r *ginext.Request) (*ginext.Response, error) {
	log := log.With().Str("handler", "AuthHandler.VerifyToken").Logger()

	req := model.TokenVerifyRequest{}
	r.MustBind(&req)

	if err := ginext.ValidateRequest(&req); err != nil {
		log.Debug().Err(err).Msg("Validation failed")
		return nil, err
	}

	verifyResp, err := h.as.VerifyToken(r.Context(), req.Token)
	if err != nil {
		log.Error().Err(err).Msg("Token verification failed")
		if err.Error() == "invalid token" || err.Error() == "expired token" {
			return nil, ginext.NewUnauthorizedError("Invalid token")
		}
		return nil, ginext.NewInternalServerError("Token verification failed")
	}

	return ginext.NewSuccessResponse(verifyResp, "Token verified successfully"), nil
}

func (h *AuthHandlerImpl) RefreshToken(r *ginext.Request) (*ginext.Response, error) {
	log := log.With().Str("handler", "AuthHandler.RefreshToken").Logger()

	req := model.RefreshTokenRequest{}
	r.MustBind(&req)

	if err := ginext.ValidateRequest(&req); err != nil {
		log.Debug().Err(err).Msg("Validation failed")
		return nil, err
	}

	authResp, err := h.as.RefreshToken(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Token refresh failed")
		if err.Error() == "invalid refresh token" {
			return nil, ginext.NewUnauthorizedError("Invalid refresh token")
		}
		return nil, ginext.NewInternalServerError("Token refresh failed")
	}

	return ginext.NewSuccessResponse(authResp, "Token refreshed successfully"), nil
}
