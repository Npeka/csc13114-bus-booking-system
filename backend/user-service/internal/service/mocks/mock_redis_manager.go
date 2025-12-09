package mocks

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

// MockRedisManager is a mock implementation of shared/db.RedisManager interface
type MockRedisManager struct {
	mock.Mock
}

func (m *MockRedisManager) GetClient() *redis.Client {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*redis.Client)
}

func (m *MockRedisManager) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisManager) HealthCheck() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisManager) GetStats() *redis.PoolStats {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*redis.PoolStats)
}

func (m *MockRedisManager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisManager) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisManager) Del(ctx context.Context, keys ...string) error {
	args := m.Called(ctx, keys)
	return args.Error(0)
}

func (m *MockRedisManager) Exists(ctx context.Context, keys ...string) (int64, error) {
	args := m.Called(ctx, keys)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisManager) Expire(ctx context.Context, key string, expiration time.Duration) error {
	args := m.Called(ctx, key, expiration)
	return args.Error(0)
}

func (m *MockRedisManager) TTL(ctx context.Context, key string) (time.Duration, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (m *MockRedisManager) HSet(ctx context.Context, key string, values ...interface{}) error {
	args := m.Called(ctx, key, values)
	return args.Error(0)
}

func (m *MockRedisManager) HGet(ctx context.Context, key, field string) (string, error) {
	args := m.Called(ctx, key, field)
	return args.String(0), args.Error(1)
}

func (m *MockRedisManager) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockRedisManager) HDel(ctx context.Context, key string, fields ...string) error {
	args := m.Called(ctx, key, fields)
	return args.Error(0)
}

func (m *MockRedisManager) LPush(ctx context.Context, key string, values ...interface{}) error {
	args := m.Called(ctx, key, values)
	return args.Error(0)
}

func (m *MockRedisManager) RPush(ctx context.Context, key string, values ...interface{}) error {
	args := m.Called(ctx, key, values)
	return args.Error(0)
}

func (m *MockRedisManager) LPop(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisManager) RPop(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisManager) LLen(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisManager) SAdd(ctx context.Context, key string, members ...interface{}) error {
	args := m.Called(ctx, key, members)
	return args.Error(0)
}

func (m *MockRedisManager) SMembers(ctx context.Context, key string) ([]string, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRedisManager) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	args := m.Called(ctx, key, member)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisManager) SRem(ctx context.Context, key string, members ...interface{}) error {
	args := m.Called(ctx, key, members)
	return args.Error(0)
}

func (m *MockRedisManager) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	args := m.Called(ctx, key, members)
	return args.Error(0)
}

func (m *MockRedisManager) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	args := m.Called(ctx, key, start, stop)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRedisManager) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	args := m.Called(ctx, key, start, stop)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]redis.Z), args.Error(1)
}

func (m *MockRedisManager) ZRem(ctx context.Context, key string, members ...interface{}) error {
	args := m.Called(ctx, key, members)
	return args.Error(0)
}

func (m *MockRedisManager) Incr(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisManager) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	args := m.Called(ctx, key, value)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisManager) Decr(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}
