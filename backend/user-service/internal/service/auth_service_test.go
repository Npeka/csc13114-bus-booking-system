package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"bus-booking/shared/constants"
	"bus-booking/user-service/config"
	"bus-booking/user-service/internal/model"
	"bus-booking/user-service/internal/service/mocks"
	"bus-booking/user-service/internal/utils"

	"firebase.google.com/go/v4/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_VerifyToken_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockJWT := new(mocks.MockJWTManager)
	mockBlacklist := new(mocks.MockTokenBlacklistManager)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			AccessTokenTTL: 15 * time.Minute,
		},
	}

	service := NewAuthService(cfg, mockJWT, nil, mockBlacklist, mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	token := "valid.access.token"
	claims := &utils.JWTClaims{
		UserID:    userID,
		Email:     "test@example.com",
		Role:      "1",
		TokenType: utils.AccessToken,
	}

	now := time.Now()
	claims.IssuedAt = jwt.NewNumericDate(now)

	user := &model.User{
		ID:       userID,
		Email:    "test@example.com",
		FullName: "Test User",
		Role:     constants.RolePassenger,
		Status:   "active",
	}

	mockJWT.On("ValidateAccessToken", token).Return(claims, nil)
	mockBlacklist.On("IsTokenBlacklisted", ctx, token).Return(false)
	mockBlacklist.On("IsUserTokensBlacklisted", ctx, userID, claims.IssuedAt.Unix()).Return(false)
	mockRepo.On("GetByID", ctx, userID).Return(user, nil)

	// Act
	result, err := service.VerifyToken(ctx, token)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID.String(), result.UserID)
	assert.Equal(t, user.Email, result.Email)
	assert.Equal(t, user.Role, result.Role)
	mockJWT.AssertExpectations(t)
	mockBlacklist.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_VerifyToken_InvalidToken(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockJWT := new(mocks.MockJWTManager)
	mockBlacklist := new(mocks.MockTokenBlacklistManager)
	cfg := &config.Config{}

	service := NewAuthService(cfg, mockJWT, nil, mockBlacklist, mockRepo)
	ctx := context.Background()

	token := "invalid.token"
	mockJWT.On("ValidateAccessToken", token).Return(nil, errors.New("invalid token"))

	// Act
	result, err := service.VerifyToken(ctx, token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid access token")
	mockJWT.AssertExpectations(t)
}

func TestAuthService_VerifyToken_TokenBlacklisted(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockJWT := new(mocks.MockJWTManager)
	mockBlacklist := new(mocks.MockTokenBlacklistManager)
	cfg := &config.Config{}

	service := NewAuthService(cfg, mockJWT, nil, mockBlacklist, mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	token := "blacklisted.token"
	claims := &utils.JWTClaims{
		UserID: userID,
	}
	now := time.Now()
	claims.IssuedAt = jwt.NewNumericDate(now)

	mockJWT.On("ValidateAccessToken", token).Return(claims, nil)
	mockBlacklist.On("IsTokenBlacklisted", ctx, token).Return(true)

	// Act
	result, err := service.VerifyToken(ctx, token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "token is blacklisted")
	mockJWT.AssertExpectations(t)
	mockBlacklist.AssertExpectations(t)
}

func TestAuthService_VerifyToken_UserTokensBlacklisted(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockJWT := new(mocks.MockJWTManager)
	mockBlacklist := new(mocks.MockTokenBlacklistManager)
	cfg := &config.Config{}

	service := NewAuthService(cfg, mockJWT, nil, mockBlacklist, mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	token := "valid.token"
	claims := &utils.JWTClaims{
		UserID: userID,
	}
	now := time.Now()
	claims.IssuedAt = jwt.NewNumericDate(now)

	mockJWT.On("ValidateAccessToken", token).Return(claims, nil)
	mockBlacklist.On("IsTokenBlacklisted", ctx, token).Return(false)
	mockBlacklist.On("IsUserTokensBlacklisted", ctx, userID, claims.IssuedAt.Unix()).Return(true)

	// Act
	result, err := service.VerifyToken(ctx, token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user tokens are blacklisted")
	mockJWT.AssertExpectations(t)
	mockBlacklist.AssertExpectations(t)
}

func TestAuthService_VerifyToken_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockJWT := new(mocks.MockJWTManager)
	mockBlacklist := new(mocks.MockTokenBlacklistManager)
	cfg := &config.Config{}

	service := NewAuthService(cfg, mockJWT, nil, mockBlacklist, mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	token := "valid.token"
	claims := &utils.JWTClaims{
		UserID: userID,
	}
	now := time.Now()
	claims.IssuedAt = jwt.NewNumericDate(now)

	mockJWT.On("ValidateAccessToken", token).Return(claims, nil)
	mockBlacklist.On("IsTokenBlacklisted", ctx, token).Return(false)
	mockBlacklist.On("IsUserTokensBlacklisted", ctx, userID, claims.IssuedAt.Unix()).Return(false)
	mockRepo.On("GetByID", ctx, userID).Return(nil, errors.New("not found"))

	// Act
	result, err := service.VerifyToken(ctx, token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user not found")
	mockJWT.AssertExpectations(t)
	mockBlacklist.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_VerifyToken_UserNotActive(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockJWT := new(mocks.MockJWTManager)
	mockBlacklist := new(mocks.MockTokenBlacklistManager)
	cfg := &config.Config{}

	service := NewAuthService(cfg, mockJWT, nil, mockBlacklist, mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	token := "valid.token"
	claims := &utils.JWTClaims{
		UserID: userID,
	}
	now := time.Now()
	claims.IssuedAt = jwt.NewNumericDate(now)

	user := &model.User{
		ID:     userID,
		Status: "suspended",
	}

	mockJWT.On("ValidateAccessToken", token).Return(claims, nil)
	mockBlacklist.On("IsTokenBlacklisted", ctx, token).Return(false)
	mockBlacklist.On("IsUserTokensBlacklisted", ctx, userID, claims.IssuedAt.Unix()).Return(false)
	mockRepo.On("GetByID", ctx, userID).Return(user, nil)

	// Act
	result, err := service.VerifyToken(ctx, token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user is not active")
	mockJWT.AssertExpectations(t)
	mockBlacklist.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_FirebaseAuth_NewUser(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockJWT := new(mocks.MockJWTManager)
	mockFirebase := new(mocks.MockFirebaseAuthClient)
	mockBlacklist := new(mocks.MockTokenBlacklistManager)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			AccessTokenTTL: 15 * time.Minute,
		},
	}

	service := &AuthServiceImpl{
		userRepo:          mockRepo,
		jwtManager:        mockJWT,
		firebaseAuth:      mockFirebase,
		config:            cfg,
		tokenBlacklistMgr: mockBlacklist,
	}

	ctx := context.Background()
	req := &model.FirebaseAuthRequest{
		IDToken: "firebase.id.token",
	}

	firebaseToken := &auth.Token{
		UID: "firebase123",
		Claims: map[string]interface{}{
			"email":          "newuser@example.com",
			"name":           "New User",
			"email_verified": true,
		},
	}

	mockFirebase.On("VerifyIDToken", ctx, req.IDToken).Return(firebaseToken, nil)
	mockRepo.On("GetByFirebaseUID", ctx, firebaseToken.UID).Return(nil, errors.New("not found"))
	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(nil)
	mockJWT.On("GenerateAccessToken", mock.AnythingOfType("uuid.UUID"), "newuser@example.com", "1").Return("access.token", nil)
	mockJWT.On("GenerateRefreshToken", mock.AnythingOfType("uuid.UUID"), "newuser@example.com", "1").Return("refresh.token", nil)

	// Act
	result, err := service.FirebaseAuth(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "access.token", result.AccessToken)
	assert.Equal(t, "refresh.token", result.RefreshToken)
	mockFirebase.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}

func TestAuthService_FirebaseAuth_ExistingUser(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockJWT := new(mocks.MockJWTManager)
	mockFirebase := new(mocks.MockFirebaseAuthClient)
	mockBlacklist := new(mocks.MockTokenBlacklistManager)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			AccessTokenTTL: 15 * time.Minute,
		},
	}

	service := &AuthServiceImpl{
		userRepo:          mockRepo,
		jwtManager:        mockJWT,
		firebaseAuth:      mockFirebase,
		config:            cfg,
		tokenBlacklistMgr: mockBlacklist,
	}

	ctx := context.Background()
	req := &model.FirebaseAuthRequest{
		IDToken: "firebase.id.token",
	}

	firebaseToken := &auth.Token{
		UID: "firebase123",
	}

	existingUser := &model.User{
		ID:          uuid.New(),
		Email:       "existing@example.com",
		FirebaseUID: &firebaseToken.UID,
		Status:      "active",
		Role:        constants.RolePassenger,
	}

	mockFirebase.On("VerifyIDToken", ctx, req.IDToken).Return(firebaseToken, nil)
	mockRepo.On("GetByFirebaseUID", ctx, firebaseToken.UID).Return(existingUser, nil)
	mockJWT.On("GenerateAccessToken", existingUser.ID, existingUser.Email, "1").Return("access.token", nil)
	mockJWT.On("GenerateRefreshToken", existingUser.ID, existingUser.Email, "1").Return("refresh.token", nil)

	// Act
	result, err := service.FirebaseAuth(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "access.token", result.AccessToken)
	mockFirebase.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}

func TestAuthService_FirebaseAuth_InvalidToken(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockJWT := new(mocks.MockJWTManager)
	mockFirebase := new(mocks.MockFirebaseAuthClient)
	mockBlacklist := new(mocks.MockTokenBlacklistManager)
	cfg := &config.Config{}

	service := &AuthServiceImpl{
		userRepo:          mockRepo,
		jwtManager:        mockJWT,
		firebaseAuth:      mockFirebase,
		config:            cfg,
		tokenBlacklistMgr: mockBlacklist,
	}

	ctx := context.Background()
	req := &model.FirebaseAuthRequest{
		IDToken: "invalid.token",
	}

	mockFirebase.On("VerifyIDToken", ctx, req.IDToken).Return(nil, errors.New("invalid token"))

	// Act
	result, err := service.FirebaseAuth(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Invalid Firebase token")
	mockFirebase.AssertExpectations(t)
}

func TestAuthService_RefreshToken_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockJWT := new(mocks.MockJWTManager)
	mockBlacklist := new(mocks.MockTokenBlacklistManager)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			AccessTokenTTL: 15 * time.Minute,
		},
	}

	service := NewAuthService(cfg, mockJWT, nil, mockBlacklist, mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	req := &model.RefreshTokenRequest{
		RefreshToken: "valid.refresh.token",
	}

	claims := &utils.JWTClaims{
		UserID: userID,
		Email:  "test@example.com",
		Role:   "1",
	}
	now := time.Now()
	claims.IssuedAt = jwt.NewNumericDate(now)

	user := &model.User{
		ID:     userID,
		Email:  "test@example.com",
		Role:   constants.RolePassenger,
		Status: "active",
	}

	mockJWT.On("ValidateRefreshToken", req.RefreshToken).Return(claims, nil)
	mockBlacklist.On("IsTokenBlacklisted", ctx, req.RefreshToken).Return(false)
	mockBlacklist.On("IsUserTokensBlacklisted", ctx, userID, claims.IssuedAt.Unix()).Return(false)
	mockRepo.On("GetByID", ctx, userID).Return(user, nil)
	mockBlacklist.On("BlacklistToken", ctx, req.RefreshToken).Return(true)
	mockJWT.On("GenerateAccessToken", userID, user.Email, "1").Return("new.access.token", nil)
	mockJWT.On("GenerateRefreshToken", userID, user.Email, "1").Return("new.refresh.token", nil)

	// Act
	result, err := service.RefreshToken(ctx, req, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "new.access.token", result.AccessToken)
	assert.Equal(t, "new.refresh.token", result.RefreshToken)
	mockJWT.AssertExpectations(t)
	mockBlacklist.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_RefreshToken_InvalidToken(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockJWT := new(mocks.MockJWTManager)
	mockBlacklist := new(mocks.MockTokenBlacklistManager)
	cfg := &config.Config{}

	service := NewAuthService(cfg, mockJWT, nil, mockBlacklist, mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	req := &model.RefreshTokenRequest{
		RefreshToken: "invalid.token",
	}

	mockJWT.On("ValidateRefreshToken", req.RefreshToken).Return(nil, errors.New("invalid token"))

	// Act
	result, err := service.RefreshToken(ctx, req, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid refresh token")
	mockJWT.AssertExpectations(t)
}

func TestAuthService_RefreshToken_UserMismatch(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockJWT := new(mocks.MockJWTManager)
	mockBlacklist := new(mocks.MockTokenBlacklistManager)
	cfg := &config.Config{}

	service := NewAuthService(cfg, mockJWT, nil, mockBlacklist, mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	differentUserID := uuid.New()
	req := &model.RefreshTokenRequest{
		RefreshToken: "valid.token",
	}

	claims := &utils.JWTClaims{
		UserID: differentUserID,
	}

	mockJWT.On("ValidateRefreshToken", req.RefreshToken).Return(claims, nil)

	// Act
	result, err := service.RefreshToken(ctx, req, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "refresh token does not match user")
	mockJWT.AssertExpectations(t)
}

func TestAuthService_Logout_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockJWT := new(mocks.MockJWTManager)
	mockBlacklist := new(mocks.MockTokenBlacklistManager)
	cfg := &config.Config{}

	service := NewAuthService(cfg, mockJWT, nil, mockBlacklist, mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	req := model.LogoutRequest{
		RefreshToken: "refresh.token",
		AccessToken:  "access.token",
	}

	claims := &utils.JWTClaims{
		UserID: userID,
	}

	mockJWT.On("ValidateRefreshToken", req.RefreshToken).Return(claims, nil)
	mockBlacklist.On("BlacklistToken", ctx, req.AccessToken).Return(true)
	mockBlacklist.On("BlacklistToken", ctx, req.RefreshToken).Return(true)

	// Act
	err := service.Logout(ctx, req, userID)

	// Assert
	assert.NoError(t, err)
	mockJWT.AssertExpectations(t)
	mockBlacklist.AssertExpectations(t)
}

func TestAuthService_Register_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockJWT := new(mocks.MockJWTManager)
	mockBlacklist := new(mocks.MockTokenBlacklistManager)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			AccessTokenTTL: 15 * time.Minute,
		},
	}

	service := NewAuthService(cfg, mockJWT, nil, mockBlacklist, mockRepo)
	ctx := context.Background()

	req := &model.RegisterRequest{
		FullName: "Test User",
		Email:    "newuser@example.com",
		Password: "password123",
	}

	// Mock email doesn't exist
	mockRepo.On("GetByEmail", ctx, req.Email).Return(nil, errors.New("not found"))
	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(nil)
	mockJWT.On("GenerateAccessToken", mock.AnythingOfType("uuid.UUID"), req.Email, "1").Return("access.token", nil)
	mockJWT.On("GenerateRefreshToken", mock.AnythingOfType("uuid.UUID"), req.Email, "1").Return("refresh.token", nil)

	// Act
	result, err := service.Register(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "access.token", result.AccessToken)
	assert.Equal(t, "refresh.token", result.RefreshToken)
	assert.NotNil(t, result.User)
	assert.Equal(t, req.Email, result.User.Email)
	assert.Equal(t, req.FullName, result.User.FullName)
	mockRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}

func TestAuthService_Register_EmailAlreadyExists(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockJWT := new(mocks.MockJWTManager)
	mockBlacklist := new(mocks.MockTokenBlacklistManager)
	cfg := &config.Config{}

	service := NewAuthService(cfg, mockJWT, nil, mockBlacklist, mockRepo)
	ctx := context.Background()

	req := &model.RegisterRequest{
		FullName: "Test User",
		Email:    "existing@example.com",
		Password: "password123",
	}

	existingUser := &model.User{
		ID:    uuid.New(),
		Email: req.Email,
	}

	mockRepo.On("GetByEmail", ctx, req.Email).Return(existingUser, nil)

	// Act
	result, err := service.Register(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Email already registered")
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockJWT := new(mocks.MockJWTManager)
	mockBlacklist := new(mocks.MockTokenBlacklistManager)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			AccessTokenTTL: 15 * time.Minute,
		},
	}

	service := NewAuthService(cfg, mockJWT, nil, mockBlacklist, mockRepo)
	ctx := context.Background()

	req := &model.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Hash the password
	hashedPassword, _ := utils.HashPassword(req.Password)

	user := &model.User{
		ID:           uuid.New(),
		Email:        req.Email,
		FullName:     "Test User",
		PasswordHash: &hashedPassword,
		Role:         constants.RolePassenger,
		Status:       "active",
	}

	mockRepo.On("GetByEmail", ctx, req.Email).Return(user, nil)
	mockJWT.On("GenerateAccessToken", user.ID, user.Email, "1").Return("access.token", nil)
	mockJWT.On("GenerateRefreshToken", user.ID, user.Email, "1").Return("refresh.token", nil)

	// Act
	result, err := service.Login(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "access.token", result.AccessToken)
	assert.Equal(t, "refresh.token", result.RefreshToken)
	mockRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockJWT := new(mocks.MockJWTManager)
	mockBlacklist := new(mocks.MockTokenBlacklistManager)
	cfg := &config.Config{}

	service := NewAuthService(cfg, mockJWT, nil, mockBlacklist, mockRepo)
	ctx := context.Background()

	req := &model.LoginRequest{
		Email:    "notfound@example.com",
		Password: "password123",
	}

	mockRepo.On("GetByEmail", ctx, req.Email).Return(nil, errors.New("not found"))

	// Act
	result, err := service.Login(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Invalid email or password")
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_NoPasswordSet(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockJWT := new(mocks.MockJWTManager)
	mockBlacklist := new(mocks.MockTokenBlacklistManager)
	cfg := &config.Config{}

	service := NewAuthService(cfg, mockJWT, nil, mockBlacklist, mockRepo)
	ctx := context.Background()

	req := &model.LoginRequest{
		Email:    "firebase@example.com",
		Password: "password123",
	}

	user := &model.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: nil, // Firebase user without password
		Status:       "active",
	}

	mockRepo.On("GetByEmail", ctx, req.Email).Return(user, nil)

	// Act
	result, err := service.Login(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Invalid email or password")
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockJWT := new(mocks.MockJWTManager)
	mockBlacklist := new(mocks.MockTokenBlacklistManager)
	cfg := &config.Config{}

	service := NewAuthService(cfg, mockJWT, nil, mockBlacklist, mockRepo)
	ctx := context.Background()

	req := &model.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	// Hash a different password
	correctPassword := "correctpassword"
	hashedPassword, _ := utils.HashPassword(correctPassword)

	user := &model.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: &hashedPassword,
		Status:       "active",
	}

	mockRepo.On("GetByEmail", ctx, req.Email).Return(user, nil)

	// Act
	result, err := service.Login(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Invalid email or password")
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotActive(t *testing.T) {
	// Arrange
	mockRepo := new(mocks.MockUserRepository)
	mockJWT := new(mocks.MockJWTManager)
	mockBlacklist := new(mocks.MockTokenBlacklistManager)
	cfg := &config.Config{}

	service := NewAuthService(cfg, mockJWT, nil, mockBlacklist, mockRepo)
	ctx := context.Background()

	req := &model.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	hashedPassword, _ := utils.HashPassword(req.Password)

	user := &model.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: &hashedPassword,
		Status:       "suspended", // Not active
	}

	mockRepo.On("GetByEmail", ctx, req.Email).Return(user, nil)

	// Act
	result, err := service.Login(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Account is not active")
	mockRepo.AssertExpectations(t)
}
