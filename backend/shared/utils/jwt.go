package utils

import (
	"time"

	"csc13114-bus-ticket-booking-system/shared/config"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents JWT claims structure
type JWTClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// JWTManager handles JWT token operations
type JWTManager struct {
	config *config.JWTConfig
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(cfg *config.JWTConfig) *JWTManager {
	return &JWTManager{
		config: cfg,
	}
}

// GenerateAccessToken generates an access token
func (jm *JWTManager) GenerateAccessToken(userID, email, role string) (string, error) {
	now := time.Now()
	claims := &JWTClaims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(jm.config.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    jm.config.Issuer,
			Audience:  []string{jm.config.Audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.config.SecretKey))
}

// GenerateRefreshToken generates a refresh token
func (jm *JWTManager) GenerateRefreshToken(userID, email, role string) (string, error) {
	now := time.Now()
	claims := &JWTClaims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(jm.config.RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    jm.config.Issuer,
			Audience:  []string{jm.config.Audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.config.RefreshSecretKey))
}

// GenerateTokenPair generates both access and refresh tokens
func (jm *JWTManager) GenerateTokenPair(userID, email, role string) (accessToken, refreshToken string, err error) {
	accessToken, err = jm.GenerateAccessToken(userID, email, role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = jm.GenerateRefreshToken(userID, email, role)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateAccessToken validates an access token
func (jm *JWTManager) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	return jm.validateToken(tokenString, jm.config.SecretKey, "access")
}

// ValidateRefreshToken validates a refresh token
func (jm *JWTManager) ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	return jm.validateToken(tokenString, jm.config.RefreshSecretKey, "refresh")
}

// validateToken validates a token with the given secret and expected type
func (jm *JWTManager) validateToken(tokenString, secret, expectedType string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenMalformed
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, jwt.ErrTokenMalformed
	}

	// Validate token type
	if claims.TokenType != expectedType {
		return nil, jwt.ErrTokenUsedBeforeIssued
	}

	// Validate issuer and audience
	if claims.Issuer != jm.config.Issuer {
		return nil, jwt.ErrTokenInvalidIssuer
	}

	if len(claims.Audience) == 0 || claims.Audience[0] != jm.config.Audience {
		return nil, jwt.ErrTokenInvalidAudience
	}

	return claims, nil
}

// RefreshAccessToken generates a new access token using a valid refresh token
func (jm *JWTManager) RefreshAccessToken(refreshTokenString string) (string, error) {
	claims, err := jm.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", err
	}

	return jm.GenerateAccessToken(claims.UserID, claims.Email, claims.Role)
}

// ExtractTokenFromHeader extracts token from Authorization header
func ExtractTokenFromHeader(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}

// GetTokenExpiration returns the expiration time of a token
func (jm *JWTManager) GetTokenExpiration(tokenString string, isRefreshToken bool) (time.Time, error) {
	var claims *JWTClaims
	var err error

	if isRefreshToken {
		claims, err = jm.ValidateRefreshToken(tokenString)
	} else {
		claims, err = jm.ValidateAccessToken(tokenString)
	}

	if err != nil {
		return time.Time{}, err
	}

	if claims.ExpiresAt == nil {
		return time.Time{}, jwt.ErrTokenMalformed
	}

	return claims.ExpiresAt.Time, nil
}
