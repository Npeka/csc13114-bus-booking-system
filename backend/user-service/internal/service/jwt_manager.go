package service

import (
	"bus-booking/user-service/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager interface {
	GenerateAccessToken(userID uuid.UUID, email, role string) (string, error)
	GenerateRefreshToken(userID uuid.UUID, email, role string) (string, error)
	ValidateAccessToken(tokenString string) (*JWTClaims, error)
	ValidateRefreshToken(tokenString string) (*JWTClaims, error)
}

type JWTManagerImpl struct {
	cfg *config.JWTConfig
}

func NewJWTManager(cfg *config.JWTConfig) JWTManager {
	return &JWTManagerImpl{cfg: cfg}
}

type JWTClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

func (jm *JWTManagerImpl) GenerateAccessToken(userID uuid.UUID, email, role string) (string, error) {
	now := time.Now()
	claims := &JWTClaims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(jm.cfg.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    jm.cfg.Issuer,
			Audience:  []string{jm.cfg.Audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.cfg.SecretKey))
}

func (jm *JWTManagerImpl) GenerateRefreshToken(userID uuid.UUID, email, role string) (string, error) {
	now := time.Now()
	claims := &JWTClaims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: RefreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(jm.cfg.RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    jm.cfg.Issuer,
			Audience:  []string{jm.cfg.Audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.cfg.RefreshSecretKey))
}

func (jm *JWTManagerImpl) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	return jm.validateToken(tokenString, jm.cfg.SecretKey, AccessToken)
}

func (jm *JWTManagerImpl) ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	return jm.validateToken(tokenString, jm.cfg.RefreshSecretKey, RefreshToken)
}

func (jm *JWTManagerImpl) validateToken(tokenString, secret string, expectedType TokenType) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		if claims.TokenType != expectedType {
			return nil, jwt.ErrTokenInvalidClaims
		}
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}
