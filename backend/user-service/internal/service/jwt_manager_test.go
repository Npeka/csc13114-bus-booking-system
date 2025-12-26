package service

import (
	"testing"
	"time"

	"bus-booking/user-service/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJWTManager(t *testing.T) {
	cfg := &config.JWTConfig{
		SecretKey:        "test-secret",
		RefreshSecretKey: "test-refresh-secret",
		AccessTokenTTL:   15 * time.Minute,
		RefreshTokenTTL:  24 * time.Hour,
		Issuer:           "test-issuer",
		Audience:         "test-audience",
	}

	manager := NewJWTManager(cfg)

	assert.NotNil(t, manager)
	assert.IsType(t, &JWTManagerImpl{}, manager)
}

func TestGenerateAccessToken_Success(t *testing.T) {
	cfg := &config.JWTConfig{
		SecretKey:      "test-secret-key",
		AccessTokenTTL: 15 * time.Minute,
		Issuer:         "test-issuer",
		Audience:       "test-audience",
	}

	manager := NewJWTManager(cfg)
	userID := uuid.New()
	email := "test@example.com"
	role := "user"

	token, err := manager.GenerateAccessToken(userID, email, role)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify the token structure
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.SecretKey), nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)
}

func TestGenerateRefreshToken_Success(t *testing.T) {
	cfg := &config.JWTConfig{
		RefreshSecretKey: "test-refresh-secret",
		RefreshTokenTTL:  24 * time.Hour,
		Issuer:           "test-issuer",
		Audience:         "test-audience",
	}

	manager := NewJWTManager(cfg)
	userID := uuid.New()
	email := "test@example.com"
	role := "admin"

	token, err := manager.GenerateRefreshToken(userID, email, role)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateAccessToken_Success(t *testing.T) {
	cfg := &config.JWTConfig{
		SecretKey:      "test-secret-key",
		AccessTokenTTL: 15 * time.Minute,
		Issuer:         "test-issuer",
		Audience:       "test-audience",
	}

	manager := NewJWTManager(cfg)
	userID := uuid.New()
	email := "test@example.com"
	role := "user"

	// Generate token
	token, err := manager.GenerateAccessToken(userID, email, role)
	require.NoError(t, err)

	// Validate token
	claims, err := manager.ValidateAccessToken(token)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, role, claims.Role)
	assert.Equal(t, AccessToken, claims.TokenType)
}

func TestValidateAccessToken_InvalidToken(t *testing.T) {
	cfg := &config.JWTConfig{
		SecretKey: "test-secret-key",
	}

	manager := NewJWTManager(cfg)

	// Test with invalid token
	claims, err := manager.ValidateAccessToken("invalid.token.here")

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateAccessToken_WrongTokenType(t *testing.T) {
	cfg := &config.JWTConfig{
		SecretKey:        "test-secret-key",
		RefreshSecretKey: "test-refresh-secret",
		RefreshTokenTTL:  24 * time.Hour,
		AccessTokenTTL:   15 * time.Minute,
		Issuer:           "test-issuer",
		Audience:         "test-audience",
	}

	manager := NewJWTManager(cfg)
	userID := uuid.New()

	// Generate refresh token
	refreshToken, err := manager.GenerateRefreshToken(userID, "test@example.com", "user")
	require.NoError(t, err)

	// Try to validate as access token (wrong secret + wrong type)
	claims, err := manager.ValidateAccessToken(refreshToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateRefreshToken_Success(t *testing.T) {
	cfg := &config.JWTConfig{
		RefreshSecretKey: "test-refresh-secret",
		RefreshTokenTTL:  24 * time.Hour,
		Issuer:           "test-issuer",
		Audience:         "test-audience",
	}

	manager := NewJWTManager(cfg)
	userID := uuid.New()
	email := "admin@example.com"
	role := "admin"

	// Generate refresh token
	token, err := manager.GenerateRefreshToken(userID, email, role)
	require.NoError(t, err)

	// Validate refresh token
	claims, err := manager.ValidateRefreshToken(token)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, role, claims.Role)
	assert.Equal(t, RefreshToken, claims.TokenType)
}

func TestValidateRefreshToken_InvalidToken(t *testing.T) {
	cfg := &config.JWTConfig{
		RefreshSecretKey: "test-refresh-secret",
	}

	manager := NewJWTManager(cfg)

	claims, err := manager.ValidateRefreshToken("completely.invalid.token")

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateRefreshToken_WrongSecret(t *testing.T) {
	cfg1 := &config.JWTConfig{
		RefreshSecretKey: "secret-1",
		RefreshTokenTTL:  24 * time.Hour,
		Issuer:           "test-issuer",
		Audience:         "test-audience",
	}

	cfg2 := &config.JWTConfig{
		RefreshSecretKey: "secret-2", // Different secret
	}

	manager1 := NewJWTManager(cfg1)
	manager2 := NewJWTManager(cfg2)

	// Generate with manager1
	token, err := manager1.GenerateRefreshToken(uuid.New(), "test@example.com", "user")
	require.NoError(t, err)

	// Try to validate with manager2 (wrong secret)
	claims, err := manager2.ValidateRefreshToken(token)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGenerateAndValidate_FullCycle(t *testing.T) {
	cfg := &config.JWTConfig{
		SecretKey:        "access-secret",
		RefreshSecretKey: "refresh-secret",
		AccessTokenTTL:   15 * time.Minute,
		RefreshTokenTTL:  24 * time.Hour,
		Issuer:           "test-service",
		Audience:         "test-client",
	}

	manager := NewJWTManager(cfg)
	userID := uuid.New()
	email := "fullcycle@test.com"
	role := "moderator"

	// Generate both tokens
	accessToken, err := manager.GenerateAccessToken(userID, email, role)
	require.NoError(t, err)

	refreshToken, err := manager.GenerateRefreshToken(userID, email, role)
	require.NoError(t, err)

	// Validate access token
	accessClaims, err := manager.ValidateAccessToken(accessToken)
	assert.NoError(t, err)
	assert.Equal(t, AccessToken, accessClaims.TokenType)
	assert.Equal(t, userID, accessClaims.UserID)

	// Validate refresh token
	refreshClaims, err := manager.ValidateRefreshToken(refreshToken)
	assert.NoError(t, err)
	assert.Equal(t, RefreshToken, refreshClaims.TokenType)
	assert.Equal(t, userID, refreshClaims.UserID)

	// Cross-validation should fail
	_, err = manager.ValidateRefreshToken(accessToken)
	assert.Error(t, err)

	_, err = manager.ValidateAccessToken(refreshToken)
	assert.Error(t, err)
}
