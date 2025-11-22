package service

import (
	"context"
	"fmt"
	"strings"

	"firebase.google.com/go/v4/auth"

	"bus-booking/shared/constants"
	"bus-booking/shared/ginext"
	"bus-booking/user-service/config"
	"bus-booking/user-service/internal/model"
	"bus-booking/user-service/internal/repository"
	"bus-booking/user-service/internal/utils"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type AuthService interface {
	VerifyToken(ctx context.Context, token string) (*model.TokenVerifyResponse, error)
	FirebaseAuth(ctx context.Context, req *model.FirebaseAuthRequest) (*model.AuthResponse, error)
	RefreshToken(ctx context.Context, req *model.RefreshTokenRequest, userID uuid.UUID) (*model.AuthResponse, error)
	Logout(ctx context.Context, req model.SignoutRequest, userID uuid.UUID) error
}

type AuthServiceImpl struct {
	userRepo          repository.UserRepository
	jwtManager        utils.JWTManager
	firebaseAuth      *auth.Client
	config            *config.Config
	tokenBlacklistMgr TokenBlacklistManager
}

func NewAuthService(
	userRepo repository.UserRepository,
	jwtManager utils.JWTManager,
	firebaseAuth *auth.Client,
	config *config.Config,
	tokenBlacklistMgr TokenBlacklistManager,
) AuthService {
	return &AuthServiceImpl{
		userRepo:          userRepo,
		jwtManager:        jwtManager,
		firebaseAuth:      firebaseAuth,
		config:            config,
		tokenBlacklistMgr: tokenBlacklistMgr,
	}
}

func (s *AuthServiceImpl) VerifyToken(ctx context.Context, accessToken string) (*model.TokenVerifyResponse, error) {
	claims, err := s.jwtManager.ValidateAccessToken(accessToken)
	if err != nil {
		return nil, ginext.NewUnauthorizedError("invalid access token")
	}

	// Check blacklist - đơn giản không cần handle error phức tạp
	if s.tokenBlacklistMgr.IsTokenBlacklisted(ctx, accessToken) {
		return nil, ginext.NewUnauthorizedError("token is blacklisted")
	}

	if s.tokenBlacklistMgr.IsUserTokensBlacklisted(ctx, claims.UserID, claims.IssuedAt.Unix()) {
		return nil, ginext.NewUnauthorizedError("user tokens are blacklisted")
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil || user == nil {
		return nil, ginext.NewUnauthorizedError("user not found")
	}

	if user.Status != "active" && user.Status != "verified" {
		return nil, ginext.NewUnauthorizedError("user is not active")
	}

	return &model.TokenVerifyResponse{
		UserID: claims.UserID.String(),
		Email:  user.Email,
		Role:   user.Role,
		Name:   user.FullName,
	}, nil
}

func (s *AuthServiceImpl) FirebaseAuth(ctx context.Context, req *model.FirebaseAuthRequest) (*model.AuthResponse, error) {
	if s.firebaseAuth == nil {
		log.Error().Msg("Firebase Auth is not initialized")
		return nil, ginext.NewInternalServerError("Firebase Auth is not available")
	}

	token, err := s.firebaseAuth.VerifyIDToken(ctx, req.IDToken)
	if err != nil {
		log.Error().Err(err).Msg("Failed to verify Firebase ID token")
		return nil, ginext.NewUnauthorizedError("Invalid Firebase token")
	}

	// Check if user already exists by Firebase UID
	user, err := s.userRepo.GetByFirebaseUID(ctx, token.UID)
	if err == nil && user != nil {
		// User exists, return auth response
		if user.Status != "active" && user.Status != "verified" {
			return nil, ginext.NewForbiddenError("Account is not active")
		}
		return s.generateAuthResponse(user)
	}

	// Extract claims from Firebase token
	email := ""
	phone := ""
	fullName := ""
	avatar := ""

	if emailClaim, exists := token.Claims["email"]; exists && emailClaim != nil {
		email, _ = emailClaim.(string)
	}
	if phoneClaim, exists := token.Claims["phone_number"]; exists && phoneClaim != nil {
		phone, _ = phoneClaim.(string)
	}
	if nameClaim, exists := token.Claims["name"]; exists && nameClaim != nil {
		fullName, _ = nameClaim.(string)
	}
	if pictureClaim, exists := token.Claims["picture"]; exists && pictureClaim != nil {
		avatar, _ = pictureClaim.(string)
	}

	// Generate full name from email if not provided
	if fullName == "" && email != "" {
		fullName = strings.Split(email, "@")[0]
	}
	// Fallback to phone number for full name
	if fullName == "" && phone != "" {
		fullName = phone
	}
	// If still empty, use firebase UID as fallback
	if fullName == "" {
		fullName = token.UID[:12]
	}

	// Check email verification status
	emailVerified := false
	if emailVerifyClaim, exists := token.Claims["email_verified"]; exists && emailVerifyClaim != nil {
		emailVerified, _ = emailVerifyClaim.(bool)
	}

	// Check phone verification status
	phoneVerified := false
	if phone != "" {
		phoneVerified = true
	}

	// Create new user
	user = &model.User{
		Email:         email,
		Phone:         phone,
		FullName:      fullName,
		Avatar:        avatar,
		Role:          constants.RolePassenger,
		Status:        "verified",
		FirebaseUID:   token.UID,
		EmailVerified: emailVerified,
		PhoneVerified: phoneVerified,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to create Firebase user")
		return nil, ginext.NewInternalServerError("Failed to create user")
	}

	return s.generateAuthResponse(user)
}

func (s *AuthServiceImpl) RefreshToken(ctx context.Context, req *model.RefreshTokenRequest, userID uuid.UUID) (*model.AuthResponse, error) {
	claims, err := s.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, ginext.NewUnauthorizedError("invalid refresh token")
	}

	if userID == uuid.Nil {
		userID = claims.UserID
		log.Debug().Str("userID", userID.String()).Msg("Using userID from refresh token claims")
	} else if claims.UserID != userID {
		return nil, ginext.NewUnauthorizedError("refresh token does not match user")
	}

	// Check blacklist
	if s.tokenBlacklistMgr.IsTokenBlacklisted(ctx, req.RefreshToken) {
		return nil, ginext.NewUnauthorizedError("refresh token has been revoked")
	}

	if s.tokenBlacklistMgr.IsUserTokensBlacklisted(ctx, claims.UserID, claims.IssuedAt.Unix()) {
		return nil, ginext.NewUnauthorizedError("all user tokens have been revoked")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, ginext.NewInternalServerError("user not found")
	}

	if user.Status != "active" && user.Status != "verified" {
		return nil, ginext.NewForbiddenError("account is not active")
	}

	// Blacklist old refresh token
	s.tokenBlacklistMgr.BlacklistToken(ctx, req.RefreshToken)

	return s.generateAuthResponse(user)
}

func (s *AuthServiceImpl) generateAuthResponse(user *model.User) (*model.AuthResponse, error) {
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, fmt.Sprintf("%d", user.Role))
	if err != nil {
		return nil, ginext.NewInternalServerError("Failed to generate access token")
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email, fmt.Sprintf("%d", user.Role))
	if err != nil {
		return nil, ginext.NewInternalServerError("Failed to generate refresh token")
	}

	return &model.AuthResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.JWT.AccessTokenTTL.Seconds()),
	}, nil
}

func (s *AuthServiceImpl) Logout(ctx context.Context, req model.SignoutRequest, userID uuid.UUID) error {
	claims, err := s.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return ginext.NewUnauthorizedError("invalid refresh token")
	}

	if claims.UserID != userID {
		return ginext.NewUnauthorizedError("refresh token does not match user")
	}

	// Đơn giản - chỉ blacklist token
	s.tokenBlacklistMgr.BlacklistToken(ctx, req.AccessToken)
	s.tokenBlacklistMgr.BlacklistToken(ctx, req.RefreshToken)

	return nil
}
