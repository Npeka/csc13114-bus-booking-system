package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"bus-booking/shared/db"
	"bus-booking/user-service/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type TokenManager interface {
	// Core blacklist operations
	Blacklist(ctx context.Context, token string) bool
	IsBlacklisted(ctx context.Context, token string) bool

	// User-wide blacklist
	BlacklistUserTokens(ctx context.Context, userID uuid.UUID) bool
	IsUserTokensBlacklisted(ctx context.Context, userID uuid.UUID, tokenIssuedAt int64) bool
}

type TokenBlacklistManagerImpl struct {
	redisClient db.RedisManager
	jwtManager  utils.JWTManager
}

func NewTokenManager(
	redisManager db.RedisManager,
	jwtManager utils.JWTManager,
) TokenManager {
	return &TokenBlacklistManagerImpl{
		redisClient: redisManager,
		jwtManager:  jwtManager,
	}
}

func (tbm *TokenBlacklistManagerImpl) Blacklist(ctx context.Context, token string) bool {
	key := fmt.Sprintf("blacklist:token:%s", token)

	ttl := tbm.calculateTokenTTL(token)
	if ttl <= 0 {
		log.Debug().Str("token", token[:10]+"...").Msg("Token already expired, skip blacklist")
		return true
	}

	if err := tbm.redisClient.Set(ctx, key, "1", ttl); err != nil {
		log.Error().Err(err).Msg("Failed to blacklist token")
		return false
	}

	return true
}

func (tbm *TokenBlacklistManagerImpl) IsBlacklisted(ctx context.Context, token string) bool {
	key := fmt.Sprintf("blacklist:token:%s", token)

	exists, err := tbm.redisClient.Exists(ctx, key)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check token blacklist")
		return false
	}

	return exists > 0
}

// BlacklistUserTokens - Blacklist tất cả token của user (force logout all devices)
func (tbm *TokenBlacklistManagerImpl) BlacklistUserTokens(ctx context.Context, userID uuid.UUID) bool {
	key := fmt.Sprintf("blacklist:user:%s", userID.String())

	// Lưu timestamp hiện tại, TTL 7 ngày
	err := tbm.redisClient.Set(ctx, key, time.Now().Unix(), 7*24*time.Hour)
	if err != nil {
		log.Warn().Err(err).Str("user_id", userID.String()).Msg("Failed to blacklist user tokens")
		return false
	}

	log.Info().Str("user_id", userID.String()).Msg("All user tokens blacklisted")
	return true
}

// IsUserTokensBlacklisted - Check token của user có bị blacklist không
func (tbm *TokenBlacklistManagerImpl) IsUserTokensBlacklisted(ctx context.Context, userID uuid.UUID, tokenIssuedAt int64) bool {
	key := fmt.Sprintf("blacklist:user:%s", userID.String())

	blacklistTime, err := tbm.redisClient.Get(ctx, key)
	if err == redis.Nil {
		return false // User không bị blacklist
	}
	if err != nil {
		log.Warn().Err(err).Str("user_id", userID.String()).Msg("Failed to check user blacklist")
		return false // Fail-safe
	}

	// Token issued trước thời điểm blacklist = bị blacklist
	blacklistTimeInt, err := strconv.Atoi(blacklistTime)
	if err != nil {
		log.Warn().Err(err).Str("user_id", userID.String()).Msg("Failed to convert blacklist time to int")
		return false // Fail-safe
	}
	return tokenIssuedAt < int64(blacklistTimeInt)
}

// calculateTokenTTL - Parse JWT token và tính TTL còn lại
func (tbm *TokenBlacklistManagerImpl) calculateTokenTTL(tokenString string) time.Duration {
	// Parse token without verification (chỉ cần claims)
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse token for TTL calculation")
		return 24 * time.Hour // Fallback TTL
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Warn().Msg("Invalid token claims format")
		return 24 * time.Hour // Fallback TTL
	}

	// Lấy exp claim
	exp, ok := claims["exp"]
	if !ok {
		log.Warn().Msg("Token missing exp claim")
		return 24 * time.Hour // Fallback TTL
	}

	// Convert exp to int64
	var expTime int64
	switch v := exp.(type) {
	case float64:
		expTime = int64(v)
	case int64:
		expTime = v
	case int:
		expTime = int64(v)
	default:
		log.Warn().Interface("exp_type", v).Msg("Invalid exp claim type")
		return 24 * time.Hour // Fallback TTL
	}

	// Tính TTL còn lại
	expiration := time.Unix(expTime, 0)
	ttl := time.Until(expiration)

	// Add thêm 5 phút buffer để tránh race condition
	ttl += 5 * time.Minute

	if ttl <= 0 {
		return 0 // Token đã expire
	}

	return ttl
}
