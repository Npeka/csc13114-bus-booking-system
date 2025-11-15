package db

import (
	"context"
	"fmt"
	"time"

	"csc13114-bus-ticket-booking-system/shared/config"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// RedisManager handles Redis connections and operations
type RedisManager struct {
	Client *redis.Client
	Config *config.RedisConfig
}

// NewRedisConnection creates a new Redis connection
func NewRedisConnection(cfg *config.RedisConfig) (*RedisManager, error) {
	// Redis client options
	options := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,

		// Connection pool settings
		PoolTimeout: 30 * time.Second,

		// TLS Config can be added here if needed
		// TLSConfig: &tls.Config{},
	}

	// Create Redis client
	client := redis.NewClient(options)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Info().
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Int("db", cfg.DB).
		Msg("Successfully connected to Redis")

	return &RedisManager{
		Client: client,
		Config: cfg,
	}, nil
}

// Close closes the Redis connection
func (rm *RedisManager) Close() error {
	if rm.Client != nil {
		return rm.Client.Close()
	}
	return nil
}

// HealthCheck performs a health check on Redis
func (rm *RedisManager) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rm.Client.Ping(ctx).Result()
	return err
}

// GetStats returns Redis client statistics
func (rm *RedisManager) GetStats() *redis.PoolStats {
	if rm.Client == nil {
		return nil
	}
	return rm.Client.PoolStats()
}

// Set sets a key-value pair with expiration
func (rm *RedisManager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rm.Client.Set(ctx, key, value, expiration).Err()
}

// Get gets a value by key
func (rm *RedisManager) Get(ctx context.Context, key string) (string, error) {
	return rm.Client.Get(ctx, key).Result()
}

// Del deletes keys
func (rm *RedisManager) Del(ctx context.Context, keys ...string) error {
	return rm.Client.Del(ctx, keys...).Err()
}

// Exists checks if keys exist
func (rm *RedisManager) Exists(ctx context.Context, keys ...string) (int64, error) {
	return rm.Client.Exists(ctx, keys...).Result()
}

// Expire sets expiration for a key
func (rm *RedisManager) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return rm.Client.Expire(ctx, key, expiration).Err()
}

// TTL gets the time to live for a key
func (rm *RedisManager) TTL(ctx context.Context, key string) (time.Duration, error) {
	return rm.Client.TTL(ctx, key).Result()
}

// HSet sets hash field
func (rm *RedisManager) HSet(ctx context.Context, key string, values ...interface{}) error {
	return rm.Client.HSet(ctx, key, values...).Err()
}

// HGet gets hash field value
func (rm *RedisManager) HGet(ctx context.Context, key, field string) (string, error) {
	return rm.Client.HGet(ctx, key, field).Result()
}

// HGetAll gets all hash fields and values
func (rm *RedisManager) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return rm.Client.HGetAll(ctx, key).Result()
}

// HDel deletes hash fields
func (rm *RedisManager) HDel(ctx context.Context, key string, fields ...string) error {
	return rm.Client.HDel(ctx, key, fields...).Err()
}

// LPush pushes elements to the head of a list
func (rm *RedisManager) LPush(ctx context.Context, key string, values ...interface{}) error {
	return rm.Client.LPush(ctx, key, values...).Err()
}

// RPush pushes elements to the tail of a list
func (rm *RedisManager) RPush(ctx context.Context, key string, values ...interface{}) error {
	return rm.Client.RPush(ctx, key, values...).Err()
}

// LPop pops an element from the head of a list
func (rm *RedisManager) LPop(ctx context.Context, key string) (string, error) {
	return rm.Client.LPop(ctx, key).Result()
}

// RPop pops an element from the tail of a list
func (rm *RedisManager) RPop(ctx context.Context, key string) (string, error) {
	return rm.Client.RPop(ctx, key).Result()
}

// LLen gets the length of a list
func (rm *RedisManager) LLen(ctx context.Context, key string) (int64, error) {
	return rm.Client.LLen(ctx, key).Result()
}

// SAdd adds members to a set
func (rm *RedisManager) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return rm.Client.SAdd(ctx, key, members...).Err()
}

// SMembers gets all members of a set
func (rm *RedisManager) SMembers(ctx context.Context, key string) ([]string, error) {
	return rm.Client.SMembers(ctx, key).Result()
}

// SIsMember checks if a member exists in a set
func (rm *RedisManager) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return rm.Client.SIsMember(ctx, key, member).Result()
}

// SRem removes members from a set
func (rm *RedisManager) SRem(ctx context.Context, key string, members ...interface{}) error {
	return rm.Client.SRem(ctx, key, members...).Err()
}

// ZAdd adds members to a sorted set
func (rm *RedisManager) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return rm.Client.ZAdd(ctx, key, members...).Err()
}

// ZRange gets members from a sorted set by range
func (rm *RedisManager) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return rm.Client.ZRange(ctx, key, start, stop).Result()
}

// ZRangeWithScores gets members with scores from a sorted set by range
func (rm *RedisManager) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return rm.Client.ZRangeWithScores(ctx, key, start, stop).Result()
}

// ZRem removes members from a sorted set
func (rm *RedisManager) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return rm.Client.ZRem(ctx, key, members...).Err()
}

// Incr increments a key's value
func (rm *RedisManager) Incr(ctx context.Context, key string) (int64, error) {
	return rm.Client.Incr(ctx, key).Result()
}

// IncrBy increments a key's value by a specific amount
func (rm *RedisManager) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return rm.Client.IncrBy(ctx, key, value).Result()
}

// Decr decrements a key's value
func (rm *RedisManager) Decr(ctx context.Context, key string) (int64, error) {
	return rm.Client.Decr(ctx, key).Result()
}

// DecrBy decrements a key's value by a specific amount
func (rm *RedisManager) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	return rm.Client.DecrBy(ctx, key, value).Result()
}

// FlushDB flushes the current database
func (rm *RedisManager) FlushDB(ctx context.Context) error {
	return rm.Client.FlushDB(ctx).Err()
}

// FlushAll flushes all databases
func (rm *RedisManager) FlushAll(ctx context.Context) error {
	return rm.Client.FlushAll(ctx).Err()
}
