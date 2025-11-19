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
	FirebaseAuth(ctx context.Context, req *model.FirebaseAuthRequest) (*model.AuthResponse, error)
	RefreshToken(ctx context.Context, req *model.RefreshTokenRequest, userID uuid.UUID) (*model.AuthResponse, error)
	VerifyToken(ctx context.Context, token string) (*model.TokenVerifyResponse, error)
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

func (s *AuthServiceImpl) VerifyToken(ctx context.Context, token string) (*model.TokenVerifyResponse, error) {
	claims, err := s.jwtManager.ValidateAccessToken(token)
	if err != nil {
		return nil, ginext.NewUnauthorizedError("invalid access token")
	}

	// Check blacklist - đơn giản không cần handle error phức tạp
	if s.tokenBlacklistMgr.IsTokenBlacklisted(ctx, token) {
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

	user, err := s.userRepo.GetByFirebaseUID(ctx, token.UID)
	if err == nil {
		if user.Status != "active" && user.Status != "verified" {
			return nil, ginext.NewForbiddenError("Account is not active")
		}
		return s.generateAuthResponse(user)
	}

	email := ""
	phone := ""
	fullName := ""
	avatar := ""

	if emailClaim, exists := token.Claims["email"]; exists && emailClaim != nil {
		email = emailClaim.(string)
	}
	if phoneClaim, exists := token.Claims["phone_number"]; exists && phoneClaim != nil {
		phone = phoneClaim.(string)
	}
	if nameClaim, exists := token.Claims["name"]; exists && nameClaim != nil {
		fullName = nameClaim.(string)
	}
	if pictureClaim, exists := token.Claims["picture"]; exists && pictureClaim != nil {
		avatar = pictureClaim.(string)
	}

	if fullName == "" && email != "" {
		fullName = strings.Split(email, "@")[0]
	}

	user = &model.User{
		Email:         email,
		Phone:         phone,
		FullName:      fullName,
		Avatar:        avatar,
		Role:          constants.RolePassenger,
		Status:        "verified",
		FirebaseUID:   token.UID,
		EmailVerified: false,
		PhoneVerified: false,
	}

	if emailVerified, exists := token.Claims["email_verified"]; exists && emailVerified != nil {
		user.EmailVerified = emailVerified.(bool)
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

	if claims.UserID != userID {
		return nil, ginext.NewUnauthorizedError("refresh token does not match user")
	}

	// Check blacklist - đơn giản
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
