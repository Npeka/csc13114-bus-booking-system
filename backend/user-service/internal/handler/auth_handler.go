package handler

import (
	"github.com/rs/zerolog/log"

	"bus-booking/shared/ginext"
	"bus-booking/user-service/internal/model"
	"bus-booking/user-service/internal/service"
)

type AuthHandler struct {
	authService service.AuthServiceInterface
}

func NewAuthHandler(authService service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Signup(r *ginext.Request) (*ginext.Response, error) {
	log := log.With().Str("handler", "AuthHandler.Signup").Logger()

	req := model.SignupRequest{}
	r.MustBind(&req)

	if err := ginext.ValidateRequest(&req); err != nil {
		log.Error().Err(err).Msg("Validation failed")
		return nil, ginext.NewBadRequestError("Invalid request data")
	}

	authResp, err := h.authService.Signup(r.Context(), &req)
	if err != nil {
		log.Warn().Err(err).Msg("Signup failed")
		return nil, err
	}

	return ginext.NewCreatedResponse(authResp, "User registered successfully"), nil
}

func (h *AuthHandler) Signin(r *ginext.Request) (*ginext.Response, error) {
	log := log.With().Str("handler", "AuthHandler.Signin").Logger()

	req := model.SigninRequest{}
	r.MustBind(&req)

	if err := ginext.ValidateRequest(&req); err != nil {
		log.Debug().Err(err).Msg("Validation failed")
		return nil, err
	}

	authResp, err := h.authService.Signin(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Signin failed")
		if err.Error() == "invalid credentials" || err.Error() == "user not found" {
			return nil, ginext.NewUnauthorizedError("Invalid email or password")
		}
		return nil, ginext.NewInternalServerError("Sign in failed")
	}

	return ginext.NewSuccessResponse(authResp, "User signed in successfully"), nil
}

func (h *AuthHandler) OAuth2Signin(r *ginext.Request) (*ginext.Response, error) {
	log := log.With().Str("handler", "AuthHandler.OAuth2Signin").Logger()

	req := model.OAuth2SigninRequest{}
	r.MustBind(&req)

	if err := ginext.ValidateRequest(&req); err != nil {
		log.Debug().Err(err).Msg("Validation failed")
		return nil, err
	}

	authResp, err := h.authService.OAuth2Signin(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("OAuth2 signin failed")
		if err.Error() == "invalid token" || err.Error() == "invalid firebase token" {
			return nil, ginext.NewUnauthorizedError("Invalid token")
		}
		return nil, ginext.NewInternalServerError("OAuth2 sign in failed")
	}

	return ginext.NewSuccessResponse(authResp, "OAuth2 sign in successful"), nil
}

func (h *AuthHandler) Signout(r *ginext.Request) (*ginext.Response, error) {
	log := log.With().Str("handler", "AuthHandler.Signout").Logger()

	req := model.SignoutRequest{}
	r.MustBind(&req)

	if err := ginext.ValidateRequest(&req); err != nil {
		log.Debug().Err(err).Msg("Validation failed")
		return nil, err
	}

	err := h.authService.Signout(r.Context(), req.RefreshToken)
	if err != nil {
		log.Error().Err(err).Msg("Signout failed")
		return nil, ginext.NewInternalServerError("Sign out failed")
	}

	return ginext.NewSuccessResponse(nil, "User signed out successfully"), nil
}

func (h *AuthHandler) VerifyToken(r *ginext.Request) (*ginext.Response, error) {
	log := log.With().Str("handler", "AuthHandler.VerifyToken").Logger()

	req := model.TokenVerifyRequest{}
	r.MustBind(&req)

	if err := ginext.ValidateRequest(&req); err != nil {
		log.Debug().Err(err).Msg("Validation failed")
		return nil, err
	}

	verifyResp, err := h.authService.VerifyToken(r.Context(), req.Token)
	if err != nil {
		log.Error().Err(err).Msg("Token verification failed")
		if err.Error() == "invalid token" || err.Error() == "expired token" {
			return nil, ginext.NewUnauthorizedError("Invalid token")
		}
		return nil, ginext.NewInternalServerError("Token verification failed")
	}

	return ginext.NewSuccessResponse(verifyResp, "Token verified successfully"), nil
}

func (h *AuthHandler) RefreshToken(r *ginext.Request) (*ginext.Response, error) {
	log := log.With().Str("handler", "AuthHandler.RefreshToken").Logger()

	req := model.RefreshTokenRequest{}
	r.MustBind(&req)

	if err := ginext.ValidateRequest(&req); err != nil {
		log.Debug().Err(err).Msg("Validation failed")
		return nil, err
	}

	authResp, err := h.authService.RefreshToken(r.Context(), &req)
	if err != nil {
		log.Error().Err(err).Msg("Token refresh failed")
		if err.Error() == "invalid refresh token" {
			return nil, ginext.NewUnauthorizedError("Invalid refresh token")
		}
		return nil, ginext.NewInternalServerError("Token refresh failed")
	}

	return ginext.NewSuccessResponse(authResp, "Token refreshed successfully"), nil
}
