package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"bus-booking/shared/db/mocks"
	"bus-booking/user-service/config"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestNewTokenManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)

	// Use real JWTManager instead of mock (no import cycle)
	jwtCfg := &config.JWTConfig{
		SecretKey:        "test-secret",
		RefreshSecretKey: "test-refresh",
		AccessTokenTTL:   15 * time.Minute,
		RefreshTokenTTL:  24 * time.Hour,
	}
	jwtManager := NewJWTManager(jwtCfg)

	manager := NewTokenManager(mockRedis, jwtManager)

	assert.NotNil(t, manager)
	assert.IsType(t, &TokenBlacklistManagerImpl{}, manager)
}

func TestBlacklist_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)
	jwtManager := NewJWTManager(&config.JWTConfig{SecretKey: "test", RefreshSecretKey: "test"})
	manager := NewTokenManager(mockRedis, jwtManager)

	ctx := context.Background()
	token := "valid.jwt.token"
	key := fmt.Sprintf("blacklist:token:%s", token)

	// Manager will try to calculate TTL - we need validateToken to work
	// For simplicity, mock redis operations
	mockRedis.EXPECT().
		Set(ctx, key, "1", gomock.Any()).
		Return(nil).
		Times(1)

	result := manager.Blacklist(ctx, token)

	assert.True(t, result)
}

func TestBlacklist_RedisError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)
	jwtManager := NewJWTManager(&config.JWTConfig{SecretKey: "test", RefreshSecretKey: "test"})
	manager := NewTokenManager(mockRedis, jwtManager)

	ctx := context.Background()
	token := "valid.jwt.token"
	key := fmt.Sprintf("blacklist:token:%s", token)

	mockRedis.EXPECT().
		Set(ctx, key, "1", gomock.Any()).
		Return(assert.AnError).
		Times(1)

	result := manager.Blacklist(ctx, token)

	assert.False(t, result)
}

func TestIsBlacklisted_True(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)
	jwtManager := NewJWTManager(&config.JWTConfig{SecretKey: "test", RefreshSecretKey: "test"})
	manager := NewTokenManager(mockRedis, jwtManager)

	ctx := context.Background()
	token := "blacklisted.token"
	key := fmt.Sprintf("blacklist:token:%s", token)

	mockRedis.EXPECT().
		Exists(ctx, key).
		Return(int64(1), nil).
		Times(1)

	result := manager.IsBlacklisted(ctx, token)

	assert.True(t, result)
}

func TestIsBlacklisted_False(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)
	jwtManager := NewJWTManager(&config.JWTConfig{SecretKey: "test", RefreshSecretKey: "test"})
	manager := NewTokenManager(mockRedis, jwtManager)

	ctx := context.Background()
	token := "valid.token"
	key := fmt.Sprintf("blacklist:token:%s", token)

	mockRedis.EXPECT().
		Exists(ctx, key).
		Return(int64(0), nil).
		Times(1)

	result := manager.IsBlacklisted(ctx, token)

	assert.False(t, result)
}

func TestIsBlacklisted_RedisError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)
	jwtManager := NewJWTManager(&config.JWTConfig{SecretKey: "test", RefreshSecretKey: "test"})
	manager := NewTokenManager(mockRedis, jwtManager)

	ctx := context.Background()
	token := "some.token"
	key := fmt.Sprintf("blacklist:token:%s", token)

	mockRedis.EXPECT().
		Exists(ctx, key).
		Return(int64(0), assert.AnError).
		Times(1)

	result := manager.IsBlacklisted(ctx, token)

	// Returns false on error for fail-safe
	assert.False(t, result)
}

func TestBlacklistUserTokens_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)
	jwtManager := NewJWTManager(&config.JWTConfig{SecretKey: "test", RefreshSecretKey: "test"})
	manager := NewTokenManager(mockRedis, jwtManager)

	ctx := context.Background()
	userID := uuid.New()
	key := fmt.Sprintf("blacklist:user:%s", userID.String())

	mockRedis.EXPECT().
		Set(ctx, key, gomock.Any(), 7*24*time.Hour).
		Return(nil).
		Times(1)

	result := manager.BlacklistUserTokens(ctx, userID)

	assert.True(t, result)
}

func TestBlacklistUserTokens_RedisError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)
	jwtManager := NewJWTManager(&config.JWTConfig{SecretKey: "test", RefreshSecretKey: "test"})
	manager := NewTokenManager(mockRedis, jwtManager)

	ctx := context.Background()
	userID := uuid.New()
	key := fmt.Sprintf("blacklist:user:%s", userID.String())

	mockRedis.EXPECT().
		Set(ctx, key, gomock.Any(), 7*24*time.Hour).
		Return(assert.AnError).
		Times(1)

	result := manager.BlacklistUserTokens(ctx, userID)

	assert.False(t, result)
}

func TestIsUserTokensBlacklisted_True(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)
	jwtManager := NewJWTManager(&config.JWTConfig{SecretKey: "test", RefreshSecretKey: "test"})
	manager := NewTokenManager(mockRedis, jwtManager)

	ctx := context.Background()
	userID := uuid.New()
	key := fmt.Sprintf("blacklist:user:%s", userID.String())

	// Token issued at timestamp 1000
	tokenIssuedAt := int64(1000)

	// User blacklisted at timestamp 2000
	blacklistTime := "2000"

	mockRedis.EXPECT().
		Get(ctx, key).
		Return(blacklistTime, nil).
		Times(1)

	result := manager.IsUserTokensBlacklisted(ctx, userID, tokenIssuedAt)

	// Token issued before blacklist time = blacklisted
	assert.True(t, result)
}

func TestIsUserTokensBlacklisted_False(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)
	jwtManager := NewJWTManager(&config.JWTConfig{SecretKey: "test", RefreshSecretKey: "test"})
	manager := NewTokenManager(mockRedis, jwtManager)

	ctx := context.Background()
	userID := uuid.New()
	key := fmt.Sprintf("blacklist:user:%s", userID.String())

	// Token issued at timestamp 3000
	tokenIssuedAt := int64(3000)

	// User blacklisted at timestamp 2000
	blacklistTime := "2000"

	mockRedis.EXPECT().
		Get(ctx, key).
		Return(blacklistTime, nil).
		Times(1)

	result := manager.IsUserTokensBlacklisted(ctx, userID, tokenIssuedAt)

	// Token issued after blacklist time = NOT blacklisted
	assert.False(t, result)
}

func TestIsUserTokensBlacklisted_NoBlacklist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)
	jwtManager := NewJWTManager(&config.JWTConfig{SecretKey: "test", RefreshSecretKey: "test"})
	manager := NewTokenManager(mockRedis, jwtManager)

	ctx := context.Background()
	userID := uuid.New()
	key := fmt.Sprintf("blacklist:user:%s", userID.String())

	mockRedis.EXPECT().
		Get(ctx, key).
		Return("", redis.Nil).
		Times(1)

	result := manager.IsUserTokensBlacklisted(ctx, userID, 1000)

	// No blacklist = not blacklisted
	assert.False(t, result)
}

func TestIsUserTokensBlacklisted_RedisError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)
	jwtManager := NewJWTManager(&config.JWTConfig{SecretKey: "test", RefreshSecretKey: "test"})
	manager := NewTokenManager(mockRedis, jwtManager)

	ctx := context.Background()
	userID := uuid.New()
	key := fmt.Sprintf("blacklist:user:%s", userID.String())

	mockRedis.EXPECT().
		Get(ctx, key).
		Return("", assert.AnError).
		Times(1)

	result := manager.IsUserTokensBlacklisted(ctx, userID, 1000)

	// Fail-safe: return false on error
	assert.False(t, result)
}

func TestIsUserTokensBlacklisted_InvalidBlacklistTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockRedisManager(ctrl)
	jwtManager := NewJWTManager(&config.JWTConfig{SecretKey: "test", RefreshSecretKey: "test"})
	manager := NewTokenManager(mockRedis, jwtManager)

	ctx := context.Background()
	userID := uuid.New()
	key := fmt.Sprintf("blacklist:user:%s", userID.String())

	// Invalid timestamp (not a number)
	mockRedis.EXPECT().
		Get(ctx, key).
		Return("invalid-timestamp", nil).
		Times(1)

	result := manager.IsUserTokensBlacklisted(ctx, userID, 1000)

	// Fail-safe: return false on parse error
	assert.False(t, result)
}
