package service

import (
	"bus-booking/shared/db"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	resetTokenPrefix = "password_reset:"
	resetTokenTTL    = 1 * time.Hour
)

type PasswordResetService interface {
	GenerateResetToken(ctx context.Context, email string) (string, error)
	ValidateResetToken(ctx context.Context, token string) (string, error)
	InvalidateResetToken(ctx context.Context, token string) error
}

type PasswordResetServiceImpl struct {
	redisClient *db.RedisManager
}

func NewPasswordResetService(redisClient *db.RedisManager) PasswordResetService {
	return &PasswordResetServiceImpl{
		redisClient: redisClient,
	}
}

// GenerateResetToken creates a secure random token and stores email in Redis
func (s *PasswordResetServiceImpl) GenerateResetToken(ctx context.Context, email string) (string, error) {
	// Generate 32-byte random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// Store email with token as key in Redis with TTL
	key := resetTokenPrefix + token
	if err := s.redisClient.Set(ctx, key, email, resetTokenTTL); err != nil {
		return "", fmt.Errorf("failed to store reset token: %w", err)
	}

	return token, nil
}

// ValidateResetToken checks if token exists and returns associated email
func (s *PasswordResetServiceImpl) ValidateResetToken(ctx context.Context, token string) (string, error) {
	key := resetTokenPrefix + token
	email, err := s.redisClient.Get(ctx, key)
	if err == redis.Nil {
		return "", fmt.Errorf("invalid or expired reset token")
	}
	if err != nil {
		return "", fmt.Errorf("failed to validate reset token: %w", err)
	}

	return email, nil
}

// InvalidateResetToken removes the token from Redis
func (s *PasswordResetServiceImpl) InvalidateResetToken(ctx context.Context, token string) error {
	key := resetTokenPrefix + token
	if err := s.redisClient.Del(ctx, key); err != nil {
		return fmt.Errorf("failed to invalidate reset token: %w", err)
	}
	return nil
}
