package service

import (
	"context"
	"errors"
	"fmt"

	"firebase.google.com/go/v4/auth"

	"bus-booking/shared/constants"
	"bus-booking/shared/ginext"
	"bus-booking/user-service/config"
	"bus-booking/user-service/internal/model"
	"bus-booking/user-service/internal/repository"
	"bus-booking/user-service/internal/utils"

	"github.com/rs/zerolog/log"
)

type AuthService interface {
	Signup(ctx context.Context, req *model.SignupRequest) (*model.AuthResponse, error)
	Signin(ctx context.Context, req *model.SigninRequest) (*model.AuthResponse, error)
	OAuth2Signin(ctx context.Context, req *model.OAuth2SigninRequest) (*model.AuthResponse, error)
	Signout(ctx context.Context, userID string) error
	VerifyToken(ctx context.Context, token string) (*model.TokenVerifyResponse, error)
	RefreshToken(ctx context.Context, req *model.RefreshTokenRequest) (*model.AuthResponse, error)
}

type AuthServiceImpl struct {
	userRepo     repository.UserRepository
	jwtManager   *utils.JWTManager
	firebaseAuth *auth.Client
	config       *config.Config
}

func NewAuthService(
	userRepo repository.UserRepository,
	jwtManager *utils.JWTManager,
	firebaseAuth *auth.Client,
	config *config.Config,
) AuthService {
	return &AuthServiceImpl{
		userRepo:     userRepo,
		jwtManager:   jwtManager,
		firebaseAuth: firebaseAuth,
		config:       config,
	}
}

// Signup handles user registration
func (s *AuthServiceImpl) Signup(ctx context.Context, req *model.SignupRequest) (*model.AuthResponse, error) {
	emailExists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check email existence")
		return nil, err
	}
	if emailExists {
		return nil, ginext.NewConflictError("email already exists")
	}

	usernameExists, err := s.userRepo.UsernameExists(ctx, req.Username)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check username existence")
		return nil, err
	}
	if usernameExists {
		return nil, ginext.NewConflictError("username already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return nil, ginext.NewInternalServerError("Failed to hash password")
	}

	role := req.Role
	if role == 0 {
		role = model.RolePassenger
	}

	user := &model.User{
		Email:         req.Email,
		Username:      req.Username,
		Password:      hashedPassword,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Phone:         req.Phone,
		Role:          role,
		Status:        "active",
		EmailVerified: false,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		return nil, ginext.NewInternalServerError("Failed to create user")
	}

	return s.generateAuthResponse(user)
}

func (s *AuthServiceImpl) Signin(ctx context.Context, req *model.SigninRequest) (*model.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Err(err).Msg("Failed to get user by email")
		return nil, ginext.NewUnauthorizedError("invalid credentials")
	}

	if user == nil {
		return nil, ginext.NewUnauthorizedError("invalid credentials")
	}

	if user.Status != "active" && user.Status != "verified" {
		return nil, ginext.NewForbiddenError("account is not active")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, ginext.NewUnauthorizedError("invalid credentials")
	}

	return s.generateAuthResponse(user)
}

func (s *AuthServiceImpl) OAuth2Signin(ctx context.Context, req *model.OAuth2SigninRequest) (*model.AuthResponse, error) {
	if req.Provider != "firebase" {
		return nil, errors.New("unsupported OAuth2 provider")
	}

	if s.firebaseAuth == nil {
		log.Error().Msg("Firebase Auth is not initialized")
		return nil, errors.New("Firebase Auth is not available")
	}

	token, err := s.firebaseAuth.VerifyIDToken(ctx, req.IDToken)
	if err != nil {
		log.Error().Err(err).Msg("Failed to verify Firebase ID token")
		return nil, errors.New("invalid token")
	}

	user, err := s.userRepo.GetByFirebaseUID(ctx, token.UID)
	if err != nil && err.Error() != constants.ErrNotFound {
		log.Error().Err(err).Msg("Failed to get user by Firebase UID")
		return nil, errors.New(constants.ErrInternalServer)
	}

	if user == nil {
		email := token.Claims["email"].(string)
		name := token.Claims["name"].(string)

		user = &model.User{
			Email:         email,
			Username:      email,
			FirstName:     name,
			LastName:      "",
			Role:          model.RolePassenger,
			Status:        "verified",
			FirebaseUID:   token.UID,
			EmailVerified: token.Claims["email_verified"].(bool),
		}

		if err := s.userRepo.Create(ctx, user); err != nil {
			log.Error().Err(err).Msg("Failed to create OAuth2 user")
			return nil, errors.New(constants.ErrInternalServer)
		}
	}

	return s.generateAuthResponse(user)
}

// Signout handles user logout
func (s *AuthServiceImpl) Signout(ctx context.Context, userID string) error {
	// In a real implementation, you might want to:
	// 1. Blacklist the current access token
	// 2. Revoke refresh tokens for the user
	// 3. Clear any session data

	// For now, we'll just log the signout
	log.Info().Str("user_id", userID).Msg("User signed out")
	return nil
}

func (s *AuthServiceImpl) VerifyToken(ctx context.Context, token string) (*model.TokenVerifyResponse, error) {
	claims, err := s.jwtManager.ValidateAccessToken(token)
	if err != nil {
		return &model.TokenVerifyResponse{Valid: false}, nil
	}

	// Get user details from database
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return &model.TokenVerifyResponse{Valid: false}, nil
	}

	if user == nil || (user.Status != "active" && user.Status != "verified") {
		return &model.TokenVerifyResponse{Valid: false}, nil
	}

	return &model.TokenVerifyResponse{
		Valid:  true,
		UserID: claims.UserID.String(),
		Email:  user.Email,
		Role:   string(user.Role),
		Name:   user.FirstName + " " + user.LastName,
	}, nil
}

func (s *AuthServiceImpl) RefreshToken(ctx context.Context, req *model.RefreshTokenRequest) (*model.AuthResponse, error) {
	claims, err := s.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	userID := claims.UserID

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by ID")
		return nil, errors.New("user not found")
	}

	if user.Status != "active" && user.Status != "verified" {
		return nil, errors.New("account is not active")
	}

	return s.generateAuthResponse(user)
}

func (s *AuthServiceImpl) generateAuthResponse(user *model.User) (*model.AuthResponse, error) {
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, fmt.Sprintf("%d", user.Role))
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate access token")
		return nil, ginext.NewInternalServerError("Failed to generate access token")
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email, fmt.Sprintf("%d", user.Role))
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate refresh token")
		return nil, ginext.NewInternalServerError("Failed to generate refresh token")
	}

	return &model.AuthResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.JWT.AccessTokenTTL.Seconds()),
	}, nil
}
