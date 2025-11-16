package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"bus-booking/shared/constants"
	"bus-booking/shared/utils"
	"bus-booking/template-service/internal/model"
	"bus-booking/template-service/internal/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// AuthServiceInterface defines the interface for auth service
type AuthServiceInterface interface {
	Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
	RefreshToken(ctx context.Context, req *model.RefreshTokenRequest) (*model.LoginResponse, error)
	ChangePassword(ctx context.Context, userID uint, req *model.ChangePasswordRequest) error
	ResetPassword(ctx context.Context, req *model.ResetPasswordRequest) error
	ConfirmResetPassword(ctx context.Context, req *model.ConfirmResetPasswordRequest) error
}

// AuthService implements AuthServiceInterface
type AuthService struct {
	userRepo   repository.UserRepositoryInterface
	jwtManager *utils.JWTManager
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo repository.UserRepositoryInterface,
	jwtManager *utils.JWTManager,
) AuthServiceInterface {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by email")
		return nil, errors.New(constants.ErrInternalServer)
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if user.Status != "active" {
		return nil, errors.New("account is not active")
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate tokens
	userUUID := uuid.MustParse(fmt.Sprintf("00000000-0000-0000-0000-%012d", user.ID))
	accessToken, err := s.jwtManager.GenerateAccessToken(userUUID, user.Email, user.Role)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate access token")
		return nil, errors.New(constants.ErrInternalServer)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(fmt.Sprintf("%d", user.ID), user.Email, user.Role)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate refresh token")
		return nil, errors.New(constants.ErrInternalServer)
	}

	return &model.LoginResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(time.Hour.Seconds()), // 1 hour
	}, nil
}

// RefreshToken generates new tokens using refresh token
func (s *AuthService) RefreshToken(ctx context.Context, req *model.RefreshTokenRequest) (*model.LoginResponse, error) {
	// Validate refresh token
	claims, err := s.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Get user
	userID, err := strconv.ParseUint(claims.UserID.String()[24:], 10, 32)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}
	user, err := s.userRepo.GetByID(ctx, uint(userID))
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by ID")
		return nil, errors.New(constants.ErrInternalServer)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check if user is active
	if user.Status != "active" {
		return nil, errors.New("account is not active")
	}

	// Generate new tokens
	userUUID2 := uuid.MustParse(fmt.Sprintf("00000000-0000-0000-0000-%012d", user.ID))
	accessToken, err := s.jwtManager.GenerateAccessToken(userUUID2, user.Email, user.Role)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate access token")
		return nil, errors.New(constants.ErrInternalServer)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(fmt.Sprintf("%d", user.ID), user.Email, user.Role)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate refresh token")
		return nil, errors.New(constants.ErrInternalServer)
	}

	return &model.LoginResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(time.Hour.Seconds()),
	}, nil
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(ctx context.Context, userID uint, req *model.ChangePasswordRequest) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by ID")
		return errors.New(constants.ErrInternalServer)
	}
	if user == nil {
		return errors.New(constants.ErrNotFound)
	}

	// Verify current password
	if !utils.CheckPasswordHash(req.CurrentPassword, user.Password) {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return errors.New(constants.ErrInternalServer)
	}

	// Update password
	user.Password = hashedPassword
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to update user password")
		return errors.New(constants.ErrInternalServer)
	}

	return nil
}

// ResetPassword initiates password reset process
func (s *AuthService) ResetPassword(ctx context.Context, req *model.ResetPasswordRequest) error {
	// Find user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by email")
		return errors.New(constants.ErrInternalServer)
	}
	if user == nil {
		// Don't reveal that user doesn't exist
		return nil
	}

	// TODO: Generate reset token and send email
	// This is a placeholder for the actual implementation
	log.Info().
		Uint("user_id", user.ID).
		Str("email", user.Email).
		Msg("Password reset requested")

	return nil
}

// ConfirmResetPassword confirms password reset with token
func (s *AuthService) ConfirmResetPassword(ctx context.Context, req *model.ConfirmResetPasswordRequest) error {
	// TODO: Validate reset token and update password
	// This is a placeholder for the actual implementation
	log.Info().Str("token", req.Token).Msg("Password reset confirmation requested")

	return errors.New("not implemented")
}
